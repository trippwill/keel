# Keel

Deterministic spatial layout for discrete character buffers

Keel is a deterministic layout engine for terminal applications. You describe a
layout hierarchy, Keel deterministically allocates space along rows and
columns, and frames render content and optional lipgloss styles. Rendering
is strict by default: if frames or content don't fit the allocation, Keel
returns an `ExtentTooSmallError` unless a `FitMode` permits fitting.

## Concepts

- `Row` / `Col` define stacks that split space along an axis.
- `Panel` is a frame identified by a `KeelID`.
- `ExtentConstraint` (`Fixed`, `Flex`, `FlexMin`, `FlexMax`, `FlexMinMax`) controls how space is
  allocated along the stack axis.
- `Size` describes the available width/height for arrange/render.
- `FitMode` controls how content fits inside a frame.
- `Context` provides `ContentProvider`, `StyleProvider`, and logging.
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
		keel.Panel(keel.Fixed(3), "header"),
		keel.Row(keel.FlexUnit(),
			keel.Panel(keel.FlexMin(1, 10), "nav"),
			keel.Panel(keel.FlexMin(2, 20), "body"),
		),
	)

	ctx := keel.Context[string]{
		ContentProvider: func(id string, _ keel.RenderInfo) (string, error) {
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
		StyleProvider: func(id string) *gloss.Style {
			if id == "header" {
				style := gloss.NewStyle().Bold(true).Padding(0, 1)
				return &style
			}
			return nil
		},
	}
	size := keel.Size{Width: 80, Height: 24}

	out, err := keel.RenderSpec(ctx, layout, size)
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
	keel.Panel(keel.FlexMinMax(1, 10, 20), "nav"),
	keel.Panel(keel.FlexMax(2, 30), "body"),
)
```

## Arranged layouts

If you render repeatedly at the same size, arrange once and re-use the arranged
tree until the width/height or layout changes.

```go
size := keel.Size{Width: 80, Height: 24}
layout, err := keel.Arrange(ctx, layout, size)
if err != nil {
	panic(err)
}

out, err := keel.Render(ctx, layout)
if err != nil {
	panic(err)
}
```

## Logging

Keel can emit render logs through a context logger. Log events include stack
allocations, frame renders, and render errors. Paths are slash-delimited slot
indices rooted at `/` (e.g. `/0/1`).

```go
logger := keel.NewFileLogger(os.Stdout)
ctx := keel.Context[string]{}.
	WithContentProvider(contentProvider).
	WithStyleProvider(styleProvider).
	WithLogger(logger.Log)
size := keel.Size{Width: 80, Height: 24}

out, err := keel.RenderSpec(ctx, layout, size)
```

The default message formats are available via `keel.LogEventFormats`, and you
can supply any `LoggerFunc` to integrate with your own logging.

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
