package keel

// SplitSpec defines a container that splits its allocation along an axis.
type SplitSpec struct {
	ExtentConstraint
	axis Axis
	rs   []Renderable
}

var (
	_ Container  = (*SplitSpec)(nil)
	_ Renderable = (*SplitSpec)(nil)
)

// Split creates a new split with the given axis and extent.
// Slots are stored as references; mutating slots after creation affects the Split.
// Panics on invalid axis.
func Split(axis Axis, extent ExtentConstraint, slots ...Renderable) *SplitSpec {
	if (axis != AxisHorizontal) && (axis != AxisVertical) {
		panic(ErrInvalidAxis)
	}

	return &SplitSpec{
		ExtentConstraint: extent,
		axis:             axis,
		rs:               slots,
	}
}

// Row creates a new horizontal split.
// Slots are stored as references; mutating slots after creation affects the Split.
func Row(size ExtentConstraint, slots ...Renderable) *SplitSpec {
	return Split(AxisHorizontal, size, slots...)
}

// Col creates a new vertical split.
// Slots are stored as references; mutating slots after creation affects the Split.
func Col(size ExtentConstraint, slots ...Renderable) *SplitSpec {
	return Split(AxisVertical, size, slots...)
}

// GetAxis implements [Container].
func (s *SplitSpec) GetAxis() Axis { return s.axis }

// Len implements [Container].
func (s *SplitSpec) Len() int { return len(s.rs) }

// Slot implements [Container].
func (s *SplitSpec) Slot(index int) (Renderable, bool) {
	if index < 0 || index >= len(s.rs) {
		return nil, false
	}

	return s.rs[index], true
}
