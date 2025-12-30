package keel

// PanelSpec is the frame spec that displays content for a specific ID.
type PanelSpec[KID KeelID] struct {
	ExtentConstraint // Size specification for the panel
	id               KID
	fit              FitMode
}

// ID implements [FrameSpec].
func (p PanelSpec[KID]) ID() KID {
	return p.id
}

// Fit implements [FrameSpec].
func (p PanelSpec[KID]) Fit() FitMode {
	return p.fit
}

var _ FrameSpec[string] = PanelSpec[string]{}

// Panel creates a new PanelSpec with the given extent and ID.
// The extent describes total allocation along the stack axis.
// Content fitting defaults to [FitExact].
func Panel[KID KeelID](extent ExtentConstraint, id KID) PanelSpec[KID] {
	return PanelFit(extent, FitExact, id)
}

// PanelClip creates a new PanelSpec with the given extent and ID.
// The extent describes total allocation along the stack axis.
// Content fitting is set to [FitClip].
func PanelClip[KID KeelID](extent ExtentConstraint, id KID) PanelSpec[KID] {
	return PanelFit(extent, FitClip, id)
}

// PanelWrap creates a new PanelSpec with the given extent and ID.
// The extent describes total allocation along the stack axis.
// Content fitting is set to [FitWrapClip].
func PanelWrap[KID KeelID](extent ExtentConstraint, id KID) PanelSpec[KID] {
	return PanelFit(extent, FitWrapClip, id)
}

// PanelWrapStrict creates a new PanelSpec with the given extent and ID.
// The extent describes total allocation along the stack axis.
// Content fitting is set to [FitWrapStrict].
func PanelWrapStrict[KID KeelID](extent ExtentConstraint, id KID) PanelSpec[KID] {
	return PanelFit(extent, FitWrapStrict, id)
}

// PanelOverflow creates a new PanelSpec with the given extent and ID.
// The extent describes total allocation along the stack axis.
// Content fitting is set to [FitOverflow].
func PanelOverflow[KID KeelID](extent ExtentConstraint, id KID) PanelSpec[KID] {
	return PanelFit(extent, FitOverflow, id)
}

// PanelFit creates a new PanelSpec with the given extent, content fit mode, and ID.
// The extent describes total allocation along the stack axis.
func PanelFit[KID KeelID](extent ExtentConstraint, fit FitMode, id KID) PanelSpec[KID] {
	return PanelSpec[KID]{
		ExtentConstraint: extent,
		id:               id,
		fit:              fit,
	}
}
