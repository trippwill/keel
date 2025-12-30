package examples

import (
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/trippwill/keel"
)

// ExampleSplit returns the example layout hierarchy used across demos and tests.
func ExampleSplit() keel.Renderable {
	return keel.Col(keel.FlexUnit(),
		keel.PanelClip(keel.Fixed(3), "header"),
		keel.Row(keel.FlexMin(1, 6),
			keel.PanelWrap(keel.FlexUnit(), "nav"),
			keel.Panel(keel.FlexMin(2, 8), "feed"),
			keel.Panel(keel.Fixed(19), "detail"),
		),
		keel.Row(keel.Fixed(3),
			keel.PanelClip(keel.FlexMax(1, 10), "status"),
			keel.Panel(keel.FlexUnit(), "help"),
		),
	)
}

// ExampleResolvedSplit resolves the example layout at the given size.
func ExampleResolvedSplit(width, height int) (keel.Resolved[string], error) {
	layout := ExampleSplit()
	ctx := keel.Context[string]{}
	size := keel.Size{Width: width, Height: height}
	return keel.Resolve[string](ctx, layout, size)
}

// ExampleSplitContentProvider returns content for the example layout.
func ExampleSplitContentProvider(id string, _ keel.RenderInfo) (string, error) {
	switch id {
	case "header":
		return "Chiplog Dashboard", nil
	case "nav":
		return "Queues\n- ingest\n- parse\n- render\n- ship", nil
	case "feed":
		return "Latest\n- build ok\n- cache warm\n- alloc pass", nil
	case "detail":
		return "Detail\nid: 42\nstatus: running", nil
	case "status":
		return "status: connected", nil
	case "help":
		return "?: help  q: quit", nil
	default:
		return "", &keel.UnknownBlockIDError{ID: id}
	}
}

var (
	headerStyle = gloss.NewStyle().
			Border(gloss.RoundedBorder()).
			Padding(0, 1).
			Bold(true)
	navStyle = gloss.NewStyle().
			Border(gloss.NormalBorder()).
			Padding(0, 1).
			Foreground(gloss.Color("6"))
	feedStyle = gloss.NewStyle().
			Border(gloss.NormalBorder()).
			Padding(0, 1).
			Foreground(gloss.Color("2"))
	detailStyle = gloss.NewStyle().
			Border(gloss.NormalBorder()).
			Padding(0, 1).
			Foreground(gloss.Color("3"))
	statusStyle = gloss.NewStyle().
			Border(gloss.NormalBorder()).
			Padding(0, 1).
			Foreground(gloss.Color("10"))
	helpStyle = gloss.NewStyle().
			Border(gloss.NormalBorder()).
			Padding(0, 1).
			Foreground(gloss.Color("4"))
)

// ExampleSplitStyleProvider returns cached styles for the example layout.
func ExampleSplitStyleProvider(id string) *gloss.Style {
	switch id {
	case "header":
		return &headerStyle
	case "nav":
		return &navStyle
	case "feed":
		return &feedStyle
	case "detail":
		return &detailStyle
	case "status":
		return &statusStyle
	case "help":
		return &helpStyle
	default:
		return nil
	}
}
