package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/trippwill/keel"
	"github.com/trippwill/keel/examples"
)

func main() {
	width := flag.Int("width", 61, "Width of the layout")
	height := flag.Int("height", 13, "Height of the layout")
	debug := flag.Bool("debug", false, "Enable debug output")
	logPath := flag.String("log", "", "Write render logs to the given file path")

	flag.Parse()

	layout := examples.ExampleSplit()
	context := keel.Context[string]{}.
		WithSize(*width, *height).
		WithContentProvider(examples.ExampleSplitContentProvider).
		WithStyleProvider(examples.ExampleSplitStyleProvider)

	if *debug {
		context = context.WithContentProvider(keel.DefaultDebugProvider[string])
	}
	if *logPath != "" {
		logger, file, err := keel.NewFileLoggerPath(*logPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				fmt.Println(err)
			}
		}()
		context = context.WithLogger(logger.Log)
	}

	rendered, err := keel.Render(layout, context)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(rendered)
}
