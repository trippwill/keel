package keel

// SplitSpec defines a stack that splits its allocation along an axis.
type SplitSpec struct {
	ExtentConstraint
	axis Axis
	rs   []Spec
}

var _ StackSpec = SplitSpec{}

// Split creates a new split with the given axis and extent.
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
// Slots are stored as references; mutating slots after creation affects the Split.
// Panics on invalid axis.
func Split(axis Axis, extent ExtentConstraint, slots ...Spec) SplitSpec {
	if (axis != AxisHorizontal) && (axis != AxisVertical) {
		panic(ErrInvalidAxis)
	}

	return SplitSpec{
		ExtentConstraint: extent,
		axis:             axis,
		rs:               slots,
	}
}

// Row creates a new horizontal split.
// Slots are stored as references; mutating slots after creation affects the Split.
func Row(size ExtentConstraint, slots ...Spec) SplitSpec {
	return Split(AxisHorizontal, size, slots...)
}

// Col creates a new vertical split.
// Slots are stored as references; mutating slots after creation affects the Split.
func Col(size ExtentConstraint, slots ...Spec) SplitSpec {
	return Split(AxisVertical, size, slots...)
}

// Axis implements [StackSpec].
func (s SplitSpec) Axis() Axis { return s.axis }

// Len implements [StackSpec].
func (s SplitSpec) Len() int { return len(s.rs) }

// Slot implements [StackSpec].
func (s SplitSpec) Slot(index int) (Spec, bool) {
	if index < 0 || index >= len(s.rs) {
		return nil, false
	}

	return s.rs[index], true
}
