package keel

import (
	"errors"
	"testing"
)

type testContainer struct {
	axis  Axis
	slots []Renderable
}

func (c testContainer) GetAxis() Axis { return c.axis }

func (c testContainer) Len() int { return len(c.slots) }

func (c testContainer) Slot(index int) (Renderable, bool) {
	if index < 0 || index >= len(c.slots) {
		return nil, false
	}
	return c.slots[index], true
}

func (c testContainer) GetExtent() ExtentConstraint { return FlexUnit() }

func TestRowResolverEmptyContainer(t *testing.T) {
	container := testContainer{axis: AxisHorizontal}

	sizes, required, err := RowResolver(10, container)
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

func TestColResolverEmptyContainer(t *testing.T) {
	container := testContainer{axis: AxisVertical}

	sizes, required, err := ColResolver(10, container)
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

func TestRowResolverNilSlot(t *testing.T) {
	container := testContainer{
		axis:  AxisHorizontal,
		slots: []Renderable{nil},
	}

	_, _, err := RowResolver(10, container)
	var slotErr *SlotError
	if !errors.As(err, &slotErr) {
		t.Fatalf("expected SlotError, got %v", err)
	}
	if slotErr.Index != 0 {
		t.Fatalf("expected index 0, got %d", slotErr.Index)
	}
	if !errors.Is(err, ErrNilSlot) {
		t.Fatalf("expected ErrNilSlot")
	}
}

func TestRenderContainerInvalidAxis(t *testing.T) {
	container := testContainer{
		axis:  Axis(99),
		slots: []Renderable{Panel(FlexUnit(), "a")},
	}

	_, err := RenderContainer(container, Context[string]{Width: 10, Height: 1})
	if !errors.Is(err, ErrInvalidAxis) {
		t.Fatalf("expected ErrInvalidAxis, got %v", err)
	}
}
