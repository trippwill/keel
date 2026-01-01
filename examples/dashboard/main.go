package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/trippwill/keel"
	"github.com/trippwill/keel/examples"
	"github.com/trippwill/keel/logging"
)

func main() {
	width := flag.Int("width", 61, "Width of the layout")
	height := flag.Int("height", 13, "Height of the layout")
	debug := flag.Bool("debug", false, "Enable debug output")
	logPath := flag.String("log", "", "Write render logs to the given file path")

	flag.Parse()

	spec := examples.ExampleSplit()
	config := keel.NewConfig()
	config.SetDebug(*debug)
	renderer := keel.NewRendererWithConfig(
		config,
		spec,
		examples.ExampleSplitStyleProvider,
		examples.ExampleSplitContentProvider,
	)
	size := keel.Size{Width: *width, Height: *height}

	if *logPath != "" {
		logger, file, err := logging.NewFileLoggerPath(*logPath)
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
		config.SetLogger(logger.Log)
	}

	rendered, err := renderer.Render(size)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(rendered)
}
