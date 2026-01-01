package keel

import "github.com/trippwill/keel/core"

var (
	ErrExtentTooSmall       = core.ErrExtentTooSmall
	ErrConfigurationInvalid = core.ErrConfigurationInvalid
)

type (
	KeelID                = core.KeelID
	Size                  = core.Size
	Spec                  = core.Spec
	FrameSpec[KID KeelID] = core.FrameSpec[KID]
	StackSpec             = core.StackSpec
	FrameInfo             = core.FrameInfo
)
