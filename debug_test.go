package keel

import (
	"testing"

	gloss "github.com/charmbracelet/lipgloss"
)

func TestDefaultDebugProviderFitsContent(t *testing.T) {
	provider := DefaultDebugProvider[string]
	info := RenderInfo[string]{
		ID:            "a",
		Width:         6,
		Height:        3,
		ContentWidth:  6,
		ContentHeight: 2,
		FrameWidth:    0,
		FrameHeight:   0,
		Clip:          ClipConstraint{},
	}

	got, err := provider(info)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	width, height := gloss.Size(got)
	if width > info.ContentWidth || height > info.ContentHeight {
		t.Fatalf("expected <= %dx%d, got %dx%d", info.ContentWidth, info.ContentHeight, width, height)
	}
}

func TestDefaultDebugProviderCompactSingleLine(t *testing.T) {
	provider := DefaultDebugProvider[string]
	info := RenderInfo[string]{
		ID:            "header",
		Width:         70,
		Height:        12,
		ContentWidth:  80,
		ContentHeight: 1,
		FrameWidth:    4,
		FrameHeight:   2,
	}

	got, err := provider(info)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "id:header|a:70x12|f:4x2|c:80x1"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestDefaultDebugProviderCompactTwoLines(t *testing.T) {
	provider := DefaultDebugProvider[string]
	info := RenderInfo[string]{
		ID:            "status",
		Width:         30,
		Height:        3,
		ContentWidth:  80,
		ContentHeight: 2,
		FrameWidth:    2,
		FrameHeight:   2,
	}

	got, err := provider(info)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "id:status\nid:status|a:30x3|f:2x2|c:80x2"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestDefaultDebugProviderExpandedLines(t *testing.T) {
	provider := DefaultDebugProvider[string]
	info := RenderInfo[string]{
		ID:            "nav",
		Width:         15,
		Height:        6,
		ContentWidth:  20,
		ContentHeight: 5,
		FrameWidth:    4,
		FrameHeight:   2,
		Clip:          ClipConstraint{Width: 10, Height: 0},
	}

	got, err := provider(info)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "id:nav\nalloc:15x6\nframe:4x2\ncontent:20x5\nclip:10x0"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
