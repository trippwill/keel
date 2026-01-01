package keel

import (
	"github.com/trippwill/keel/core"
	"github.com/trippwill/keel/engine"
)

// Exact creates a new PanelSpec with the given extent and ID.
// The extent describes total allocation along the stack axis.
// Content fitting defaults to [core.FitExact].
func Exact[KID KeelID](extent ExtentConstraint, id KID) FrameSpec[KID] {
	return engine.NewPanelSpec(extent, core.FitExact, id)
}

// Clip creates a new PanelSpec with the given extent and ID.
// The extent describes total allocation along the stack axis.
// Content fitting is set to [core.FitClip].
func Clip[KID KeelID](extent ExtentConstraint, id KID) FrameSpec[KID] {
	return engine.NewPanelSpec(extent, core.FitClip, id)
}

// Wrap creates a new PanelSpec with the given extent and ID.
// The extent describes total allocation along the stack axis.
// Content fitting is set to [core.FitWrapClip].
func Wrap[KID KeelID](extent ExtentConstraint, id KID) FrameSpec[KID] {
	return engine.NewPanelSpec(extent, core.FitWrapClip, id)
}

// WrapStrict creates a new PanelSpec with the given extent and ID.
// The extent describes total allocation along the stack axis.
// Content fitting is set to [core.FitWrapStrict].
func WrapStrict[KID KeelID](extent ExtentConstraint, id KID) FrameSpec[KID] {
	return engine.NewPanelSpec(extent, core.FitWrapStrict, id)
}

// Overflow creates a new PanelSpec with the given extent and ID.
// The extent describes total allocation along the stack axis.
// Content fitting is set to [core.FitOverflow].
func Overflow[KID KeelID](extent ExtentConstraint, id KID) FrameSpec[KID] {
	return engine.NewPanelSpec(extent, core.FitOverflow, id)
}
