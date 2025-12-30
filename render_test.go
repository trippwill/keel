package keel

import (
	"errors"
	"testing"

	gloss "github.com/charmbracelet/lipgloss"
)

func makeContentProvider(result string) func(string, RenderInfo) (string, error) {
	return func(_ string, _ RenderInfo) (string, error) {
		return result, nil
	}
}

type dummyRenderable struct {
	extent ExtentConstraint
}

func (d dummyRenderable) GetExtent() ExtentConstraint { return d.extent }

type logEntry struct {
	event LogEvent
	path  string
	msg   string
}

type flakyContainer struct {
	axis  Axis
	slots []Renderable
	calls []int
}

func (c flakyContainer) GetExtent() ExtentConstraint { return FlexUnit() }

func (c flakyContainer) GetAxis() Axis { return c.axis }

func (c flakyContainer) Len() int { return len(c.slots) }

func (c flakyContainer) Slot(index int) (Renderable, bool) {
	if index < 0 || index >= len(c.slots) {
		return nil, false
	}
	c.calls[index]++
	if c.calls[index] > 1 {
		return nil, true
	}
	return c.slots[index], true
}

func TestRender_UnknownRenderable(t *testing.T) {
	ctx := Context[string]{}
	size := Size{Width: 1, Height: 1}
	_, err := Render(ctx, dummyRenderable{extent: FlexUnit()}, size)
	if !errors.Is(err, ErrUnknownRenderable) {
		t.Fatalf("expected ErrUnknownRenderable, got %v", err)
	}
}

func TestRenderContainer_Empty(t *testing.T) {
	ctx := Context[string]{}
	size := Size{Width: 5, Height: 2}
	got, err := Render(ctx, Row(FlexUnit()), size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty output, got %q", got)
	}
}

func TestRenderContainer_AllocTooSmallVertical(t *testing.T) {
	split := Col(
		FlexUnit(),
		Panel(FlexMin(1, 3), "a"),
		Panel(FlexMin(1, 3), "b"),
	)

	ctx := Context[string]{
		ContentProvider: makeContentProvider(""),
	}
	size := Size{Width: 2, Height: 4}

	_, err := Render(ctx, split, size)
	var tooSmall *ExtentTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected ExtentTooSmallError, got %v", err)
	}
	if tooSmall.Axis != AxisVertical {
		t.Fatalf("expected vertical axis, got %v", tooSmall.Axis)
	}
	if tooSmall.Source != "vertical split" {
		t.Fatalf("expected source %q, got %q", "vertical split", tooSmall.Source)
	}
	if tooSmall.Reason != "allocation" {
		t.Fatalf("expected reason %q, got %q", "allocation", tooSmall.Reason)
	}
}

func TestRenderContainer_UnstableSlot(t *testing.T) {
	container := flakyContainer{
		axis:  AxisHorizontal,
		slots: []Renderable{Panel(FlexUnit(), "a")},
		calls: make([]int, 1),
	}

	ctx := Context[string]{}
	size := Size{Width: 1, Height: 1}
	_, err := Render(ctx, container, size)
	var slotErr *SlotError
	if !errors.As(err, &slotErr) {
		t.Fatalf("expected SlotError, got %v", err)
	}
	if slotErr.Index != 0 {
		t.Fatalf("expected index 0, got %d", slotErr.Index)
	}
	if !errors.Is(err, ErrNilSlot) {
		t.Fatalf("expected ErrNilSlot")
	}
}

func TestRenderContainer_InvalidAxis(t *testing.T) {
	container := flakyContainer{
		axis:  Axis(99),
		slots: []Renderable{Panel(FlexUnit(), "a")},
		calls: make([]int, 1),
	}

	ctx := Context[string]{}
	size := Size{Width: 1, Height: 1}
	_, err := RenderContainer(ctx, container, size)
	if !errors.Is(err, ErrInvalidAxis) {
		t.Fatalf("expected ErrInvalidAxis, got %v", err)
	}
}

func TestRenderContainer_ResolverSlotError(t *testing.T) {
	split := Row(FlexUnit(), nil)

	ctx := Context[string]{}
	size := Size{Width: 1, Height: 1}
	_, err := Render(ctx, split, size)
	var slotErr *SlotError
	if !errors.As(err, &slotErr) {
		t.Fatalf("expected SlotError, got %v", err)
	}
	if slotErr.Index != 0 {
		t.Fatalf("expected index 0, got %d", slotErr.Index)
	}
	if !errors.Is(err, ErrNilSlot) {
		t.Fatalf("expected ErrNilSlot")
	}
}

