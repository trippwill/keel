package engine

import (
	"testing"

	"github.com/trippwill/keel/core"
)

func TestSplitSpecAccessors(t *testing.T) {
	extent := core.ExtentConstraint{Kind: core.ExtentFlex, Units: 1, MinCells: 0, MaxCells: 0}
	child := NewPanelSpec(core.ExtentConstraint{Kind: core.ExtentFixed, Units: 2, MinCells: 2, MaxCells: 0}, core.FitExact, "a")
	spec := NewSplitSpec(core.AxisHorizontal, extent, child)
	if got := spec.Axis(); got != core.AxisHorizontal {
		t.Fatalf("unexpected axis: %v", got)
	}
	if got := spec.Extent(); got != extent {
		t.Fatalf("unexpected extent: %+v", got)
	}
	if got := spec.Len(); got != 1 {
		t.Fatalf("expected 1 slot, got %d", got)
	}
	slot, ok := spec.Slot(0)
	if !ok || slot == nil {
		t.Fatalf("expected slot")
	}
	if _, ok := spec.Slot(1); ok {
		t.Fatalf("expected out-of-range slot")
	}
}

func TestSplitSpecInvalidAxisPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic")
		}
	}()
	_ = NewSplitSpec(core.Axis(99), core.ExtentConstraint{Kind: core.ExtentFlex, Units: 1}, nil)
}
