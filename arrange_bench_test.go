package keel_test

import (
	"fmt"
	"testing"

	"github.com/trippwill/keel"
)

type benchStack struct {
	slots []keel.Spec
}

func (s *benchStack) Len() int {
	return len(s.slots)
}

func (s *benchStack) Slot(i int) (keel.Spec, bool) {
	if i < 0 || i >= len(s.slots) {
		return nil, false
	}
	return s.slots[i], true
}

func (s *benchStack) Axis() keel.Axis {
	return keel.AxisHorizontal
}

func (s *benchStack) Extent() keel.ExtentConstraint {
	return keel.FlexMin(1, 1)
}

type benchSpec struct {
	extent keel.ExtentConstraint
}

func (b benchSpec) Extent() keel.ExtentConstraint {
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
				if _, _, err := keel.ArrangeStack(total, stack); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func benchSlots(count int) ([]keel.Spec, int) {
	slots := make([]keel.Spec, count)
	required := 0
	for i := range count {
		var extent keel.ExtentConstraint
		if i%3 == 0 {
			extent = keel.Fixed(3)
		} else {
			extent = keel.FlexMin(1, 1)
		}
		slots[i] = benchSpec{extent: extent}
		switch extent.Kind {
		case keel.ExtentFixed:
			required += extent.Units
		case keel.ExtentFlex:
			required += extent.MinCells
		}
	}
	return slots, required
}