func TestRender_LoggerEvents(t *testing.T) {
	var entries []logEntry
	logger := func(event LogEvent, path, msg string) {
		entries = append(entries, logEntry{
			event: event,
			path:  path,
			msg:   msg,
		})
	}

	layout := Row(FlexUnit(),
		Panel(FlexUnit(), "a"),
	)
	ctx := Context[string]{
		ContentProvider: makeContentProvider(""),
	}.WithLogger(logger)
	size := Size{Width: 3, Height: 1}

	_, err := Render(ctx, layout, size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) < 2 {
		t.Fatalf("expected at least 2 log entries, got %d", len(entries))
	}

	if entries[0].event != LogEventContainerAlloc {
		t.Fatalf("expected container.alloc first, got %q", entries[0].event)
	}
	if entries[0].path != "/" {
		t.Fatalf("expected root path, got %q", entries[0].path)
	}

	found := false
	for _, entry := range entries {
		if entry.event == LogEventBlockRender && entry.path == "/0" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected block.render for path /0")
	}
}

func TestRender_LoggerError(t *testing.T) {
	var entries []logEntry
	logger := func(event LogEvent, path, msg string) {
		entries = append(entries, logEntry{
			event: event,
			path:  path,
			msg:   msg,
		})
	}

	ctx := Context[string]{}.WithLogger(logger)
	size := Size{Width: 1, Height: 1}
	_, err := Render(ctx, Panel(FlexUnit(), "a"), size)
	if err == nil {
		t.Fatalf("expected error")
	}

	found := false
	for _, entry := range entries {
		if entry.event == LogEventRenderError {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected render.error entry")
	}
}

func TestRenderSplit_AllocTooSmall(t *testing.T) {
	split := Row(
		FlexUnit(),
		Panel(FlexMin(1, 3), "a"),
		Panel(FlexMin(1, 3), "b"),
	)

	ctx := Context[string]{
		ContentProvider: makeContentProvider(""),
	}
	size := Size{Width: 4, Height: 1}

	_, err := Render(ctx, split, size)
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

func TestRenderSplit_EmptySlots(t *testing.T) {
	cases := []struct {
		name   string
		layout Renderable
	}{
		{name: "row", layout: Row(FlexUnit())},
		{name: "col", layout: Col(FlexUnit())},
	}

	ctx := Context[string]{
		ContentProvider: makeContentProvider(""),
	}
	size := Size{Width: 10, Height: 3}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Render(ctx, tc.layout, size)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != "" {
				t.Fatalf("expected empty output, got %q", got)
			}
		})
	}
}

