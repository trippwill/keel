package keel

import (
	"errors"
	"testing"

	gloss "github.com/charmbracelet/lipgloss"
)

func makeContentProvider(result string) func(RenderInfo[string]) (string, error) {
	return func(info RenderInfo[string]) (string, error) {
		return result, nil
	}
}

func TestRenderSplit_AllocTooSmall(t *testing.T) {
	split := Row(
		FlexUnit(),
		Panel(FlexMin(1, 3), "a"),
		Panel(FlexMin(1, 3), "b"),
	)

	ctx := Context[string]{
		Width:           4,
		Height:          1,
		ContentProvider: makeContentProvider(""),
	}

	_, err := Render(split, ctx)
	var tooSmall *ExtentTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected ExtentTooSmallError, got %v", err)
	}
	if tooSmall.Axis != AxisHorizontal {
		t.Fatalf("expected horizontal axis, got %v", tooSmall.Axis)
	}
	if tooSmall.Source != "horizontal split" {
		t.Fatalf("expected source %q, got %q", "horizontal split", tooSmall.Source)
	}
	if tooSmall.Reason != "allocation" {
		t.Fatalf("expected reason %q, got %q", "allocation", tooSmall.Reason)
	}
}

func TestRenderSplit_ChildChromeTooTall(t *testing.T) {
	split := Row(
		FlexUnit(),
		Panel(Fixed(2), "a"),
	)

	ctx := Context[string]{
		Width:  2,
		Height: 1,
		StyleProvider: func(id string) *gloss.Style {
			if id != "a" {
				return nil
			}
			s := gloss.NewStyle().Border(gloss.NormalBorder())
			return &s
		},
		ContentProvider: makeContentProvider(""),
	}

	_, err := Render(split, ctx)
	var tooSmall *ExtentTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected ExtentTooSmallError, got %v", err)
	}
	if tooSmall.Axis != AxisVertical {
		t.Fatalf("expected vertical axis, got %v", tooSmall.Axis)
	}
	if tooSmall.Reason != "frame" {
		t.Fatalf("expected reason %q, got %q", "frame", tooSmall.Reason)
	}
}

func TestRenderSplit_ChildChromeTooWide(t *testing.T) {
	split := Row(
		FlexUnit(),
		Panel(Fixed(1), "a"),
	)

	ctx := Context[string]{
		Width:  1,
		Height: 3,
		StyleProvider: func(id string) *gloss.Style {
			if id != "a" {
				return nil
			}
			s := gloss.NewStyle().Border(gloss.NormalBorder())
			return &s
		},
		ContentProvider: makeContentProvider(""),
	}

	_, err := Render(split, ctx)
	var tooSmall *ExtentTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected ExtentTooSmallError, got %v", err)
	}
	if tooSmall.Axis != AxisHorizontal {
		t.Fatalf("expected horizontal axis, got %v", tooSmall.Axis)
	}
	if tooSmall.Reason != "frame" {
		t.Fatalf("expected reason %q, got %q", "frame", tooSmall.Reason)
	}
}

