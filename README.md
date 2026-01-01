# Keel

Deterministic spatial layout for discrete character buffers

Keel is a deterministic layout engine for terminal applications. You describe a
layout hierarchy, Keel deterministically allocates space along rows and
columns, and frames render content and optional lipgloss styles. Rendering
is strict by default: if frames or content don't fit the allocation, Keel
returns an `ExtentTooSmallError` unless a `FitMode` permits fitting.

## Concepts

- `Row` / `Col` define stacks that split space along an axis.
- Frame constructors (`Exact`, `Clip`, `Wrap`, `WrapStrict`, `Overflow`) identify frames by `KeelID`.
- `ExtentConstraint` (`Fixed`, `Flex`, `FlexMin`, `FlexMax`, `FlexMinMax`) controls how space is
  allocated along the stack axis.
- `Size` describes the available width/height for arrange/render.
- Fit modes (`Exact`, `Clip`, `Wrap`, `WrapStrict`, `Overflow`) control how content fits inside a frame.
- `Renderer` provides `ContentProvider`, `StyleProvider`, and render configuration.
- Flex max caps are soft: if all flex slots hit their max and space remains,
  the remainder is distributed ignoring max caps.

## Example

```go
package main

import (
	"fmt"

	gloss "github.com/charmbracelet/lipgloss"
	"github.com/trippwill/keel"
)

func main() {
	layout := keel.Col(keel.FlexUnit(),
		keel.Exact(keel.Fixed(3), "header"),
		keel.Row(keel.FlexUnit(),
			keel.Exact(keel.FlexMin(1, 10), "nav"),
			keel.Exact(keel.FlexMin(2, 20), "body"),
		),
	)

	renderer := keel.NewRenderer(
		layout,
		func(id string) *gloss.Style {
			if id == "header" {
				style := gloss.NewStyle().Bold(true).Padding(0, 1)
				return &style
			}
			return nil
		},
		func(id string, _ keel.FrameInfo) (string, error) {
			switch id {
			case "header":
				return "Dashboard", nil
			case "nav":
				return "nav", nil
			case "body":
				return "content", nil
			default:
				return "", &keel.UnknownFrameIDError{ID: id}
			}
		},
	)
	size := keel.Size{Width: 80, Height: 24}

	out, err := renderer.Render(size)
	if err != nil {
		panic(err)
	}

	fmt.Println(out)
}
```

There is a runnable demo in `examples/dashboard` that uses the shared fixtures
in `examples`.

Here's a small example using soft max caps:

```go
layout := keel.Row(keel.FlexUnit(),
	keel.Exact(keel.FlexMinMax(1, 10, 20), "nav"),
	keel.Exact(keel.FlexMax(2, 30), "body"),
)
```

## Arranged layouts

Renderers cache the arranged layout for the last size. Call `Render` with the
current size; it will re-arrange only when the size changes. If you mutate a
spec in place, call `renderer.Invalidate()` to force a re-arrange. For a new
spec, construct a new renderer.

```go
size := keel.Size{Width: 80, Height: 24}
out, err := renderer.Render(size)
if err != nil {
	panic(err)
}
```

## Logging

Keel can emit render logs through the renderer config logger. Log events include stack
allocations, frame renders, and render errors. Paths are slash-delimited slot
indices rooted at `/` (e.g. `/0/1`).

```go
// import "github.com/trippwill/keel/logging"
logger := logging.NewFileLogger(os.Stdout)
renderer := keel.NewRenderer(layout, styleProvider, contentProvider)
renderer.Config().SetLogger(logger.Log)
size := keel.Size{Width: 80, Height: 24}

out, err := renderer.Render(size)
```

The default message formats are available via `logging.LogEventFormats`, and you
can supply any `logging.LoggerFunc` to integrate with your own logging.

## Limitations

Keel does not perform intrinsic measurement, or stateful
rendering. It exists solely to map hierarchical layout intent onto terminal
geometry.

1. No intrinsic sizing
   - Frames don't ask "how big do you want to be?"

2. No focus / input model
   - Keel never knows about: cursor, focus, keybindings
   - Those go in the engine layer, keyed by KeelID

## Development

**One time setup**
- install [Mise](https://mise.jdx.dev/getting-started.html#installing-mise-cli)
- `mise trust`
- `mise install` to set up a new development environment
  - May require network access to fetch dependencies

**Development commands**

- `mise run demo` to run the example dashboard
- `mise run test` to run tests (no cache)
- `mise run bench` to run benchmarks
- `mise run precommit` to run fmt, vet, build, and tests
- `mise run bench-report` to update `current_bench_result.txt` and `BENCHMARKS.md`

Or standard Go commands:

- `go test ./...`
- `go test ./... -bench='BenchmarkRender|BenchmarkArrange' -benchmem`
- `go generate ./...`
- `go fmt ./...`
- `go vet ./...`
- `go build ./...`