func TestRenderSplit_SlotChromeTooTall(t *testing.T) {
	split := Row(
		FlexUnit(),
		Panel(Fixed(2), "a"),
	)

	ctx := Context[string]{
		StyleProvider: func(id string) *gloss.Style {
			if id != "a" {
				return nil
			}
			s := gloss.NewStyle().Border(gloss.NormalBorder())
			return &s
		},
		ContentProvider: makeContentProvider(""),
	}
	size := Size{Width: 2, Height: 1}

	_, err := Render(ctx, split, size)
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

func TestRenderSplit_SlotChromeTooWide(t *testing.T) {
	split := Row(
		FlexUnit(),
		Panel(Fixed(1), "a"),
	)

	ctx := Context[string]{
		StyleProvider: func(id string) *gloss.Style {
			if id != "a" {
				return nil
			}
			s := gloss.NewStyle().Border(gloss.NormalBorder())
			return &s
		},
		ContentProvider: makeContentProvider(""),
	}
	size := Size{Width: 1, Height: 3}

	_, err := Render(ctx, split, size)
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
	ctx := Context[string]{}
	size := Size{Width: 1, Height: 1}
	_, err := Render(ctx, panel, size)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRenderPanel_ContentProviderInfo(t *testing.T) {
	panel := PanelClip(FlexUnit(), "a")
	style := gloss.NewStyle().
		Border(gloss.NormalBorder()).
		Padding(1, 2).
		Margin(1, 1)
	var got RenderInfo
	var gotID string
	calls := 0
	ctx := Context[string]{
		StyleProvider: func(id string) *gloss.Style {
			return &style
		},
		ContentProvider: func(id string, info RenderInfo) (string, error) {
			gotID = id
			got = info
			calls++
			return "ok", nil
		},
	}
	size := Size{Width: 20, Height: 10}

	_, err := Render(ctx, panel, size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 content call, got %d", calls)
	}

	frameWidth, frameHeight := style.GetFrameSize()
	if gotID != "a" {
		t.Fatalf("expected ID %q, got %q", "a", gotID)
	}
	if got.Width != size.Width || got.Height != size.Height {
		t.Fatalf("expected %dx%d, got %dx%d", size.Width, size.Height, got.Width, got.Height)
	}
	if got.FrameWidth != frameWidth || got.FrameHeight != frameHeight {
		t.Fatalf("expected frame %dx%d, got %dx%d", frameWidth, frameHeight, got.FrameWidth, got.FrameHeight)
	}
	if got.ContentWidth != size.Width-frameWidth || got.ContentHeight != size.Height-frameHeight {
		t.Fatalf(
			"expected content %dx%d, got %dx%d",
			size.Width-frameWidth,
			size.Height-frameHeight,
			got.ContentWidth,
			got.ContentHeight,
		)
	}
	if got.Fit != FitClip {
		t.Fatalf("expected fit %v, got %v", FitClip, got.Fit)
	}
}

func TestRenderPanel_ChromeTooLarge(t *testing.T) {
	panel := Panel(FlexUnit(), "a")
	ctx := Context[string]{
		StyleProvider: func(id string) *gloss.Style {
			s := gloss.NewStyle().Border(gloss.NormalBorder())
			return &s
		},
		ContentProvider: makeContentProvider(""),
	}
	size := Size{Width: 1, Height: 1}

	_, err := Render(ctx, panel, size)
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
		StyleProvider: func(id string) *gloss.Style {
			return &style
		},
		ContentProvider: makeContentProvider("hi"),
	}
	size := Size{Width: 5, Height: 1}

	got, err := Render(ctx, panel, size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := style.Width(size.Width).Height(size.Height).Render("hi")
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestRenderPanel_ContentTooWide(t *testing.T) {
	panel := Panel(FlexUnit(), "a")
	ctx := Context[string]{
		ContentProvider: makeContentProvider("abcd"),
	}
	size := Size{Width: 2, Height: 1}

	_, err := Render(ctx, panel, size)
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

func TestRenderPanel_FitClipTruncatesContent(t *testing.T) {
	panel := PanelClip(FlexUnit(), "a")
	ctx := Context[string]{
		ContentProvider: makeContentProvider("abcd"),
	}
	size := Size{Width: 2, Height: 1}

	got, err := Render(ctx, panel, size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "ab" {
		t.Fatalf("expected %q, got %q", "ab", got)
	}
}

func TestRenderPanel_FitClipLargerThanContentStillFits(t *testing.T) {
	panel := PanelClip(FlexUnit(), "a")
	ctx := Context[string]{
		ContentProvider: makeContentProvider("ok"),
	}
	size := Size{Width: 4, Height: 1}

	got, err := Render(ctx, panel, size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	width, height := gloss.Size(got)
	if width != size.Width || height != size.Height {
		t.Fatalf("expected %dx%d, got %dx%d", size.Width, size.Height, width, height)
	}
}

func TestRenderPanel_FitWrapTruncatesHeight(t *testing.T) {
	panel := PanelWrap(FlexUnit(), "a")
	ctx := Context[string]{
		ContentProvider: makeContentProvider("abcd efgh"),
	}
	size := Size{Width: 4, Height: 1}

	got, err := Render(ctx, panel, size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "abcd" {
		t.Fatalf("expected %q, got %q", "abcd", got)
	}
}

func TestRenderPanel_FitOverflowAllowsOverflow(t *testing.T) {
	panel := PanelOverflow(FlexUnit(), "a")
	ctx := Context[string]{
		ContentProvider: makeContentProvider("abcd"),
	}
	size := Size{Width: 2, Height: 1}

	got, err := Render(ctx, panel, size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	width, height := gloss.Size(got)
	if width != size.Width || height <= size.Height {
		t.Fatalf("expected overflow output %dx%d, got %dx%d", size.Width, size.Height+1, width, height)
	}
}

func TestRenderPanel_TransformAffectsSize(t *testing.T) {
	panel := Panel(FlexUnit(), "a")
	ctx := Context[string]{
		StyleProvider: func(id string) *gloss.Style {
			s := gloss.NewStyle().Transform(func(s string) string {
				return s + "xx"
			})
			return &s
		},
		ContentProvider: makeContentProvider("ab"),
	}
	size := Size{Width: 3, Height: 1}

	_, err := Render(ctx, panel, size)
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
		StyleProvider: func(id string) *gloss.Style {
			s := gloss.NewStyle().Transform(func(s string) string {
				return s + "\nX"
			})
			return &s
		},
		ContentProvider: makeContentProvider("hi"),
	}
	size := Size{Width: 2, Height: 1}

	_, err := Render(ctx, panel, size)
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

func TestRenderPanel_FitWrapStrictTooTall(t *testing.T) {
	panel := PanelWrapStrict(FlexUnit(), "a")
	ctx := Context[string]{
		ContentProvider: makeContentProvider("abcd efgh"),
	}
	size := Size{Width: 4, Height: 1}

	_, err := Render(ctx, panel, size)
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

func TestRenderPanel_OutputMatchesAllocation(t *testing.T) {
	panel := Panel(FlexUnit(), "a")
	ctx := Context[string]{
		StyleProvider: func(id string) *gloss.Style {
			s := gloss.NewStyle().
				Border(gloss.NormalBorder()).
				Padding(1).
				Margin(1)
			return &s
		},
		ContentProvider: makeContentProvider("hi"),
	}
	size := Size{Width: 10, Height: 8}

	got, err := Render(ctx, panel, size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	width, height := gloss.Size(got)
	if width != size.Width || height != size.Height {
		t.Fatalf("expected %dx%d, got %dx%d", size.Width, size.Height, width, height)
	}
}
