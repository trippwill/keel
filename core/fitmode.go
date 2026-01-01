//go:generate stringer -type=FitMode -trimprefix=Fit
package core

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
