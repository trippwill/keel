package engine

import (
	"fmt"
	"testing"
)

type benchStack struct {
	slots []Spec
}

func (s *benchStack) Len() int {
	return len(s.slots)
}

func (s *benchStack) Slot(i int) (Spec, bool) {
	if i < 0 || i >= len(s.slots) {
		return nil, false
	}
	return s.slots[i], true
}

func (s *benchStack) Axis() Axis {
	return AxisHorizontal
}

func (s *benchStack) Extent() ExtentConstraint {
	return ExtentConstraint{Kind: ExtentFlex, MinCells: 1, MaxCells: 1}
}

type benchSpec struct {
	extent ExtentConstraint
}

func (b benchSpec) Extent() ExtentConstraint {
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

func benchSlots(count int) ([]Spec, int) {
	slots := make([]Spec, count)
	required := 0
	for i := range count {
		var extent ExtentConstraint
		if i%3 == 0 {
			extent = ExtentConstraint{Kind: ExtentFixed, Units: 3}
		} else {
			extent = ExtentConstraint{Kind: ExtentFlex, Units: 1, MinCells: 1}
		}
		slots[i] = benchSpec{extent: extent}
		switch extent.Kind {
		case ExtentFixed:
			required += extent.Units
		case ExtentFlex:
			required += extent.MinCells
		}
	}
	return slots, required
}
