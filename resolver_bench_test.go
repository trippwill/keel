package keel_test

import (
	"fmt"
	"testing"

	"github.com/trippwill/chiplog/keel"
)

type benchContainer struct {
	slots []keel.Renderable
}

func (c *benchContainer) Len() int {
	return len(c.slots)
}

func (c *benchContainer) Slot(i int) (keel.Renderable, bool) {
	if i < 0 || i >= len(c.slots) {
		return nil, false
	}
	return c.slots[i], true
}

func (c *benchContainer) GetAxis() keel.Axis {
	return keel.AxisHorizontal
}

func (c *benchContainer) GetExtent() keel.ExtentConstraint {
	return keel.FlexMin(1, 1)
}

type benchRenderable struct {
	extent keel.ExtentConstraint
}

func (b benchRenderable) GetExtent() keel.ExtentConstraint {
	return b.extent
}

func BenchmarkResolveContainer(b *testing.B) {
	for _, count := range []int{3, 8, 32, 128} {
		b.Run(fmt.Sprintf("n=%d", count), func(b *testing.B) {
			slots, required := benchSlots(count)
			total := required + count
			container := &benchContainer{slots: slots}

			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				if _, _, err := keel.RowResolver(total, container); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func benchSlots(count int) ([]keel.Renderable, int) {
	slots := make([]keel.Renderable, count)
	required := 0
	for i := range count {
		var extent keel.ExtentConstraint
		if i%3 == 0 {
			extent = keel.Fixed(3)
		} else {
			extent = keel.FlexMin(1, 1)
		}
		slots[i] = benchRenderable{extent: extent}
		switch extent.Kind {
		case keel.ExtentFixed:
			required += extent.Units
		case keel.ExtentFlex:
			required += extent.MinCells
		}
	}
	return slots, required
}
