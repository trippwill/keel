package core

import (
	"testing"

	gloss "github.com/charmbracelet/lipgloss"
)

func TestDebugProviderFitsContent(t *testing.T) {
	provider := DebugContentProvider[string]
	id := "a"
	info := FrameInfo{
		Width:         6,
		Height:        3,
		ContentWidth:  6,
		ContentHeight: 2,
		FrameWidth:    0,
		FrameHeight:   0,
		Fit:           FitExact,
	}

	got, err := provider(id, info)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	width, height := gloss.Size(got)
	if width > info.ContentWidth || height > info.ContentHeight {
		t.Fatalf("expected <= %dx%d, got %dx%d", info.ContentWidth, info.ContentHeight, width, height)
	}
}

func TestDebugProviderCompactSingleLine(t *testing.T) {
	provider := DebugContentProvider[string]
	id := "header"
	info := FrameInfo{
		Width:         70,
		Height:        12,
		ContentWidth:  80,
		ContentHeight: 1,
		FrameWidth:    4,
		FrameHeight:   2,
	}

	got, err := provider(id, info)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "id:header|a:70x12|f:4x2|c:80x1|ft:Exact"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestDebugProviderCompactTwoLines(t *testing.T) {
	provider := DebugContentProvider[string]
	id := "status"
	info := FrameInfo{
		Width:         30,
		Height:        3,
		ContentWidth:  80,
		ContentHeight: 2,
		FrameWidth:    2,
		FrameHeight:   2,
	}

	got, err := provider(id, info)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "id:status\nid:status|a:30x3|f:2x2|c:80x2|ft:Exact"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestDebugProviderExpandedLines(t *testing.T) {
	provider := DebugContentProvider[string]
	id := "nav"
	info := FrameInfo{
		Width:         15,
		Height:        6,
		ContentWidth:  20,
		ContentHeight: 5,
		FrameWidth:    4,
		FrameHeight:   2,
		Fit:           FitWrapClip,
	}

	got, err := provider(id, info)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "id:nav\nalloc:15x6\nframe:4x2\ncontent:20x5\nfit:WrapClip"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
