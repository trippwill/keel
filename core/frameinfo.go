package core

// FrameInfo describes the allocated space for a [FrameSpec] render pass.
type FrameInfo struct {
	Width, Height               int     // Total allocated size
	ContentWidth, ContentHeight int     // Inner content box size
	FrameWidth, FrameHeight     int     // Total frame size (padding + border + margin)
	Fit                         FitMode // Fit mode for content
}
