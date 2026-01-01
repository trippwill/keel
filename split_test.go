package keel

import (
	"errors"
	"testing"

	"github.com/trippwill/keel/engine"
)

type testStack struct {
	axis  engine.Axis
	slots []Spec
}

func (s testStack) Axis() engine.Axis { return s.axis }

func (s testStack) Len() int { return len(s.slots) }

func (s testStack) Slot(index int) (Spec, bool) {
	if index < 0 || index >= len(s.slots) {
		return nil, false
	}
	return s.slots[index], true
}

func (s testStack) Extent() engine.ExtentConstraint { return FlexUnit() }

func TestArrangeStackEmptyStack(t *testing.T) {
	stack := testStack{axis: engine.AxisHorizontal}

	sizes, required, err := engine.ArrangeStack(10, stack)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sizes) != 0 {
		t.Fatalf("expected empty sizes, got %v", sizes)
	}
	if required != 0 {
		t.Fatalf("expected required 0, got %d", required)
	}
}

func TestArrangeStackEmptyStackVertical(t *testing.T) {
	stack := testStack{axis: engine.AxisVertical}

	sizes, required, err := engine.ArrangeStack(10, stack)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sizes) != 0 {
		t.Fatalf("expected empty sizes, got %v", sizes)
	}
	if required != 0 {
		t.Fatalf("expected required 0, got %d", required)
	}
}

func TestArrangeStackNilChild(t *testing.T) {
	stack := testStack{
		axis:  engine.AxisHorizontal,
		slots: []Spec{nil},
	}

	_, _, err := engine.ArrangeStack(10, stack)
	var slotErr *engine.SlotError
	if !errors.As(err, &slotErr) {
		t.Fatalf("expected SlotError, got %v", err)
	}
	if slotErr.Index != 0 {
		t.Fatalf("expected index 0, got %d", slotErr.Index)
	}
	if !errors.Is(err, engine.ErrNilSlot) {
		t.Fatalf("expected ErrNilSlot")
	}
}

func TestRenderStackInvalidAxis(t *testing.T) {
	stack := testStack{
		axis:  engine.Axis(99),
		slots: []Spec{Panel(FlexUnit(), "a")},
	}

	renderer := NewRenderer[string](stack, nil, nil)
	size := Size{Width: 10, Height: 1}
	_, err := renderer.Render(size)
	if !errors.Is(err, engine.ErrInvalidAxis) {
		t.Fatalf("expected ErrInvalidAxis, got %v", err)
	}
}
