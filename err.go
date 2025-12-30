package keel

import (
	"errors"
	"fmt"
)

var (
	// ErrConfigurationInvalid indicates an invalid layout configuration.
	ErrConfigurationInvalid = errors.New("configuration invalid")
	// ErrContentProviderMissing indicates a missing content provider.
	ErrContentProviderMissing = errors.New("content provider missing")
	// ErrUnknownFrameID indicates a content/style request for an unknown ID.
	ErrUnknownFrameID = errors.New("unknown frame id")
	// ErrInvalidAxis indicates an invalid axis value.
	ErrInvalidAxis = errors.New("invalid axis")
	// ErrEmptySlots indicates a stack with no slots.
	ErrEmptySlots = errors.New("empty slots")
	// ErrInvalidTotal indicates an invalid total allocation.
	ErrInvalidTotal = errors.New("invalid total")
	// ErrEmptyExtents indicates a missing set of extents.
	ErrEmptyExtents = errors.New("empty extents")
	// ErrInvalidExtentKind indicates an invalid extent kind.
	ErrInvalidExtentKind = errors.New("invalid extent kind")
	// ErrInvalidExtentUnits indicates invalid extent units.
	ErrInvalidExtentUnits = errors.New("invalid extent units")
	// ErrInvalidExtentMinCells indicates invalid minimum cells for an extent.
	ErrInvalidExtentMinCells = errors.New("invalid extent min cells")
	// ErrInvalidExtentMaxCells indicates invalid maximum cells for an extent.
	ErrInvalidExtentMaxCells = errors.New("invalid extent max cells")
	// ErrInvalidExtentMin indicates invalid minimum requirements for an extent.
	ErrInvalidExtentMin = errors.New("invalid extent min")
	// ErrInvalidExtentMax indicates invalid maximum requirements for an extent.
	ErrInvalidExtentMax = errors.New("invalid extent max")
	// ErrNilSlot indicates a nil slot entry.
	ErrNilSlot = errors.New("nil slot")
	// ErrUnknownSpec indicates a Spec with an unsupported type.
	ErrUnknownSpec = errors.New("unknown spec")
	// ErrExtentTooSmall indicates insufficient extent for an allocation.
	ErrExtentTooSmall = errors.New("extent too small")
)

// ExtentTooSmallError includes context about which allocation failed.
// It wraps ErrExtentTooSmall for errors.Is checks.
type ExtentTooSmallError struct {
	Axis       Axis
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
	return fmt.Sprintf(
		"extent too small on %s axis%s%s: need %d, have %d",
		e.Axis,
		source,
		reason,
		e.Need,
		e.Have,
	)
}

func (e *ExtentTooSmallError) Unwrap() error {
	return ErrExtentTooSmall
}

// ConfigError wraps a configuration issue with a specific reason.
// It unwraps to the underlying reason and matches ErrConfigurationInvalid.
type ConfigError struct {
	Reason error
}

func (e *ConfigError) Error() string {
	if e.Reason == nil {
		return ErrConfigurationInvalid.Error()
	}
	return fmt.Sprintf("%s: %s", ErrConfigurationInvalid, e.Reason)
}

func (e *ConfigError) Unwrap() error {
	return e.Reason
}

func (e *ConfigError) Is(target error) bool {
	if target == ErrConfigurationInvalid {
		return true
	}
	if e.Reason == nil {
		return false
	}
	return errors.Is(e.Reason, target)
}

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

// ExtentError describes a validation issue for a specific extent.
// It wraps ErrConfigurationInvalid and the underlying reason.
type ExtentError struct {
	Index  int
	Reason error
}

func (e *ExtentError) Error() string {
	if e.Reason == nil {
		return fmt.Sprintf("%s: extent %d", ErrConfigurationInvalid, e.Index)
	}
	return fmt.Sprintf("%s: extent %d: %s", ErrConfigurationInvalid, e.Index, e.Reason)
}

func (e *ExtentError) Unwrap() error {
	return e.Reason
}

func (e *ExtentError) Is(target error) bool {
	if target == ErrConfigurationInvalid {
		return true
	}
	if e.Reason == nil {
		return false
	}
	return errors.Is(e.Reason, target)
}

// SlotError describes a validation issue for a specific slot.
// It wraps ErrConfigurationInvalid and the underlying reason.
type SlotError struct {
	Index  int
	Reason error
}

func (e *SlotError) Error() string {
	if e.Reason == nil {
		return fmt.Sprintf("%s: slot %d", ErrConfigurationInvalid, e.Index)
	}
	return fmt.Sprintf("%s: slot %d: %s", ErrConfigurationInvalid, e.Index, e.Reason)
}

func (e *SlotError) Unwrap() error {
	return e.Reason
}

func (e *SlotError) Is(target error) bool {
	if target == ErrConfigurationInvalid {
		return true
	}
	if e.Reason == nil {
		return false
	}
	return errors.Is(e.Reason, target)
}
