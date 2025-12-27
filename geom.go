//go:generate stringer -type=Axis -trimprefix=Axis
//go:generate stringer -type=ExtentKind
package keel

// Axis represents a layout axis: horizontal or vertical.
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
type ExtentConstraint struct {
	Kind     ExtentKind
	Units    int
	MinCells int // Minimum total cells to reserve on this axis
}

// FlexUnit returns a single flexible unit of space.
func FlexUnit() ExtentConstraint {
	return ExtentConstraint{ExtentFlex, 1, 0}
}

// Fixed creates a fixed [ExtentConstraint] with the given units in cells.
func Fixed(units int) ExtentConstraint {
	return ExtentConstraint{ExtentFixed, units, units}
}

// Flex creates a flexible [ExtentConstraint] with the given units in flex space.
func Flex(units int) ExtentConstraint {
	return ExtentConstraint{ExtentFlex, units, 0}
}

// FlexMin creates a flexible [ExtentConstraint] with the given units in flex space,
// and reserves at least minReserved total cells along the container axis.
func FlexMin(units int, minReserved int) ExtentConstraint {
	return ExtentConstraint{ExtentFlex, units, minReserved}
}

// ClipConstraint caps the width and height of the content box.
// Clipping is applied before rendering and the clipped content must fit.
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
