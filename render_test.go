package keel

import (
	"errors"
	"testing"

	gloss "github.com/charmbracelet/lipgloss"
)

func TestRenderSplit_AllocTooSmall(t *testing.T) {
	split := Row("root",
		Slot[string]{Size: SizeSpec{Kind: SizeFlex, Units: 1, ContentMin: 3}, Node: Panel("a")},
		Slot[string]{Size: SizeSpec{Kind: SizeFlex, Units: 1, ContentMin: 3}, Node: Panel("b")},
	)

	ctx := Context[string]{
		Width:  4,
		Height: 1,
		ContentProvider: func(id string) (string, error) {
			return "", nil
		},
	}

	_, err := split.Render(ctx)
	var tooSmall *TargetTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected TargetTooSmallError, got %v", err)
	}
	if tooSmall.Axis != AxisHorizontal {
		t.Fatalf("expected horizontal axis, got %v", tooSmall.Axis)
	}
	if tooSmall.Need != 6 || tooSmall.Have != 4 {
		t.Fatalf("expected need 6 have 4, got need %d have %d", tooSmall.Need, tooSmall.Have)
	}
}

func TestRenderSplit_ContainerChromeTooLarge(t *testing.T) {
	split := Row("root",
		Slot[string]{Size: Flex(1), Node: Panel("a")},
	)

	ctx := Context[string]{
		Width:  1,
		Height: 1,
		StyleProvider: func(id string, kind NodeKind) *gloss.Style {
			if id == "root" && kind == NodeContainer {
				s := gloss.NewStyle().Border(gloss.NormalBorder())
				return &s
			}
			return nil
		},
		ContentProvider: func(id string) (string, error) {
			return "", nil
		},
	}

	_, err := split.Render(ctx)
	var tooSmall *TargetTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected TargetTooSmallError, got %v", err)
	}
	if tooSmall.Axis != AxisHorizontal {
		t.Fatalf("expected horizontal axis, got %v", tooSmall.Axis)
	}
}

func TestRenderSplit_ChildChromeTooTall(t *testing.T) {
	split := Row("root",
		Slot[string]{Size: Fixed(2), Node: Panel("a")},
	)

	ctx := Context[string]{
		Width:  2,
		Height: 1,
		StyleProvider: func(id string, kind NodeKind) *gloss.Style {
			if id == "a" && kind == NodeContent {
				s := gloss.NewStyle().Border(gloss.NormalBorder())
				return &s
			}
			return nil
		},
		ContentProvider: func(id string) (string, error) {
			return "", nil
		},
	}

	_, err := split.Render(ctx)
	var tooSmall *TargetTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected TargetTooSmallError, got %v", err)
	}
	if tooSmall.Axis != AxisVertical {
		t.Fatalf("expected vertical axis, got %v", tooSmall.Axis)
	}
}

func TestRenderSplit_ChildChromeTooWide(t *testing.T) {
	split := Row("root",
		Slot[string]{Size: Fixed(1), Node: Panel("a")},
	)

	ctx := Context[string]{
		Width:  1,
		Height: 3,
		StyleProvider: func(id string, kind NodeKind) *gloss.Style {
			if id == "a" && kind == NodeContent {
				s := gloss.NewStyle().Border(gloss.NormalBorder())
				return &s
			}
			return nil
		},
		ContentProvider: func(id string) (string, error) {
			return "", nil
		},
	}

	_, err := split.Render(ctx)
	var tooSmall *TargetTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected TargetTooSmallError, got %v", err)
	}
	if tooSmall.Axis != AxisHorizontal {
		t.Fatalf("expected horizontal axis, got %v", tooSmall.Axis)
	}
}

func TestRenderPanel_ContentProviderRequired(t *testing.T) {
	panel := Panel("a")
	ctx := Context[string]{Width: 1, Height: 1}
	_, err := panel.Render(ctx)
	if !errors.Is(err, ErrConfigurationInvalid) {
		t.Fatalf("expected configuration error, got %v", err)
	}
}

func TestRenderPanel_ChromeTooLarge(t *testing.T) {
	panel := Panel("a")
	ctx := Context[string]{
		Width:  1,
		Height: 1,
		StyleProvider: func(id string, kind NodeKind) *gloss.Style {
			s := gloss.NewStyle().Border(gloss.NormalBorder())
			return &s
		},
		ContentProvider: func(id string) (string, error) {
			return "", nil
		},
	}

	_, err := panel.Render(ctx)
	var tooSmall *TargetTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected TargetTooSmallError, got %v", err)
	}
}

func TestRenderPanel_StyleApplied(t *testing.T) {
	panel := Panel("a")
	style := gloss.NewStyle().Bold(true)
	ctx := Context[string]{
		Width:  5,
		Height: 1,
		StyleProvider: func(id string, kind NodeKind) *gloss.Style {
			return &style
		},
		ContentProvider: func(id string) (string, error) {
			return "hi", nil
		},
	}

	got, err := panel.Render(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := style.Render("hi")
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
