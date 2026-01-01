package core

import (
	"errors"
	"testing"
)

func TestConfigErrorIs(t *testing.T) {
	cases := []struct {
		name   string
		err    error
		target error
		want   bool
	}{
		{
			name:   "matches configuration invalid",
			err:    &ConfigError{Reason: ErrInvalidAxis},
			target: ErrConfigurationInvalid,
			want:   true,
		},
		{
			name:   "matches underlying reason",
			err:    &ConfigError{Reason: ErrInvalidAxis},
			target: ErrInvalidAxis,
			want:   true,
		},
		{
			name:   "nil reason matches only configuration invalid",
			err:    &ConfigError{},
			target: ErrInvalidAxis,
			want:   false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := errors.Is(tc.err, tc.target); got != tc.want {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
}

func TestExtentErrorIs(t *testing.T) {
	cases := []struct {
		name   string
		err    error
		target error
		want   bool
	}{
		{
			name:   "matches configuration invalid",
			err:    &ExtentError{Index: 2, Reason: ErrInvalidExtentUnits},
			target: ErrConfigurationInvalid,
			want:   true,
		},
		{
			name:   "matches underlying reason",
			err:    &ExtentError{Index: 2, Reason: ErrInvalidExtentUnits},
			target: ErrInvalidExtentUnits,
			want:   true,
		},
		{
			name:   "nil reason matches only configuration invalid",
			err:    &ExtentError{Index: 2},
			target: ErrInvalidExtentUnits,
			want:   false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := errors.Is(tc.err, tc.target); got != tc.want {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
}

func TestSlotErrorIs(t *testing.T) {
	cases := []struct {
		name   string
		err    error
		target error
		want   bool
	}{
		{
			name:   "matches configuration invalid",
			err:    &SlotError{Index: 1, Reason: ErrNilSlot},
			target: ErrConfigurationInvalid,
			want:   true,
		},
		{
			name:   "matches underlying reason",
			err:    &SlotError{Index: 1, Reason: ErrNilSlot},
			target: ErrNilSlot,
			want:   true,
		},
		{
			name:   "nil reason matches only configuration invalid",
			err:    &SlotError{Index: 1},
			target: ErrNilSlot,
			want:   false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := errors.Is(tc.err, tc.target); got != tc.want {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
}

func TestExtentTooSmallErrorIs(t *testing.T) {
	err := &ExtentTooSmallError{
		Axis: AxisHorizontal,
		Need: 10,
		Have: 5,
	}
	if !errors.Is(err, ErrExtentTooSmall) {
		t.Fatalf("expected ErrExtentTooSmall")
	}
}
