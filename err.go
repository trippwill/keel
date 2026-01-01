package keel

import (
	"errors"
	"fmt"
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
	// ErrEmptySlots indicates a stack with no slots.
	ErrEmptySlots = errors.New("empty slots")
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
