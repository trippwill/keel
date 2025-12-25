package keel

import (
	"fmt"

	gloss "github.com/charmbracelet/lipgloss"
)

func ExampleSplit() *SplitSpec[string] {
	type S = Slot[string]
	return Split("root", AxisVertical,
		S{Fixed(30), Row("row1",
			S{Flex(1), Panel("panel1")},
			S{Flex(2), Panel("panel2")},
		)},
	)
}

func ExampleSplitContentProvider(id string) (string, error) {
	switch id {
	case "panel1":
		return "Hello,", nil
	case "panel2":
		return "World!", nil
	default:
		return "", fmt.Errorf("unexpected panel ID: %s", id)
	}
}

func ExampleSplitStyleProvider(id string, kind NodeKind) *gloss.Style {
	// Decorate all containers with a border and background
	if kind == NodeContainer {
		s := gloss.NewStyle().
			Background(gloss.Color("8")).
			Padding(1, 2).
			Border(gloss.DoubleBorder())

		return &s
	}

	switch id {
	case "panel1":
		s := gloss.NewStyle().
			Foreground(gloss.Color("2")).
			Bold(true)
		return &s
	case "panel2":
		s := gloss.NewStyle().
			Foreground(gloss.Color("4")).
			Italic(true)
		return &s
	default:
		return nil
	}
}
