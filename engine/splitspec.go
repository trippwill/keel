package engine

import "github.com/trippwill/keel/core"

// SplitSpec defines a stack that splits its allocation along an axis.
type SplitSpec struct {
	core.ExtentConstraint
	axis core.Axis
	rs   []core.Spec
}

// NewSplitSpec creates a new split with the given axis and extent.
//
// Arguments:
//
//	axis:   Axis along which to split (horizontal or vertical)
//	extent: Total extent constraint for the split along the stack axis
//	slots:  Slot specifications to include in the split
//
// Returns:
//   - A new [SplitSpec] configured with the provided arguments.
//
// Slots are stored as references; mutating slots after creation affects the NewSplitSpec.
// Panics on invalid axis.
func NewSplitSpec(axis core.Axis, extent core.ExtentConstraint, slots ...core.Spec) SplitSpec {
	if (axis != core.AxisHorizontal) && (axis != core.AxisVertical) {
		panic(core.ErrInvalidAxis)
	}

	return SplitSpec{
		ExtentConstraint: extent,
		axis:             axis,
		rs:               slots,
	}
}

// Axis implements [StackSpec].
func (s SplitSpec) Axis() core.Axis { return s.axis }

// Len implements [StackSpec].
func (s SplitSpec) Len() int { return len(s.rs) }

// Slot implements [StackSpec].
func (s SplitSpec) Slot(index int) (core.Spec, bool) {
	if index < 0 || index >= len(s.rs) {
		return nil, false
	}

	return s.rs[index], true
}
