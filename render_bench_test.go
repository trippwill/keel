package keel_test

import (
	"testing"

	"github.com/trippwill/chiplog/keel"
	"github.com/trippwill/chiplog/keel/examples"
)

func BenchmarkRenderExampleSplit(b *testing.B) {
	layout := examples.ExampleSplit()
	ctx := keel.Context[string]{
		Width:           70,
		Height:          13,
		StyleProvider:   examples.ExampleSplitStyleProvider,
		ContentProvider: examples.ExampleSplitContentProvider,
	}

	b.ReportAllocs()
	for b.Loop() {
		if _, err := keel.Render(layout, ctx); err != nil {
			b.Fatal(err)
		}
	}
}
