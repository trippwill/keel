package engine

import "testing"

type testFrame struct {
	ExtentConstraint
	id  string
	fit FitMode
}

func (f testFrame) ID() string { return f.id }

func (f testFrame) Fit() FitMode {
	if f.fit == 0 {
		return FitExact
	}
	return f.fit
}

type testStack struct {
	ExtentConstraint
	axis  Axis
	slots []Spec
}

func (s testStack) Axis() Axis { return s.axis }

func (s testStack) Len() int { return len(s.slots) }

func (s testStack) Slot(index int) (Spec, bool) {
	if index < 0 || index >= len(s.slots) {
		return nil, false
	}
	return s.slots[index], true
}

func fixed(units int) ExtentConstraint {
	return ExtentConstraint{Kind: ExtentFixed, Units: units, MinCells: units, MaxCells: 0}
}

func flex(units int) ExtentConstraint {
	return ExtentConstraint{Kind: ExtentFlex, Units: units, MinCells: 0, MaxCells: 0}
}

func TestArrangeBuildsRects(t *testing.T) {
	layout := testStack{
		ExtentConstraint: flex(1),
		axis:             AxisHorizontal,
		slots: []Spec{
			testFrame{ExtentConstraint: fixed(3), id: "a"},
			testStack{
				ExtentConstraint: flex(1),
				axis:             AxisVertical,
				slots: []Spec{
					testFrame{ExtentConstraint: fixed(2), id: "b"},
					testFrame{ExtentConstraint: flex(1), id: "c"},
				},
			},
		},
	}
	size := Size{Width: 10, Height: 5}
	arranged, err := Arrange[string](layout, size, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if arranged.Root.Rect.Width != 10 || arranged.Root.Rect.Height != 5 {
		t.Fatalf("expected root 10x5, got %dx%d", arranged.Root.Rect.Width, arranged.Root.Rect.Height)
	}

	if len(arranged.Root.Slots) != 2 {
		t.Fatalf("expected 2 slots, got %d", len(arranged.Root.Slots))
	}

	left := arranged.Root.Slots[0]
	right := arranged.Root.Slots[1]

	if left.Rect.X != 0 || left.Rect.Width != 3 || left.Rect.Height != 5 {
		t.Fatalf("unexpected left rect: %+v", left.Rect)
	}
	if right.Rect.X != 3 || right.Rect.Width != 7 || right.Rect.Height != 5 {
		t.Fatalf("unexpected right rect: %+v", right.Rect)
	}

	if len(right.Slots) != 2 {
		t.Fatalf("expected 2 slots under right, got %d", len(right.Slots))
	}

	top := right.Slots[0]
	bottom := right.Slots[1]
	if top.Rect.Y != 0 || top.Rect.Height != 2 {
		t.Fatalf("unexpected top rect: %+v", top.Rect)
	}
	if bottom.Rect.Y != 2 || bottom.Rect.Height != 3 {
		t.Fatalf("unexpected bottom rect: %+v", bottom.Rect)
	}
}
