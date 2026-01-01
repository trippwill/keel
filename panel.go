package keel

import "github.com/trippwill/keel/engine"

// PanelSpec is the frame spec that displays content for a specific ID.
type PanelSpec[KID KeelID] struct {
	engine.ExtentConstraint // Size specification for the panel
	id                      KID
	fit                     engine.FitMode
}

// ID implements [FrameSpec].
func (p PanelSpec[KID]) ID() KID {
	return p.id
}

// Fit implements [FrameSpec].
func (p PanelSpec[KID]) Fit() engine.FitMode {
	return p.fit
}

var _ engine.FrameSpec[string] = PanelSpec[string]{}

// Panel creates a new PanelSpec with the given extent and ID.
// The extent describes total allocation along the stack axis.
// Content fitting defaults to [FitExact].
func Panel[KID KeelID](extent engine.ExtentConstraint, id KID) PanelSpec[KID] {
	return PanelFit(extent, engine.FitExact, id)
}

// PanelClip creates a new PanelSpec with the given extent and ID.
// The extent describes total allocation along the stack axis.
// Content fitting is set to [FitClip].
func PanelClip[KID KeelID](extent engine.ExtentConstraint, id KID) PanelSpec[KID] {
	return PanelFit(extent, engine.FitClip, id)
}

// PanelWrap creates a new PanelSpec with the given extent and ID.
// The extent describes total allocation along the stack axis.
// Content fitting is set to [FitWrapClip].
func PanelWrap[KID KeelID](extent engine.ExtentConstraint, id KID) PanelSpec[KID] {
	return PanelFit(extent, engine.FitWrapClip, id)
}

// PanelWrapStrict creates a new PanelSpec with the given extent and ID.
// The extent describes total allocation along the stack axis.
// Content fitting is set to [FitWrapStrict].
func PanelWrapStrict[KID KeelID](extent engine.ExtentConstraint, id KID) PanelSpec[KID] {
	return PanelFit(extent, engine.FitWrapStrict, id)
}

// PanelOverflow creates a new PanelSpec with the given extent and ID.
// The extent describes total allocation along the stack axis.
// Content fitting is set to [FitOverflow].
func PanelOverflow[KID KeelID](extent engine.ExtentConstraint, id KID) PanelSpec[KID] {
	return PanelFit(extent, engine.FitOverflow, id)
}

// PanelFit creates a new PanelSpec with the given extent, content fit mode, and ID.
//
// Arguments:
//
//	extent: describes total allocation along the stack axis.
//	fit: content fitting mode
//	id: the [KeelID] to assign to the panel
//
// Returns:
//   - A new [PanelSpec] instance configured with the provided arguments.
func PanelFit[KID KeelID](extent engine.ExtentConstraint, fit engine.FitMode, id KID) PanelSpec[KID] {
	return PanelSpec[KID]{
		ExtentConstraint: extent,
		id:               id,
		fit:              fit,
	}
}
