package keel

import "github.com/trippwill/keel/engine"

// SplitSpec defines a stack that splits its allocation along an axis.
type SplitSpec struct {
	engine.ExtentConstraint
	axis engine.Axis
	rs   []Spec
}

var _ engine.StackSpec = SplitSpec{}

// Split creates a new split with the given axis and extent.
//
// Arguments:
//
//	axis:   engine.Axis along which to split (horizontal or vertical)
//	extent: Total extent constraint for the split along the stack axis
//	slots:  Slot specifications to include in the split
//
// Returns:
//   - A new [SplitSpec] configured with the provided arguments.
//
// Slots are stored as references; mutating slots after creation affects the Split.
// Panics on invalid axis.
func Split(axis engine.Axis, extent engine.ExtentConstraint, slots ...Spec) SplitSpec {
	if (axis != engine.AxisHorizontal) && (axis != engine.AxisVertical) {
		panic(engine.ErrInvalidAxis)
	}

	return SplitSpec{
		ExtentConstraint: extent,
		axis:             axis,
		rs:               slots,
	}
}

// Row creates a new horizontal split.
// Slots are stored as references; mutating slots after creation affects the Split.
func Row(size engine.ExtentConstraint, slots ...Spec) SplitSpec {
	return Split(engine.AxisHorizontal, size, slots...)
}

// Col creates a new vertical split.
// Slots are stored as references; mutating slots after creation affects the Split.
func Col(size engine.ExtentConstraint, slots ...Spec) SplitSpec {
	return Split(engine.AxisVertical, size, slots...)
}

// Axis implements [StackSpec].
func (s SplitSpec) Axis() engine.Axis { return s.axis }

// Len implements [StackSpec].
func (s SplitSpec) Len() int { return len(s.rs) }

// Slot implements [StackSpec].
func (s SplitSpec) Slot(index int) (Spec, bool) {
	if index < 0 || index >= len(s.rs) {
		return nil, false
	}

	return s.rs[index], true
}
