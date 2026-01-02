package keel

import (
	"errors"

	"github.com/trippwill/keel/core"
)

func convertError(err error) error {
	if err == nil {
		return nil
	}
	var specErr *SpecError
	if errors.As(err, &specErr) {
		return err
	}
	var tooSmall *ExtentTooSmallError
	if errors.As(err, &tooSmall) {
		return err
	}
	var coreTooSmall *core.ExtentTooSmallError
	if errors.As(err, &coreTooSmall) {
		return &ExtentTooSmallError{
			Axis:   coreTooSmall.Axis.String(),
			Need:   coreTooSmall.Need,
			Have:   coreTooSmall.Have,
			Source: coreTooSmall.Source,
			Reason: coreTooSmall.Reason,
		}
	}
	var slotErr *core.SlotError
	if errors.As(err, &slotErr) {
		return newSpecError(SpecKindSlot, slotErr.Index, slotErr.Reason)
	}
	var extentErr *core.ExtentError
	if errors.As(err, &extentErr) {
		return newSpecError(SpecKindExtent, extentErr.Index, extentErr.Reason)
	}
	var configErr *core.ConfigError
	if errors.As(err, &configErr) {
		return newSpecError(kindForReason(configErr.Reason), -1, configErr.Reason)
	}

	switch {
	case errors.Is(err, core.ErrUnknownSpec):
		return newSpecError(SpecKindSpec, -1, core.ErrUnknownSpec)
	case errors.Is(err, core.ErrInvalidAxis):
		return newSpecError(SpecKindAxis, -1, core.ErrInvalidAxis)
	case errors.Is(err, core.ErrNilSlot):
		return newSpecError(SpecKindSlot, -1, core.ErrNilSlot)
	case errors.Is(err, core.ErrInvalidTotal):
		return newSpecError(SpecKindExtent, -1, core.ErrInvalidTotal)
	case errors.Is(err, core.ErrInvalidExtentKind):
		return newSpecError(SpecKindExtent, -1, core.ErrInvalidExtentKind)
	case errors.Is(err, core.ErrInvalidExtentUnits):
		return newSpecError(SpecKindExtent, -1, core.ErrInvalidExtentUnits)
	case errors.Is(err, core.ErrInvalidExtentMinCells):
		return newSpecError(SpecKindExtent, -1, core.ErrInvalidExtentMinCells)
	case errors.Is(err, core.ErrInvalidExtentMaxCells):
		return newSpecError(SpecKindExtent, -1, core.ErrInvalidExtentMaxCells)
	case errors.Is(err, core.ErrInvalidExtentMin):
		return newSpecError(SpecKindExtent, -1, core.ErrInvalidExtentMin)
	case errors.Is(err, core.ErrInvalidExtentMax):
		return newSpecError(SpecKindExtent, -1, core.ErrInvalidExtentMax)
	case errors.Is(err, core.ErrConfigurationInvalid):
		return &SpecError{Kind: SpecKindConfig, Index: -1}
	default:
		return err
	}
}

func newSpecError(kind string, index int, reason error) *SpecError {
	spec := &SpecError{Kind: kind, Index: index}
	if reason != nil {
		spec.Reason = reason.Error()
	}
	return spec
}

func kindForReason(reason error) string {
	if reason == nil {
		return SpecKindConfig
	}
	switch {
	case errors.Is(reason, core.ErrInvalidAxis):
		return SpecKindAxis
	case errors.Is(reason, core.ErrUnknownSpec):
		return SpecKindSpec
	case errors.Is(reason, core.ErrNilSlot):
		return SpecKindSlot
	case errors.Is(reason, core.ErrInvalidTotal),
		errors.Is(reason, core.ErrInvalidExtentKind),
		errors.Is(reason, core.ErrInvalidExtentUnits),
		errors.Is(reason, core.ErrInvalidExtentMinCells),
		errors.Is(reason, core.ErrInvalidExtentMaxCells),
		errors.Is(reason, core.ErrInvalidExtentMin),
		errors.Is(reason, core.ErrInvalidExtentMax):
		return SpecKindExtent
	default:
		return SpecKindConfig
	}
}
