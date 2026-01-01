package keel

import (
	"testing"

	gloss "github.com/charmbracelet/lipgloss"
	"github.com/trippwill/keel/logging"
)

func TestRenderConfigAndRenderer(t *testing.T) {
	style := func(id string) *gloss.Style { return nil }
	content := func(id string, info FrameInfo) (string, error) { return "ok", nil }
	logger := func(event logging.LogEvent, path, msg string) {}

	config := NewConfig()
	config.SetLogger(logger)
	config.SetDebug(true)

	renderer := NewRendererWithConfig(config, Row(FlexUnit()), style, content)
	if renderer.Config().Logger() == nil {
		t.Fatalf("expected logger")
	}
	if !renderer.Config().Debug() {
		t.Fatalf("expected debug enabled")
	}
}
