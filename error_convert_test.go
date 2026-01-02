package keel

import (
	"errors"
	"testing"

	"github.com/trippwill/keel/core"
)

func TestConvertErrorNil(t *testing.T) {
	if convertError(nil) != nil {
		t.Fatalf("expected nil")
	}
}

func TestConvertErrorPassesThroughKeelErrors(t *testing.T) {
	specErr := &SpecError{Kind: SpecKindAxis, Index: -1, Reason: "bad"}
	if convertError(specErr) != specErr {
		t.Fatalf("expected SpecError passthrough")
	}
	tooSmall := &ExtentTooSmallError{Axis: "Horizontal"}
	if convertError(tooSmall) != tooSmall {
		t.Fatalf("expected ExtentTooSmallError passthrough")
	}
}

func TestConvertErrorCoreExtentTooSmall(t *testing.T) {
	coreErr := &core.ExtentTooSmallError{
		Axis:   core.AxisVertical,
		Need:   10,
		Have:   4,
		Source: "vertical split",
		Reason: "allocation",
	}
	out := convertError(coreErr)
	tooSmall, ok := out.(*ExtentTooSmallError)
	if !ok {
		t.Fatalf("expected ExtentTooSmallError, got %T", out)
	}
	if tooSmall.Axis != "Vertical" || tooSmall.Need != 10 || tooSmall.Have != 4 || tooSmall.Source != "vertical split" || tooSmall.Reason != "allocation" {
		t.Fatalf("unexpected fields: %+v", tooSmall)
	}
	if !errors.Is(out, ErrExtentTooSmall) {
		t.Fatalf("expected ErrExtentTooSmall")
	}
}

func TestConvertErrorCoreSlotError(t *testing.T) {
	coreErr := &core.SlotError{Index: 2, Reason: core.ErrNilSlot}
	out := convertError(coreErr)
	specErr, ok := out.(*SpecError)
	if !ok {
		t.Fatalf("expected SpecError, got %T", out)
	}
	if specErr.Kind != SpecKindSlot || specErr.Index != 2 {
		t.Fatalf("unexpected SpecError: %+v", specErr)
	}
}

func TestConvertErrorCoreExtentError(t *testing.T) {
	coreErr := &core.ExtentError{Index: 1, Reason: core.ErrInvalidExtentUnits}
	out := convertError(coreErr)
	specErr, ok := out.(*SpecError)
	if !ok {
		t.Fatalf("expected SpecError, got %T", out)
	}
	if specErr.Kind != SpecKindExtent || specErr.Index != 1 {
		t.Fatalf("unexpected SpecError: %+v", specErr)
	}
}

func TestConvertErrorConfigReason(t *testing.T) {
	coreErr := &core.ConfigError{Reason: core.ErrInvalidAxis}
	out := convertError(coreErr)
	specErr, ok := out.(*SpecError)
	if !ok {
		t.Fatalf("expected SpecError, got %T", out)
	}
	if specErr.Kind != SpecKindAxis {
		t.Fatalf("expected axis kind, got %q", specErr.Kind)
	}
}

func TestConvertErrorConfigNilReason(t *testing.T) {
	coreErr := &core.ConfigError{}
	out := convertError(coreErr)
	specErr, ok := out.(*SpecError)
	if !ok {
		t.Fatalf("expected SpecError, got %T", out)
	}
	if specErr.Kind != SpecKindConfig {
		t.Fatalf("expected config kind, got %q", specErr.Kind)
	}
}

func TestConvertErrorFallbacks(t *testing.T) {
	cases := []struct {
		name string
		err  error
		kind string
	}{
		{"unknown spec", core.ErrUnknownSpec, SpecKindSpec},
		{"invalid axis", core.ErrInvalidAxis, SpecKindAxis},
		{"nil slot", core.ErrNilSlot, SpecKindSlot},
		{"invalid total", core.ErrInvalidTotal, SpecKindExtent},
		{"invalid extent kind", core.ErrInvalidExtentKind, SpecKindExtent},
		{"invalid extent units", core.ErrInvalidExtentUnits, SpecKindExtent},
		{"invalid extent min cells", core.ErrInvalidExtentMinCells, SpecKindExtent},
		{"invalid extent max cells", core.ErrInvalidExtentMaxCells, SpecKindExtent},
		{"invalid extent min", core.ErrInvalidExtentMin, SpecKindExtent},
		{"invalid extent max", core.ErrInvalidExtentMax, SpecKindExtent},
		{"config invalid", core.ErrConfigurationInvalid, SpecKindConfig},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := convertError(tc.err)
			specErr, ok := out.(*SpecError)
			if !ok {
				t.Fatalf("expected SpecError, got %T", out)
			}
			if specErr.Kind != tc.kind {
				t.Fatalf("expected kind %q, got %q", tc.kind, specErr.Kind)
			}
			if !errors.Is(out, ErrConfigurationInvalid) {
				t.Fatalf("expected ErrConfigurationInvalid")
			}
		})
	}
}

func TestConvertErrorUnknownPassesThrough(t *testing.T) {
	unknown := errors.New("unknown")
	if convertError(unknown) != unknown {
		t.Fatalf("expected passthrough")
	}
}
