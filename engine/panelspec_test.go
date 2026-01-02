package engine

import (
	"testing"

	"github.com/trippwill/keel/core"
)

func TestPanelSpecAccessors(t *testing.T) {
	extent := core.ExtentConstraint{Kind: core.ExtentFixed, Units: 3, MinCells: 3, MaxCells: 0}
	spec := NewPanelSpec(extent, core.FitClip, "id")
	if got := spec.Extent(); got != extent {
		t.Fatalf("unexpected extent: %+v", got)
	}
	if got := spec.Fit(); got != core.FitClip {
		t.Fatalf("unexpected fit: %v", got)
	}
	if got := spec.ID(); got != "id" {
		t.Fatalf("unexpected id: %v", got)
	}
}
