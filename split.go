package keel

// SplitSpec defines a layout container that splits space along an axis.
type SplitSpec[KID KeelID] struct {
	id    KID
	axis  Axis
	slots []Slot[KID]
}

var (
	_ Container[string]  = (*SplitSpec[string])(nil)
	_ Renderable[string] = (*SplitSpec[string])(nil)
)

// Split creates a new split with the given axis.
// Empty slots will panic.
func Split[KID KeelID](id KID, axis Axis, slots ...Slot[KID]) *SplitSpec[KID] {
	if (axis != AxisHorizontal) && (axis != AxisVertical) {
		panic("invalid axis for Split")
	}

	return &SplitSpec[KID]{
		id:    id,
		axis:  axis,
		slots: copySlots(slots),
	}
}

// Row creates a new horizontal split.
// Empty slots will panic.
func Row[KID KeelID](id KID, slots ...Slot[KID]) *SplitSpec[KID] {
	return &SplitSpec[KID]{
		id:    id,
		axis:  AxisHorizontal,
		slots: copySlots(slots),
	}
}

// Col creates a new vertical split.
// Empty slots will panic.
func Col[KID KeelID](id KID, slots ...Slot[KID]) *SplitSpec[KID] {
	return &SplitSpec[KID]{
		id:    id,
		axis:  AxisVertical,
		slots: copySlots(slots),
	}
}

func (s *SplitSpec[KID]) GetID() KID    { return s.id }
func (s *SplitSpec[KID]) GetAxis() Axis { return s.axis }
func (s *SplitSpec[KID]) Len() int      { return len(s.slots) }

func (s *SplitSpec[KID]) Slot(index int) (SizeSpec, Renderable[KID], bool) {
	if index < 0 || index >= len(s.slots) {
		return SizeSpec{}, nil, false
	}

	slot := s.slots[index]
	return slot.Size, slot.Node, true
}

// Render implements [Renderable].
func (s *SplitSpec[KID]) Render(cxt Context[KID]) (string, error) {
	return renderSplit(s, cxt)
}
