package keel

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrRendererMissing indicates a missing renderer.
	ErrRendererMissing = errors.New("renderer missing")
	// ErrSpecMissing indicates a missing layout spec.
	ErrSpecMissing = errors.New("spec missing")
	// ErrContentProviderMissing indicates a missing content provider.
	ErrContentProviderMissing = errors.New("content provider missing")
	// ErrUnknownFrameID indicates a content/style request for an unknown ID.
	ErrUnknownFrameID = errors.New("unknown frame id")
)

// ContentProviderMissingError indicates a missing content provider for a frame ID.
// It wraps ErrContentProviderMissing for errors.Is checks.
type ContentProviderMissingError struct {
	ID any
}

func (e *ContentProviderMissingError) Error() string {
	return fmt.Sprintf("%s: %v", ErrContentProviderMissing, e.ID)
}

func (e *ContentProviderMissingError) Unwrap() error {
	return ErrContentProviderMissing
}

// UnknownFrameIDError indicates a request for an unknown frame ID.
// It wraps ErrUnknownFrameID for errors.Is checks.
type UnknownFrameIDError struct {
	ID any
}

func (e *UnknownFrameIDError) Error() string {
	return fmt.Sprintf("%s: %v", ErrUnknownFrameID, e.ID)
}

func (e *UnknownFrameIDError) Unwrap() error {
	return ErrUnknownFrameID
}

// ExtentTooSmallError includes context about which allocation failed.
// It wraps ErrExtentTooSmall for errors.Is checks.
type ExtentTooSmallError struct {
	Axis       string
	Need, Have int
	Source     string
	Reason     string
}

func (e *ExtentTooSmallError) Error() string {
	source := ""
	if e.Source != "" {
		source = " for " + e.Source
	}
	reason := ""
	if e.Reason != "" {
		reason = " (" + e.Reason + ")"
	}
	axis := e.Axis
	if axis == "" {
		axis = "unknown"
	}
	return fmt.Sprintf(
		"extent too small on %s axis%s%s: need %d, have %d",
		axis,
		source,
		reason,
		e.Need,
		e.Have,
	)
}

func (e *ExtentTooSmallError) Unwrap() error {
	return ErrExtentTooSmall
}

// SpecError describes a configuration issue with a spec, slot, or extent.
// It wraps ErrConfigurationInvalid for errors.Is checks.
type SpecError struct {
	Kind   string
	Index  int
	Reason string
}

const (
	// SpecKindSpec identifies a configuration issue with a spec value.
	SpecKindSpec = "spec"
	// SpecKindAxis identifies a configuration issue with a stack axis.
	SpecKindAxis = "axis"
	// SpecKindSlot identifies a configuration issue with a slot entry.
	SpecKindSlot = "slot"
	// SpecKindExtent identifies a configuration issue with an extent constraint.
	SpecKindExtent = "extent"
	// SpecKindConfig identifies an otherwise unspecified configuration issue.
	SpecKindConfig = "config"
)

func (e *SpecError) Error() string {
	parts := make([]string, 0, 2)
	if e.Kind != "" {
		if e.Index >= 0 {
			parts = append(parts, fmt.Sprintf("%s %d", e.Kind, e.Index))
		} else {
			parts = append(parts, e.Kind)
		}
	}
	if e.Reason != "" {
		parts = append(parts, e.Reason)
	}
	if len(parts) == 0 {
		return ErrConfigurationInvalid.Error()
	}
	return fmt.Sprintf("%s: %s", ErrConfigurationInvalid, strings.Join(parts, ": "))
}

func (e *SpecError) Unwrap() error {
	return ErrConfigurationInvalid
}
