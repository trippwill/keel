package keel

import "testing"

func TestArrangeBuildsRects(t *testing.T) {
	layout := Row(FlexUnit(),
		Panel(Fixed(3), "a"),
		Col(FlexUnit(),
			Panel(Fixed(2), "b"),
			Panel(FlexUnit(), "c"),
		),
	)

	renderer := NewRenderer[string](layout, nil, nil)
	size := Size{Width: 10, Height: 5}
	arranged, err := arrange(renderer, layout, size)
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

func TestRenderMatchesLayoutRender(t *testing.T) {
	layout := Row(FlexUnit(),
		Panel(Fixed(2), "a"),
		Panel(FlexUnit(), "b"),
	)

	renderer := NewRenderer[string](layout, nil, func(id string, _ FrameInfo) (string, error) {
		switch id {
		case "a":
			return "aa", nil
		case "b":
			return "bbb", nil
		default:
			return "", &UnknownFrameIDError{ID: id}
		}
	})
	size := Size{Width: 5, Height: 1}

	want, err := renderer.Render(size)
	if err != nil {
		t.Fatalf("unexpected render error: %v", err)
	}

	arranged, err := arrange(renderer, layout, size)
	if err != nil {
		t.Fatalf("unexpected arrange error: %v", err)
	}

	got, err := renderer.renderLayout(arranged)
	if err != nil {
		t.Fatalf("unexpected layout render error: %v", err)
	}
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
