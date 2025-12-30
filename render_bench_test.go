package keel_test

import (
	"testing"

	"github.com/trippwill/keel"
	"github.com/trippwill/keel/examples"
)

func BenchmarkRenderExampleSplit(b *testing.B) {
	layout := examples.ExampleSplit()
	ctx := keel.Context[string]{
		StyleProvider:   examples.ExampleSplitStyleProvider,
		ContentProvider: examples.ExampleSplitContentProvider,
	}
	size := keel.Size{Width: 70, Height: 13}

	b.ReportAllocs()
	for b.Loop() {
		if _, err := keel.Render(ctx, layout, size); err != nil {
			b.Fatal(err)
		}
	}
}
