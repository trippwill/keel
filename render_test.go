package keel

import (
	"errors"
	"testing"

	gloss "github.com/charmbracelet/lipgloss"
	"github.com/trippwill/keel/core"
	"github.com/trippwill/keel/engine"
	"github.com/trippwill/keel/logging"
)

func makeContentProvider(result string) func(string, FrameInfo) (string, error) {
	return func(_ string, _ FrameInfo) (string, error) {
		return result, nil
	}
}

type dummySpec struct {
	extent ExtentConstraint
}

func (d dummySpec) Extent() ExtentConstraint { return d.extent }

type logEntry struct {
	event logging.LogEvent
	path  string
	msg   string
}

type flakyStack struct {
	axis  core.Axis
	slots []Spec
	calls []int
}

func (s flakyStack) Extent() ExtentConstraint { return FlexUnit() }

func (s flakyStack) Axis() core.Axis { return s.axis }

func (s flakyStack) Len() int { return len(s.slots) }

func (s flakyStack) Slot(index int) (Spec, bool) {
	if index < 0 || index >= len(s.slots) {
		return nil, false
	}
	s.calls[index]++
	if s.calls[index] > 1 {
		return nil, true
	}
	return s.slots[index], true
}

type countingStack struct {
	axis  core.Axis
	slots []Spec
	calls int
}

func (s countingStack) Extent() ExtentConstraint { return FlexUnit() }

func (s countingStack) Axis() core.Axis { return s.axis }

func (s countingStack) Len() int { return len(s.slots) }

func (s *countingStack) Slot(index int) (Spec, bool) {
	if index < 0 || index >= len(s.slots) {
		return nil, false
	}
	s.calls++
	return s.slots[index], true
}

func TestRender_UnknownSpec(t *testing.T) {
	renderer := NewRenderer[string](dummySpec{extent: FlexUnit()}, nil, nil)
	_, err := renderer.Render(Size{Width: 1, Height: 1})
	if !errors.Is(err, ErrUnknownSpec) {
		t.Fatalf("expected ErrUnknownSpec, got %v", err)
	}
}

func TestRenderStack_Empty(t *testing.T) {
	renderer := NewRenderer[string](Row(FlexUnit()), nil, nil)
	size := Size{Width: 5, Height: 2}
	got, err := renderer.Render(size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty output, got %q", got)
	}
}

func TestRenderMatchesLayoutRender(t *testing.T) {
	layout := Row(FlexUnit(),
		Exact(Fixed(2), "a"),
		Exact(FlexUnit(), "b"),
	)

	renderer := NewRenderer(layout, nil, func(id string, _ FrameInfo) (string, error) {
		switch id {
		case "a":
			return "aa", nil
		case "b":
			return "bbb", nil
		default:
			return "", &UnknownFrameIDError{ID: id}
		}
	})
	size := Size{Width: 5, Height: 1}

	want, err := renderer.Render(size)
	if err != nil {
		t.Fatalf("unexpected render error: %v", err)
	}

	arranged, err := engine.Arrange[string](layout, size, nil)
	if err != nil {
		t.Fatalf("unexpected arrange error: %v", err)
	}

	got, err := renderer.renderLayout(arranged)
	if err != nil {
		t.Fatalf("unexpected layout render error: %v", err)
	}
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestRenderStack_AllocTooSmallVertical(t *testing.T) {
	split := Col(
		FlexUnit(),
		Exact(FlexMin(1, 3), "a"),
		Exact(FlexMin(1, 3), "b"),
	)

	renderer := NewRenderer(split, nil, makeContentProvider(""))
	size := Size{Width: 2, Height: 4}

	_, err := renderer.Render(size)
	var tooSmall *ExtentTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected ExtentTooSmallError, got %v", err)
	}
	if tooSmall.Axis != core.AxisVertical {
		t.Fatalf("expected vertical axis, got %v", tooSmall.Axis)
	}
	if tooSmall.Source != "vertical split" {
		t.Fatalf("expected source %q, got %q", "vertical split", tooSmall.Source)
	}
	if tooSmall.Reason != "allocation" {
		t.Fatalf("expected reason %q, got %q", "allocation", tooSmall.Reason)
	}
}

