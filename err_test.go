package keel

import (
	"errors"
	"strings"
	"testing"
)

func TestExtentTooSmallErrorFormatting(t *testing.T) {
	err := &ExtentTooSmallError{
		Axis:   "Horizontal",
		Need:   10,
		Have:   5,
		Source: "frame header",
		Reason: "content",
	}
	want := "extent too small on Horizontal axis for frame header (content): need 10, have 5"
	if err.Error() != want {
		t.Fatalf("expected %q, got %q", want, err.Error())
	}
	if !errors.Is(err, ErrExtentTooSmall) {
		t.Fatalf("expected ErrExtentTooSmall")
	}
}

func TestExtentTooSmallErrorUnknownAxis(t *testing.T) {
	err := &ExtentTooSmallError{Need: 1, Have: 0}
	if err.Error() == "" || !strings.HasPrefix(err.Error(), "extent too small on") {
		t.Fatalf("unexpected error string: %q", err.Error())
	}
}

func TestSpecErrorFormatting(t *testing.T) {
	cases := []struct {
		name string
		err  *SpecError
		want string
	}{
		{"kind index reason", &SpecError{Kind: SpecKindSlot, Index: 2, Reason: "nil slot"}, "configuration invalid: slot 2: nil slot"},
		{"kind only", &SpecError{Kind: SpecKindAxis, Index: -1}, "configuration invalid: axis"},
		{"reason only", &SpecError{Reason: "bad"}, "configuration invalid: bad"},
		{"empty", &SpecError{}, "configuration invalid"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err.Error() != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, tc.err.Error())
			}
			if !errors.Is(tc.err, ErrConfigurationInvalid) {
				t.Fatalf("expected ErrConfigurationInvalid")
			}
		})
	}
}

func TestContentProviderMissingError(t *testing.T) {
	err := &ContentProviderMissingError{ID: "a"}
	if !errors.Is(err, ErrContentProviderMissing) {
		t.Fatalf("expected ErrContentProviderMissing")
	}
}

func TestUnknownFrameIDError(t *testing.T) {
	err := &UnknownFrameIDError{ID: "a"}
	if !errors.Is(err, ErrUnknownFrameID) {
		t.Fatalf("expected ErrUnknownFrameID")
	}
}
