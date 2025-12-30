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

	spec := examples.ExampleSplit()
	context := keel.NewContext(
		examples.ExampleSplitStyleProvider,
		examples.ExampleSplitContentProvider)
	size := keel.Size{Width: *width, Height: *height}

	if *debug {
		context = context.WithContentProvider(keel.DefaultDebugProvider)
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

	layout, err := keel.Arrange(context, spec, size)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rendered, err := keel.Render(context, layout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(rendered)
}
