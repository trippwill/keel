package keel

// Renderable is an element that can be rendered within a layout tree.
type Renderable interface {
	GetExtent() ExtentConstraint // Returns the desired total extent along an axis.
}

// Block is a [Renderable] with an ID, used for content and styling.
type Block[KID KeelID] interface {
	Renderable
	GetID() KID
	GetClip() ClipConstraint // Returns the content clip (max content size).
}

// Container is a [Renderable] that contains Slots for other [Renderable]s.
type Container interface {
	Renderable
	GetAxis() Axis                     // Layout axis of the container
	Len() int                          // Number of slots in the container
	Slot(index int) (Renderable, bool) // Slot access (ok=false when out of range); must be stable during a render call
}