func TestRenderStack_UnstableSlot(t *testing.T) {
	stack := flakyStack{
		axis:  core.AxisHorizontal,
		slots: []Spec{Exact(FlexUnit(), "a")},
		calls: make([]int, 1),
	}

	renderer := NewRenderer[string](stack, nil, nil)
	_, err := renderer.Render(Size{Width: 1, Height: 1})
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

func TestRenderStack_InvalidAxis(t *testing.T) {
	stack := testStack{
		axis:  core.Axis(99),
		slots: []Spec{Exact(FlexUnit(), "a")},
	}

	renderer := NewRenderer[string](stack, nil, nil)
	size := Size{Width: 1, Height: 1}
	_, err := renderer.Render(size)
	if !errors.Is(err, ErrInvalidAxis) {
		t.Fatalf("expected ErrInvalidAxis, got %v", err)
	}
}

func TestRenderStack_ArrangeSlotError(t *testing.T) {
	split := Row(FlexUnit(), nil)

	renderer := NewRenderer[string](split, nil, nil)
	size := Size{Width: 1, Height: 1}
	_, err := renderer.Render(size)
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

func TestRenderInvalidateClearsCachedLayout(t *testing.T) {
	spec := &countingStack{
		axis:  core.AxisHorizontal,
		slots: []Spec{Exact(FlexUnit(), "a")},
	}
	renderer := NewRenderer(spec, nil, makeContentProvider("ok"))
	size := Size{Width: 2, Height: 1}

	_, err := renderer.Render(size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	initialCalls := spec.calls

	_, err = renderer.Render(size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if spec.calls != initialCalls {
		t.Fatalf("expected cached layout reuse, got %d calls", spec.calls)
	}

	renderer.Invalidate()
	_, err = renderer.Render(size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if spec.calls == initialCalls {
		t.Fatalf("expected layout re-arrange after invalidate")
	}
}

func TestRender_LoggerEvents(t *testing.T) {
	var entries []logEntry
	logger := func(event logging.LogEvent, path, msg string) {
		entries = append(entries, logEntry{
			event: event,
			path:  path,
			msg:   msg,
		})
	}

	layout := Row(FlexUnit(),
		Exact(FlexUnit(), "a"),
	)
	renderer := NewRenderer(layout, nil, makeContentProvider(""))
	renderer.Config().SetLogger(logger)
	size := Size{Width: 3, Height: 1}

	_, err := renderer.Render(size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) < 2 {
		t.Fatalf("expected at least 2 log entries, got %d", len(entries))
	}

	if entries[0].event != logging.LogEventStackAlloc {
		t.Fatalf("expected stack.alloc first, got %q", entries[0].event)
	}
	if entries[0].path != "/" {
		t.Fatalf("expected root path, got %q", entries[0].path)
	}

	found := false
	for _, entry := range entries {
		if entry.event == logging.LogEventFrameRender && entry.path == "/0" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected frame.render for path /0")
	}
}

func TestRender_LoggerError(t *testing.T) {
	var entries []logEntry
	logger := func(event logging.LogEvent, path, msg string) {
		entries = append(entries, logEntry{
			event: event,
			path:  path,
			msg:   msg,
		})
	}

	renderer := NewRenderer[string](Exact(FlexUnit(), "a"), nil, nil)
	renderer.Config().SetLogger(logger)
	size := Size{Width: 1, Height: 1}
	_, err := renderer.Render(size)
	if err == nil {
		t.Fatalf("expected error")
	}

	found := false
	for _, entry := range entries {
		if entry.event == logging.LogEventRenderError {
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
		Exact(FlexMin(1, 3), "a"),
		Exact(FlexMin(1, 3), "b"),
	)

	renderer := NewRenderer(split, nil, makeContentProvider(""))
	size := Size{Width: 4, Height: 1}

	_, err := renderer.Render(size)
	var tooSmall *ExtentTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected ExtentTooSmallError, got %v", err)
	}
	if tooSmall.Axis != core.AxisHorizontal {
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
		layout Spec
	}{
		{name: "row", layout: Row(FlexUnit())},
		{name: "col", layout: Col(FlexUnit())},
	}

	size := Size{Width: 10, Height: 3}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			renderer := NewRenderer(tc.layout, nil, makeContentProvider(""))
			got, err := renderer.Render(size)
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
		Exact(Fixed(2), "a"),
	)

	renderer := NewRenderer(split, func(id string) *gloss.Style {
		if id != "a" {
			return nil
		}
		s := gloss.NewStyle().Border(gloss.NormalBorder())
		return &s
	}, makeContentProvider(""))
	size := Size{Width: 2, Height: 1}

	_, err := renderer.Render(size)
	var tooSmall *ExtentTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected ExtentTooSmallError, got %v", err)
	}
	if tooSmall.Axis != core.AxisVertical {
		t.Fatalf("expected vertical axis, got %v", tooSmall.Axis)
	}
	if tooSmall.Reason != "frame" {
		t.Fatalf("expected reason %q, got %q", "frame", tooSmall.Reason)
	}
}

func TestRenderSplit_SlotChromeTooWide(t *testing.T) {
	split := Row(
		FlexUnit(),
		Exact(Fixed(1), "a"),
	)

	renderer := NewRenderer(split, func(id string) *gloss.Style {
		if id != "a" {
			return nil
		}
		s := gloss.NewStyle().Border(gloss.NormalBorder())
		return &s
	}, makeContentProvider(""))
	size := Size{Width: 1, Height: 3}

	_, err := renderer.Render(size)
	var tooSmall *ExtentTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected ExtentTooSmallError, got %v", err)
	}
	if tooSmall.Axis != core.AxisHorizontal {
		t.Fatalf("expected horizontal axis, got %v", tooSmall.Axis)
	}
	if tooSmall.Reason != "frame" {
		t.Fatalf("expected reason %q, got %q", "frame", tooSmall.Reason)
	}
}

func TestRenderPanel_ContentProviderRequired(t *testing.T) {
	panel := Exact(FlexUnit(), "a")
	renderer := NewRenderer[string](panel, nil, nil)
	size := Size{Width: 1, Height: 1}
	_, err := renderer.Render(size)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRenderPanel_ContentProviderInfo(t *testing.T) {
	panel := Clip(FlexUnit(), "a")
	style := gloss.NewStyle().
		Border(gloss.NormalBorder()).
		Padding(1, 2).
		Margin(1, 1)
	var got FrameInfo
	var gotID string
	calls := 0
	renderer := NewRenderer(panel, func(id string) *gloss.Style {
		return &style
	}, func(id string, info FrameInfo) (string, error) {
		gotID = id
		got = info
		calls++
		return "ok", nil
	})
	size := Size{Width: 20, Height: 10}

	_, err := renderer.Render(size)
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
	if got.Fit != core.FitClip {
		t.Fatalf("expected fit %v, got %v", core.FitClip, got.Fit)
	}
}

func TestRenderPanel_ChromeTooLarge(t *testing.T) {
	panel := Exact(FlexUnit(), "a")
	renderer := NewRenderer(panel, func(id string) *gloss.Style {
		s := gloss.NewStyle().Border(gloss.NormalBorder())
		return &s
	}, makeContentProvider(""))
	size := Size{Width: 1, Height: 1}

	_, err := renderer.Render(size)
	var tooSmall *ExtentTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected ExtentTooSmallError, got %v", err)
	}
	if tooSmall.Reason != "frame" {
		t.Fatalf("expected reason %q, got %q", "frame", tooSmall.Reason)
	}
}

func TestRenderPanel_StyleApplied(t *testing.T) {
	panel := Exact(FlexUnit(), "a")
	style := gloss.NewStyle().Bold(true)
	renderer := NewRenderer(panel, func(id string) *gloss.Style {
		return &style
	}, makeContentProvider("hi"))
	size := Size{Width: 5, Height: 1}

	got, err := renderer.Render(size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := style.Width(size.Width).Height(size.Height).Render("hi")
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestRenderPanel_ContentTooWide(t *testing.T) {
	panel := Exact(FlexUnit(), "a")
	renderer := NewRenderer(panel, nil, makeContentProvider("abcd"))
	size := Size{Width: 2, Height: 1}

	_, err := renderer.Render(size)
	var tooSmall *ExtentTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected ExtentTooSmallError, got %v", err)
	}
	if tooSmall.Axis != core.AxisHorizontal {
		t.Fatalf("expected horizontal axis, got %v", tooSmall.Axis)
	}
	if tooSmall.Source != "frame a" {
		t.Fatalf("expected source %q, got %q", "frame a", tooSmall.Source)
	}
	if tooSmall.Reason != "content" {
		t.Fatalf("expected reason %q, got %q", "content", tooSmall.Reason)
	}
}

func TestRenderPanel_FitClipTruncatesContent(t *testing.T) {
	panel := Clip(FlexUnit(), "a")
	renderer := NewRenderer(panel, nil, makeContentProvider("abcd"))
	size := Size{Width: 2, Height: 1}

	got, err := renderer.Render(size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "ab" {
		t.Fatalf("expected %q, got %q", "ab", got)
	}
}

func TestRenderPanel_FitClipLargerThanContentStillFits(t *testing.T) {
	panel := Clip(FlexUnit(), "a")
	renderer := NewRenderer(panel, nil, makeContentProvider("ok"))
	size := Size{Width: 4, Height: 1}

	got, err := renderer.Render(size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	width, height := gloss.Size(got)
	if width != size.Width || height != size.Height {
		t.Fatalf("expected %dx%d, got %dx%d", size.Width, size.Height, width, height)
	}
}

func TestRenderPanel_FitWrapTruncatesHeight(t *testing.T) {
	panel := Wrap(FlexUnit(), "a")
	renderer := NewRenderer(panel, nil, makeContentProvider("abcd efgh"))
	size := Size{Width: 4, Height: 1}

	got, err := renderer.Render(size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "abcd" {
		t.Fatalf("expected %q, got %q", "abcd", got)
	}
}

func TestRenderPanel_FitOverflowAllowsOverflow(t *testing.T) {
	panel := Overflow(FlexUnit(), "a")
	renderer := NewRenderer(panel, nil, makeContentProvider("abcd"))
	size := Size{Width: 2, Height: 1}

	got, err := renderer.Render(size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	width, height := gloss.Size(got)
	if width != size.Width || height <= size.Height {
		t.Fatalf("expected overflow output %dx%d, got %dx%d", size.Width, size.Height+1, width, height)
	}
}

func TestRenderPanel_TransformAffectsSize(t *testing.T) {
	panel := Exact(FlexUnit(), "a")
	renderer := NewRenderer(panel, func(id string) *gloss.Style {
		s := gloss.NewStyle().Transform(func(s string) string {
			return s + "xx"
		})
		return &s
	}, makeContentProvider("ab"))
	size := Size{Width: 3, Height: 1}

	_, err := renderer.Render(size)
	var tooSmall *ExtentTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected ExtentTooSmallError, got %v", err)
	}
	if tooSmall.Axis != core.AxisHorizontal {
		t.Fatalf("expected horizontal axis, got %v", tooSmall.Axis)
	}
	if tooSmall.Reason != "content" {
		t.Fatalf("expected reason %q, got %q", "content", tooSmall.Reason)
	}
}

func TestRenderPanel_TransformAffectsHeight(t *testing.T) {
	panel := Exact(FlexUnit(), "a")
	renderer := NewRenderer(panel, func(id string) *gloss.Style {
		s := gloss.NewStyle().Transform(func(s string) string {
			return s + "\nX"
		})
		return &s
	}, makeContentProvider("hi"))
	size := Size{Width: 2, Height: 1}

	_, err := renderer.Render(size)
	var tooSmall *ExtentTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected ExtentTooSmallError, got %v", err)
	}
	if tooSmall.Axis != core.AxisVertical {
		t.Fatalf("expected vertical axis, got %v", tooSmall.Axis)
	}
	if tooSmall.Reason != "content" {
		t.Fatalf("expected reason %q, got %q", "content", tooSmall.Reason)
	}
}

func TestRenderPanel_FitWrapStrictTooTall(t *testing.T) {
	panel := WrapStrict(FlexUnit(), "a")
	renderer := NewRenderer(panel, nil, makeContentProvider("abcd efgh"))
	size := Size{Width: 4, Height: 1}

	_, err := renderer.Render(size)
	var tooSmall *ExtentTooSmallError
	if !errors.As(err, &tooSmall) {
		t.Fatalf("expected ExtentTooSmallError, got %v", err)
	}
	if tooSmall.Axis != core.AxisVertical {
		t.Fatalf("expected vertical axis, got %v", tooSmall.Axis)
	}
	if tooSmall.Reason != "content" {
		t.Fatalf("expected reason %q, got %q", "content", tooSmall.Reason)
	}
}

func TestRenderPanel_OutputMatchesAllocation(t *testing.T) {
	panel := Exact(FlexUnit(), "a")
	renderer := NewRenderer(panel, func(id string) *gloss.Style {
		s := gloss.NewStyle().
			Border(gloss.NormalBorder()).
			Padding(1).
			Margin(1)
		return &s
	}, makeContentProvider("hi"))
	size := Size{Width: 10, Height: 8}

	got, err := renderer.Render(size)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	width, height := gloss.Size(got)
	if width != size.Width || height != size.Height {
		t.Fatalf("expected %dx%d, got %dx%d", size.Width, size.Height, width, height)
	}
}
