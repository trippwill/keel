# Keel

Deterministic spatial layout for discrete character buffers

Keel is a deterministic layout engine for terminal applications. You describe a
layout hierarchy, Keel deterministically allocates space along rows and
columns, and panels render content and optional lipgloss styles. Rendering
is strict: if frames or content (after clipping) don't fit the allocation, Keel
returns an `ExtentTooSmallError`.

## Concepts

- `Row` / `Col` define containers that split space along an axis.
- `Panel` is a block identified by a `KeelID`.
- `ExtentConstraint` (`Fixed`, `Flex`, `FlexMin`) controls how space is
  allocated along the container axis.
- `ClipConstraint` caps content size before rendering.
- `Context` provides size, `ContentProvider`, and `StyleProvider`.

## Example

```go
package main

import (
	"fmt"

	gloss "github.com/charmbracelet/lipgloss"
	"github.com/trippwill/chiplog/keel"
)

func main() {
	layout := keel.Col(keel.FlexUnit(),
		keel.Panel(keel.Fixed(3), "header"),
		keel.Row(keel.FlexUnit(),
			keel.Panel(keel.FlexMin(1, 10), "nav"),
			keel.Panel(keel.FlexMin(2, 20), "body"),
		),
	)

	ctx := keel.Context[string]{
		Width:  80,
		Height: 24,
		ContentProvider: func(id string, _ keel.RenderInfo) (string, error) {
			switch id {
			case "header":
				return "Dashboard", nil
			case "nav":
				return "nav", nil
			case "body":
				return "content", nil
			default:
				return "", &keel.UnknownBlockIDError{ID: id}
			}
		},
		StyleProvider: func(id string) *gloss.Style {
			if id == "header" {
				style := gloss.NewStyle().Bold(true).Padding(0, 1)
				return &style
			}
			return nil
		},
	}

	out, err := keel.Render(layout, ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(out)
}
```

There is a runnable demo in `examples/dashboard` that uses the shared fixtures
in `examples`.

## Logging

Keel can emit render logs through a context logger. Log events include container
allocations, block renders, and render errors. Paths are slash-delimited slot
indices rooted at `/` (e.g. `/0/1`).

```go
logger := keel.NewFileLogger(os.Stdout)
ctx := keel.Context[string]{}.
	WithSize(80, 24).
	WithContentProvider(contentProvider).
	WithStyleProvider(styleProvider).
	WithLogger(logger.Log)

out, err := keel.Render(layout, ctx)
```

The default message formats are available via `keel.LogEventFormats`, and you
can supply any `LoggerFunc` to integrate with your own logging.

## Limitations

Keel does not perform intrinsic measurement, constraint solving, or stateful
rendering. It exists solely to map hierarchical layout intent onto terminal
geometry.

1. No intrinsic sizing
   - Panels don't ask "how big do you want to be?"

2. No focus / input model
   - Keel never knows about: cursor, focus, keybindings
   - Those go in the engine layer, keyed by KeelID

## Development

- `mise install` to set up a new development environment
  - May require network access to fetch dependencies
- `mise run test-keel` to run tests
- `mise run bench-keel` to run benchmarks
- `mise run precommit-keel` to run pre-commit checks
  - Always run this before opening a PR
  - May require network access to sync the workspace
  - Runs `go fmt`, `go vet`, `go test`, and `go generate`

Or standard Go commands:

- `go test ./...`
- `go test ./... -bench=BenchmarkRenderExampleSplit -benchmem`
- `go generate ./...`
