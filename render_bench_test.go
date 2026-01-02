package keel_test

import (
	"testing"

	"github.com/trippwill/keel"
	"github.com/trippwill/keel/examples"
)

func BenchmarkRenderExampleSplit(b *testing.B) {
	layout := examples.ExampleSplit()
	renderer := keel.NewRenderer(
		layout,
		examples.ExampleSplitStyleProvider,
		examples.ExampleSplitContentProvider,
	)
	size := keel.Size{Width: 70, Height: 13}

	b.ReportAllocs()
	for b.Loop() {
		if _, err := renderer.Render(size); err != nil {
			b.Fatal(err)
		}
	}
}
