package engine

import "github.com/trippwill/keel/core"

// PanelSpec is the frame spec that displays content for a specific ID.
type PanelSpec[KID core.KeelID] struct {
	core.ExtentConstraint // Size specification for the panel
	id                    KID
	fit                   core.FitMode
}

// NewPanelSpec creates a new PanelSpec with the given extent, content fit mode, and ID.
//
// Arguments:
//
//	extent: describes total allocation along the stack axis.
//	fit: content fitting mode
//	id: the [KeelID] to assign to the panel
//
// Returns:
//   - A new [PanelSpec] instance configured with the provided arguments.
func NewPanelSpec[KID core.KeelID](extent core.ExtentConstraint, fit core.FitMode, id KID) PanelSpec[KID] {
	return PanelSpec[KID]{
		ExtentConstraint: extent,
		id:               id,
		fit:              fit,
	}
}

// ID implements [FrameSpec].
func (p PanelSpec[KID]) ID() KID {
	return p.id
}

// Fit implements [FrameSpec].
func (p PanelSpec[KID]) Fit() core.FitMode {
	return p.fit
}

var _ core.FrameSpec[string] = PanelSpec[string]{}
