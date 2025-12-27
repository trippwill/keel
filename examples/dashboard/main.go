package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/trippwill/chiplog/keel"
	"github.com/trippwill/chiplog/keel/examples"
)

func main() {
	width := flag.Int("width", 61, "Width of the layout")
	height := flag.Int("height", 13, "Height of the layout")
	debug := flag.Bool("debug", false, "Enable debug output")

	flag.Parse()

	layout := examples.ExampleSplit()
	context := keel.Context[string]{}.
		WithSize(*width, *height).
		WithContentProvider(examples.ExampleSplitContentProvider).
		WithStyleProvider(examples.ExampleSplitStyleProvider)

	if *debug {
		context = context.WithContentProvider(keel.DefaultDebugProvider[string])
	}

	rendered, err := keel.Render(layout, context)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(rendered)
}
