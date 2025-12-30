package keel

// Spec is an element that participates in the layout hierarchy.
// The extent returned is the total allocation along the stack axis,
// including any frame space for frames.
type Spec interface {
	Extent() ExtentConstraint // Returns the desired total extent along an axis.
}

// FrameSpec is a [Spec] with an ID, used for content and styling.
// Frames are the only specs that produce output content.
type FrameSpec[KID KeelID] interface {
	Spec
	ID() KID
	Fit() FitMode // Returns the content fit mode.
}

// StackSpec is a [Spec] that splits its allocation across slots.
// Slot access must be stable for the duration of an arrange pass.
type StackSpec interface {
	Spec
	Axis() Axis                  // Layout axis of the stack
	Len() int                    // Number of slots in the stack
	Slot(index int) (Spec, bool) // Slot access (ok=false when out of range); must be stable during an arrange call
}
