package keel

// Renderable is an element that participates in the layout hierarchy.
// The extent returned is the total allocation along the container axis,
// including any frame space for blocks.
type Renderable interface {
	GetExtent() ExtentConstraint // Returns the desired total extent along an axis.
}

// Block is a [Renderable] with an ID, used for content and styling.
// Blocks are the only renderables that produce output content.
type Block[KID KeelID] interface {
	Renderable
	GetID() KID
	GetFit() FitMode // Returns the content fit mode.
}

// Container is a [Renderable] that splits its allocation across slots.
// Slot access must be stable for the duration of a resolve pass.
type Container interface {
	Renderable
	GetAxis() Axis                     // Layout axis of the container
	Len() int                          // Number of slots in the container
	Slot(index int) (Renderable, bool) // Slot access (ok=false when out of range); must be stable during a render call
}
