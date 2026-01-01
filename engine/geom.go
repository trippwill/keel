//go:generate stringer -type=Axis -trimprefix=Axis
//go:generate stringer -type=ExtentKind -trimprefix=Extent
//go:generate stringer -type=FitMode -trimprefix=Fit
package engine

// Axis represents a layout axis used by stacks to split space.
type Axis uint8

const (
	// AxisHorizontal lays out content left-to-right.
	AxisHorizontal Axis = 0
	// AxisVertical lays out content top-to-bottom.
	AxisVertical Axis = 1
)

// Size describes a width/height pair in cells.
type Size struct {
	Width, Height int
}

// FitMode represents how content should fit within a [FrameSpec]'s content box.
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

// FrameInfo describes the allocated space for a [FrameSpec] render pass.
type FrameInfo struct {
	Width, Height               int     // Total allocated size
	ContentWidth, ContentHeight int     // Inner content box size
	FrameWidth, FrameHeight     int     // Total frame size (padding + border + margin)
	Fit                         FitMode // Fit mode for content
}
