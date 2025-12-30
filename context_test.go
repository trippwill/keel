package keel

import (
	"testing"

	gloss "github.com/charmbracelet/lipgloss"
)

func TestContextWithProviders(t *testing.T) {
	style := func(id string) *gloss.Style { return nil }
	content := func(id string, info RenderInfo) (string, error) { return "ok", nil }
	logger := func(event LogEvent, path, msg string) {}

	orig := Context[string]{}
	withStyle := orig.WithStyleProvider(style)
	withContent := withStyle.WithContentProvider(content)
	withLogger := withContent.WithLogger(logger)

	if withStyle.StyleProvider == nil {
		t.Fatalf("expected style provider")
	}
	if withContent.ContentProvider == nil {
		t.Fatalf("expected content provider")
	}
	if orig.StyleProvider != nil || orig.ContentProvider != nil {
		t.Fatalf("expected original providers unchanged")
	}
	if withLogger.Logger == nil {
		t.Fatalf("expected logger")
	}
}
