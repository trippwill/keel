package engine

import (
	"fmt"
	"testing"

	"github.com/trippwill/keel/core"
)

type benchStack struct {
	slots []core.Spec
}

func (s *benchStack) Len() int {
	return len(s.slots)
}

func (s *benchStack) Slot(i int) (core.Spec, bool) {
	if i < 0 || i >= len(s.slots) {
		return nil, false
	}
	return s.slots[i], true
}

func (s *benchStack) Axis() core.Axis {
	return core.AxisHorizontal
}

func (s *benchStack) Extent() core.ExtentConstraint {
	return core.ExtentConstraint{Kind: core.ExtentFlex, MinCells: 1, MaxCells: 1}
}

type benchSpec struct {
	extent core.ExtentConstraint
}

func (b benchSpec) Extent() core.ExtentConstraint {
	return b.extent
}

func BenchmarkArrangeStack(b *testing.B) {
	for _, count := range []int{3, 8, 32, 128} {
		b.Run(fmt.Sprintf("n=%d", count), func(b *testing.B) {
			slots, required := benchSlots(count)
			total := required + count
			stack := &benchStack{slots: slots}

			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				if _, _, err := ArrangeStack(total, stack); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func benchSlots(count int) ([]core.Spec, int) {
	slots := make([]core.Spec, count)
	required := 0
	for i := range count {
		var extent core.ExtentConstraint
		if i%3 == 0 {
			extent = core.ExtentConstraint{Kind: core.ExtentFixed, Units: 3}
		} else {
			extent = core.ExtentConstraint{Kind: core.ExtentFlex, Units: 1, MinCells: 1}
		}
		slots[i] = benchSpec{extent: extent}
		switch extent.Kind {
		case core.ExtentFixed:
			required += extent.Units
		case core.ExtentFlex:
			required += extent.MinCells
		}
	}
	return slots, required
}
