package keel

import (
	"testing"

	gloss "github.com/charmbracelet/lipgloss"
)

func TestRendererSetters(t *testing.T) {
	renderer := NewRenderer[string](Row(FlexUnit()), nil, nil)
	style := gloss.NewStyle().Bold(true)
	renderer.SetStyleProvider(func(id string) *gloss.Style { return &style })
	frame := Exact(Fixed(1), "a")
	if styleFor(renderer, frame) == nil {
		t.Fatalf("expected style")
	}
	renderer.SetContentProvider(func(id string, info FrameInfo) (string, error) { return "ok", nil })
	out, err := contentFor(renderer, "a", FrameInfo{})
	if err != nil || out != "ok" {
		t.Fatalf("expected content provider to run")
	}
}

func TestRendererSetConfigInvalidates(t *testing.T) {
	renderer := NewRenderer(Row(FlexUnit(), Exact(Fixed(1), "a")), nil, func(id string, info FrameInfo) (string, error) { return "x", nil })
	if _, err := renderer.Render(Size{Width: 1, Height: 1}); err != nil {
		t.Fatalf("unexpected render error: %v", err)
	}
	if !renderer.hasLayout {
		t.Fatalf("expected cached layout")
	}
	renderer.SetConfig(nil)
	if renderer.hasLayout {
		t.Fatalf("expected layout invalidated")
	}
	if renderer.config == nil {
		t.Fatalf("expected config allocated")
	}
}

func TestRendererDebugUsesDebugProvider(t *testing.T) {
	renderer := NewRenderer[string](Exact(Fixed(1), "a"), nil, nil)
	renderer.Config().SetDebug(true)
	out, err := renderer.Render(Size{Width: 1, Height: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Fatalf("expected debug output")
	}
}
