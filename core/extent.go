//go:generate stringer -type=ExtentKind -trimprefix=Extent
package core

// ExtentKind represents whether an [ExtentConstraint] is fixed or flexible.
type ExtentKind uint8

const (
	// ExtentFixed represents a fixed-size extent.
	ExtentFixed ExtentKind = iota
	// ExtentFlex represents a flexible extent.
	ExtentFlex
)

// ExtentConstraint defines how much total space a [Spec] should take along an axis.
// For frames, this is the allocation for content plus any frame (padding, border, margin).
// For stacks, this is the space available to distribute across slots.
type ExtentConstraint struct {
	Kind     ExtentKind
	Units    int
	MinCells int // Minimum total cells to reserve on this axis (0 = no min)
	MaxCells int // Maximum total cells to reserve on this axis (0 = no max)
}

// Extent implements the [Spec] interface.
func (e ExtentConstraint) Extent() ExtentConstraint {
	return e
}
