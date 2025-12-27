//go:generate stringer -type=Axis -trimprefix=Axis
//go:generate stringer -type=ExtentKind
package keel

// Axis represents a layout axis used by containers to split space.
type Axis uint8

const (
	// AxisHorizontal lays out content left-to-right.
	AxisHorizontal Axis = 0
	// AxisVertical lays out content top-to-bottom.
	AxisVertical Axis = 1
)

// ExtentKind represents whether an [ExtentConstraint] is fixed or flexible.
type ExtentKind uint8

const (
	// ExtentFixed represents a fixed-size extent.
	ExtentFixed ExtentKind = iota
	// ExtentFlex represents a flexible extent.
	ExtentFlex
)

// ExtentConstraint defines how much total space a renderable should take along an axis.
// For blocks, this is the allocation for content plus any frame (padding, border, margin).
// For containers, this is the space available to distribute across slots.
type ExtentConstraint struct {
	Kind     ExtentKind
	Units    int
	MinCells int // Minimum total cells to reserve on this axis
}

// FlexUnit returns a single flexible unit of space with no minimum.
func FlexUnit() ExtentConstraint {
	return ExtentConstraint{ExtentFlex, 1, 0}
}

// Fixed creates a fixed [ExtentConstraint] with the given units in cells.
// Fixed extents must be at least their MinCells value.
func Fixed(units int) ExtentConstraint {
	return ExtentConstraint{ExtentFixed, units, units}
}

// Flex creates a flexible [ExtentConstraint] with the given units in flex space.
// Flex extents receive space after fixed and minimum allocations are satisfied.
func Flex(units int) ExtentConstraint {
	return ExtentConstraint{ExtentFlex, units, 0}
}

// FlexMin creates a flexible [ExtentConstraint] with the given units in flex space,
// and reserves at least minReserved total cells along the container axis.
func FlexMin(units int, minReserved int) ExtentConstraint {
	return ExtentConstraint{ExtentFlex, units, minReserved}
}

// ClipConstraint caps the width and height of the content box.
// Clipping is applied to content before sizing checks, and the clipped content
// must still fit within the content box.
// Zero values indicate no constraint.
type ClipConstraint struct {
	Width, Height int
}

// Clip returns a clip constraint for width and height.
func Clip(width, height int) ClipConstraint {
	return ClipConstraint{Width: width, Height: height}
}

// ClipWidth returns a clip constraint that limits width only.
func ClipWidth(width int) ClipConstraint {
	return ClipConstraint{Width: width, Height: 0}
}

// ClipHeight returns a clip constraint that limits height only.
func ClipHeight(height int) ClipConstraint {
	return ClipConstraint{Width: 0, Height: height}
}

// GetExtent implements the [Renderable] interface.
func (e *ExtentConstraint) GetExtent() ExtentConstraint {
	return *e
}

// GetClip implements the [Block] interface.
func (mc *ClipConstraint) GetClip() ClipConstraint {
	return *mc
}

// RenderInfo describes the allocated space for a [Block] render pass.
type RenderInfo struct {
	Width, Height               int            // Total allocated size
	ContentWidth, ContentHeight int            // Inner content box size
	FrameWidth, FrameHeight     int            // Total frame size (padding + border + margin)
	Clip                        ClipConstraint // Clipping constraint if any
}
