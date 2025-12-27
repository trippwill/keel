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

func BenchmarkResolveExtentAt(b *testing.B) {
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

func BenchmarkResolveCachedExtents(b *testing.B) {
	for _, count := range []int{3, 8, 32, 128} {
		b.Run(fmt.Sprintf("n=%d", count), func(b *testing.B) {
			slots, required := benchSlots(count)
			total := required + count

			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				extents := make([]keel.ExtentConstraint, count)
				for i := range slots {
					extents[i] = slots[i].GetExtent()
				}
				if _, _, err := resolveSlice(total, extents); err != nil {
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

func resolveSlice(total int, extents []keel.ExtentConstraint) ([]int, int, error) {
	if total < 0 {
		return nil, 0, &keel.ConfigError{Reason: keel.ErrInvalidTotal}
	}
	if len(extents) == 0 {
		return nil, 0, &keel.ConfigError{Reason: keel.ErrEmptyExtents}
	}

	sizes := make([]int, len(extents))
	flexUnits := 0
	required := 0
	hasFlex := false

	for i, spec := range extents {
		if spec.Units <= 0 {
			return nil, required, &keel.ExtentError{Index: i, Reason: keel.ErrInvalidExtentUnits}
		}
		if spec.MinCells < 0 {
			return nil, required, &keel.ExtentError{Index: i, Reason: keel.ErrInvalidExtentMinCells}
		}

		switch spec.Kind {
		case keel.ExtentFixed:
			if spec.Units < spec.MinCells {
				return nil, required, &keel.ExtentError{Index: i, Reason: keel.ErrInvalidExtentMin}
			}
			sizes[i] = spec.Units
		case keel.ExtentFlex:
			sizes[i] = spec.MinCells
			flexUnits += spec.Units
			hasFlex = true
		default:
			return nil, required, &keel.ExtentError{Index: i, Reason: keel.ErrInvalidExtentKind}
		}

		required += sizes[i]
	}

	if required > total {
		return nil, required, keel.ErrExtentTooSmall
	}

	leftover := total - required
	if !hasFlex {
		if leftover > 0 {
			sizes[len(sizes)-1] += leftover
		}
		return sizes, required, nil
	}

	remainder := leftover
	for i, spec := range extents {
		if spec.Kind != keel.ExtentFlex {
			continue
		}
		add := 0
		if flexUnits > 0 {
			add = leftover * spec.Units / flexUnits
		}
		sizes[i] += add
		remainder -= add
	}

	for i := 0; i < len(extents) && remainder > 0; i++ {
		if extents[i].Kind != keel.ExtentFlex {
			continue
		}
		sizes[i]++
		remainder--
	}

	return sizes, required, nil
}