func TestRenderPanel_ContentProviderRequired(t *testing.T) {
	panel := Panel(FlexUnit(), "a")
	ctx := Context[string]{Width: 1, Height: 1}
	_, err := Render(panel, ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRenderPanel_ContentProviderInfo(t *testing.T) {
	panel := PanelClip(FlexUnit(), Clip(3, 2), "a")
	style := gloss.NewStyle().
		Border(gloss.NormalBorder()).
		Padding(1, 2).
		Margin(1, 1)
	var got RenderInfo[string]
	calls := 0
	ctx := Context[string]{
		Width:  20,
		Height: 10,
		StyleProvider: func(id string) *gloss.Style {
			return &style
		},
		ContentProvider: func(info RenderInfo[string]) (string, error) {
			got = info
			calls++
			return "ok", nil
		},
	}

	_, err := Render(panel, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 content call, got %d", calls)
	}

	frameWidth, frameHeight := style.GetFrameSize()
	if got.ID != "a" {
		t.Fatalf("expected ID %q, got %q", "a", got.ID)
	}
	if got.Width != ctx.Width || got.Height != ctx.Height {
		t.Fatalf("expected %dx%d, got %dx%d", ctx.Width, ctx.Height, got.Width, got.Height)
	}
	if got.FrameWidth != frameWidth || got.FrameHeight != frameHeight {
		t.Fatalf("expected frame %dx%d, got %dx%d", frameWidth, frameHeight, got.FrameWidth, got.FrameHeight)
	}
	if got.ContentWidth != ctx.Width-frameWidth || got.ContentHeight != ctx.Height-frameHeight {
		t.Fatalf(
			"expected content %dx%d, got %dx%d",
			ctx.Width-frameWidth,
			ctx.Height-frameHeight,
			got.ContentWidth,
			got.ContentHeight,
		)
	}
	if got.Clip != (ClipConstraint{Width: 3, Height: 2}) {
		t.Fatalf("expected clip %+v, got %+v", ClipConstraint{Width: 3, Height: 2}, got.Clip)
	}
}

func TestRenderPanel_ChromeTooLarge(t *testing.T) {
	panel := Panel(FlexUnit(), "a")
	ctx := Context[string]{
		Width:  1,
		Height: 1,
		StyleProvider: func(id string) *gloss.Style {
			s := gloss.NewStyle().Border(gloss.NormalBorder())
			return &s
		},
		ContentProvider: makeContentProvider(""),
	}

	_, err := Render(panel, ctx)
	var tooSmall *ExtentTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected ExtentTooSmallError, got %v", err)
	}
	if tooSmall.Reason != "frame" {
		t.Fatalf("expected reason %q, got %q", "frame", tooSmall.Reason)
	}
}

func TestRenderPanel_StyleApplied(t *testing.T) {
	panel := Panel(FlexUnit(), "a")
	style := gloss.NewStyle().Bold(true)
	ctx := Context[string]{
		Width:  5,
		Height: 1,
		StyleProvider: func(id string) *gloss.Style {
			return &style
		},
		ContentProvider: makeContentProvider("hi"),
	}

	got, err := Render(panel, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := style.Width(ctx.Width).Height(ctx.Height).Render("hi")
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestRenderPanel_ContentTooWide(t *testing.T) {
	panel := Panel(FlexUnit(), "a")
	ctx := Context[string]{
		Width:           2,
		Height:          1,
		ContentProvider: makeContentProvider("abcd"),
	}

	_, err := Render(panel, ctx)
	var tooSmall *ExtentTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected ExtentTooSmallError, got %v", err)
	}
	if tooSmall.Axis != AxisHorizontal {
		t.Fatalf("expected horizontal axis, got %v", tooSmall.Axis)
	}
	if tooSmall.Source != "block a" {
		t.Fatalf("expected source %q, got %q", "block a", tooSmall.Source)
	}
	if tooSmall.Reason != "content" {
		t.Fatalf("expected reason %q, got %q", "content", tooSmall.Reason)
	}
}

func TestRenderPanel_ClipAppliesToContent(t *testing.T) {
	panel := PanelClip(FlexUnit(), ClipWidth(2), "a")
	ctx := Context[string]{
		Width:           2,
		Height:          1,
		ContentProvider: makeContentProvider("abcd"),
	}

	got, err := Render(panel, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "ab" {
		t.Fatalf("expected %q, got %q", "ab", got)
	}
}

func TestRenderPanel_ClipLargerThanContentStillFits(t *testing.T) {
	panel := PanelClip(FlexUnit(), ClipWidth(10), "a")
	ctx := Context[string]{
		Width:           4,
		Height:          1,
		ContentProvider: makeContentProvider("ok"),
	}

	got, err := Render(panel, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	width, height := gloss.Size(got)
	if width != ctx.Width || height != ctx.Height {
		t.Fatalf("expected %dx%d, got %dx%d", ctx.Width, ctx.Height, width, height)
	}
}

func TestRenderPanel_ClipTooWide(t *testing.T) {
	panel := PanelClip(FlexUnit(), ClipWidth(3), "a")
	ctx := Context[string]{
		Width:           2,
		Height:          1,
		ContentProvider: makeContentProvider("abcd"),
	}

	_, err := Render(panel, ctx)
	var tooSmall *ExtentTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected ExtentTooSmallError, got %v", err)
	}
	if tooSmall.Reason != "content" {
		t.Fatalf("expected reason %q, got %q", "content", tooSmall.Reason)
	}
}

func TestRenderPanel_TransformAffectsSize(t *testing.T) {
	panel := Panel(FlexUnit(), "a")
	ctx := Context[string]{
		Width:  3,
		Height: 1,
		StyleProvider: func(id string) *gloss.Style {
			s := gloss.NewStyle().Transform(func(s string) string {
				return s + "xx"
			})
			return &s
		},
		ContentProvider: makeContentProvider("ab"),
	}

	_, err := Render(panel, ctx)
	var tooSmall *ExtentTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected ExtentTooSmallError, got %v", err)
	}
	if tooSmall.Axis != AxisHorizontal {
		t.Fatalf("expected horizontal axis, got %v", tooSmall.Axis)
	}
	if tooSmall.Reason != "content" {
		t.Fatalf("expected reason %q, got %q", "content", tooSmall.Reason)
	}
}

func TestRenderPanel_TransformAffectsHeight(t *testing.T) {
	panel := Panel(FlexUnit(), "a")
	ctx := Context[string]{
		Width:  2,
		Height: 1,
		StyleProvider: func(id string) *gloss.Style {
			s := gloss.NewStyle().Transform(func(s string) string {
				return s + "\nX"
			})
			return &s
		},
		ContentProvider: makeContentProvider("hi"),
	}

	_, err := Render(panel, ctx)
	var tooSmall *ExtentTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected ExtentTooSmallError, got %v", err)
	}
	if tooSmall.Axis != AxisVertical {
		t.Fatalf("expected vertical axis, got %v", tooSmall.Axis)
	}
	if tooSmall.Reason != "content" {
		t.Fatalf("expected reason %q, got %q", "content", tooSmall.Reason)
	}
}

func TestRenderPanel_ClipWithFrameTooWide(t *testing.T) {
	panel := PanelClip(FlexUnit(), ClipWidth(5), "a")
	ctx := Context[string]{
		Width:  6,
		Height: 3,
		StyleProvider: func(id string) *gloss.Style {
			s := gloss.NewStyle().Border(gloss.NormalBorder())
			return &s
		},
		ContentProvider: makeContentProvider("abcdef"),
	}

	_, err := Render(panel, ctx)
	var tooSmall *ExtentTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected ExtentTooSmallError, got %v", err)
	}
	if tooSmall.Reason != "content" {
		t.Fatalf("expected reason %q, got %q", "content", tooSmall.Reason)
	}
}

func TestRenderPanel_OutputMatchesAllocation(t *testing.T) {
	panel := Panel(FlexUnit(), "a")
	ctx := Context[string]{
		Width:  10,
		Height: 8,
		StyleProvider: func(id string) *gloss.Style {
			s := gloss.NewStyle().
				Border(gloss.NormalBorder()).
				Padding(1).
				Margin(1)
			return &s
		},
		ContentProvider: makeContentProvider("hi"),
	}

	got, err := Render(panel, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	width, height := gloss.Size(got)
	if width != ctx.Width || height != ctx.Height {
		t.Fatalf("expected %dx%d, got %dx%d", ctx.Width, ctx.Height, width, height)
	}
}
