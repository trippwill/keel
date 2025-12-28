//go:generate stringer -type=Axis -trimprefix=Axis
//go:generate stringer -type=ExtentKind -trimprefix=Extent
//go:generate stringer -type=FitMode -trimprefix=Fit
package keel

// Axis represents a layout axis used by containers to split space.
type Axis uint8

const (
	// AxisHorizontal lays out content left-to-right.
	AxisHorizontal Axis = 0
	// AxisVertical lays out content top-to-bottom.
	AxisVertical Axis = 1
)

// FitMode represents how content should fit within a [Block]'s content box.
type FitMode uint8

const (
	// FitExact performs no fitting and errors if content exceeds the content box.
	// This is the zero-value default.
	FitExact FitMode = iota
	// FitWrapClip wraps to the content box width, then clips vertically to fit.
	FitWrapClip
	// FitWrapStrict wraps to the content box width and errors if the wrapped
	// content exceeds the content box height.
	FitWrapStrict
	// FitClip clips content to the content box in both dimensions.
	FitClip
	// FitOverflow allows content to overflow (lipgloss default behavior).
	FitOverflow
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

// GetExtent implements the [Renderable] interface.
func (e ExtentConstraint) GetExtent() ExtentConstraint {
	return e
}

// RenderInfo describes the allocated space for a [Block] render pass.
type RenderInfo struct {
	Width, Height               int            // Total allocated size
	ContentWidth, ContentHeight int            // Inner content box size
	FrameWidth, FrameHeight     int            // Total frame size (padding + border + margin)
	Fit                         FitMode        // Fit mode for content
}
