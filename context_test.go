package keel

import (
	"testing"

	gloss "github.com/charmbracelet/lipgloss"
)

func TestContextWithSize(t *testing.T) {
	orig := Context[string]{Width: 10, Height: 5}
	next := orig.WithSize(3, 4)

	if next.Width != 3 || next.Height != 4 {
		t.Fatalf("expected 3x4, got %dx%d", next.Width, next.Height)
	}
	if orig.Width != 10 || orig.Height != 5 {
		t.Fatalf("expected original unchanged, got %dx%d", orig.Width, orig.Height)
	}
}

func TestContextWithWidth(t *testing.T) {
	orig := Context[string]{Width: 10, Height: 5}
	next := orig.WithWidth(7)

	if next.Width != 7 || next.Height != 5 {
		t.Fatalf("expected 7x5, got %dx%d", next.Width, next.Height)
	}
	if orig.Width != 10 || orig.Height != 5 {
		t.Fatalf("expected original unchanged, got %dx%d", orig.Width, orig.Height)
	}
}

func TestContextWithHeight(t *testing.T) {
	orig := Context[string]{Width: 10, Height: 5}
	next := orig.WithHeight(9)

	if next.Width != 10 || next.Height != 9 {
		t.Fatalf("expected 10x9, got %dx%d", next.Width, next.Height)
	}
	if orig.Width != 10 || orig.Height != 5 {
		t.Fatalf("expected original unchanged, got %dx%d", orig.Width, orig.Height)
	}
}

func TestContextWithProviders(t *testing.T) {
	style := func(id string) *gloss.Style { return nil }
	content := func(id string, info RenderInfo) (string, error) { return "ok", nil }
	logger := func(event LogEvent, path, msg string) {}

	orig := Context[string]{Width: 1, Height: 2}
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
	if withLogger.Width != 1 || withLogger.Height != 2 {
		t.Fatalf("expected dimensions preserved, got %dx%d", withLogger.Width, withLogger.Height)
	}
}
