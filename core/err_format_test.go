package core

import (
	"errors"
	"testing"
)

func TestExtentTooSmallErrorFormat(t *testing.T) {
	err := &ExtentTooSmallError{
		Axis:   AxisVertical,
		Need:   10,
		Have:   4,
		Source: "vertical split",
		Reason: "allocation",
	}
	want := "extent too small on Vertical axis for vertical split (allocation): need 10, have 4"
	if err.Error() != want {
		t.Fatalf("expected %q, got %q", want, err.Error())
	}
	if !errors.Is(err, ErrExtentTooSmall) {
		t.Fatalf("expected ErrExtentTooSmall")
	}
}

func TestConfigErrorFormat(t *testing.T) {
	err := &ConfigError{Reason: ErrInvalidAxis}
	want := "configuration invalid: invalid axis"
	if err.Error() != want {
		t.Fatalf("expected %q, got %q", want, err.Error())
	}

	err = &ConfigError{}
	if err.Error() != ErrConfigurationInvalid.Error() {
		t.Fatalf("expected configuration invalid")
	}
}

func TestExtentErrorFormat(t *testing.T) {
	err := &ExtentError{Index: 3, Reason: ErrInvalidExtentUnits}
	want := "configuration invalid: extent 3: invalid extent units"
	if err.Error() != want {
		t.Fatalf("expected %q, got %q", want, err.Error())
	}

	err = &ExtentError{Index: 2}
	want = "configuration invalid: extent 2"
	if err.Error() != want {
		t.Fatalf("expected %q, got %q", want, err.Error())
	}
}

func TestSlotErrorFormat(t *testing.T) {
	err := &SlotError{Index: 1, Reason: ErrNilSlot}
	want := "configuration invalid: slot 1: nil slot"
	if err.Error() != want {
		t.Fatalf("expected %q, got %q", want, err.Error())
	}

	err = &SlotError{Index: 4}
	want = "configuration invalid: slot 4"
	if err.Error() != want {
		t.Fatalf("expected %q, got %q", want, err.Error())
	}
}
