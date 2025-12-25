package keel

// Renderable is an element that can be rendered within a layout tree.
type Renderable[KID KeelID] interface {
	GetID() KID                          // Unique identifier for the Renderable
	Render(Context[KID]) (string, error) // Renders the element within the given context
}

// Container is a Renderable that contains Slots for other Renderables.
type Container[KID KeelID] interface {
	GetAxis() Axis                                    // Layout axis of the container
	Len() int                                         // Number of slots in the container
	Slot(index int) (SizeSpec, Renderable[KID], bool) // Slot access (ok=false when out of range)
}

// Slot defines a space in a Container for a Renderable.
type Slot[KID KeelID] struct {
	Size SizeSpec        // Size specification for the slot
	Node Renderable[KID] // Renderable contained in the slot
}

func copySlots[KID KeelID](slots []Slot[KID]) []Slot[KID] {
	if len(slots) == 0 {
		panic("slots cannot be empty")
	}

	copied := make([]Slot[KID], len(slots))
	for i, slot := range slots {
		if slot.Node == nil {
			panic("Slot node cannot be nil")
		}
		copied[i] = slot
	}

	return copied
}
