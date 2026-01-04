package keel

import (
	"io"
	"log/slog"
	"testing"

	gloss "github.com/charmbracelet/lipgloss"
)

func TestRenderConfigAndRenderer(t *testing.T) {
	style := func(id string) *gloss.Style { return nil }
	content := func(id string, info FrameInfo) (string, error) { return "ok", nil }
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))

	config := NewConfig()
	config.SetLogger(logger)
	config.SetDebug(true)

	renderer := NewRendererWithConfig(config, Row(FlexUnit()), style, content)
	if renderer.Config().Logger() != logger {
		t.Fatalf("expected logger")
	}
	if !renderer.Config().Debug() {
		t.Fatalf("expected debug enabled")
	}
}
