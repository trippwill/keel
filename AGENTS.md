# Keel Package Guidelines

## Scope & Structure
- This directory is the standalone Go module for the layout/render engine (`go.mod`).
- Core APIs live in `renderable.go`, `panel.go`, `split.go`, `geom.go`, `resolver.go`, `render.go`, `debug.go`, and `err.go`.
- Examples and shared fixtures live in the `examples` subpackage; a runnable demo is in `examples/dashboard`.

## Build, Test, and Development Commands
- `go test ./...` runs unit tests for the module.
- `go test ./... -bench=BenchmarkRenderExampleSplit -benchmem` runs the benchmark.
- `go generate ./...` regenerates stringer output for enums (see `geom.go`).
- `mise run test-keel` runs module tests with repo-configured caches.
- `mise run precommit-keel` runs `go work sync`, generate, fmt, vet, build, and tests for `keel`.
- `mise run bench-keel` runs the benchmark using the shared task.
- `mise run bench-report` updates `BENCHMARKS.md` with a single benchmark run.

## Coding Style & Naming Conventions
- Follow `gofmt` output and standard Go conventions; exported symbols must have doc comments.
- Keep error strings lowercase and concise; wrap with typed errors from `keel/err.go`.
- Treat styles from `StyleProvider` as immutable; cached styles are expected.
- `ContentProvider` receives the block ID plus `RenderInfo`; ensure content respects the content box and the block's `FitMode`.
- Prefer small, composable helpers for allocation and rendering steps.

## Testing Guidelines
- Use `_test.go` files in this module or `keel_test` as needed; `keel_test` is preferred when internals are not required.
- `examples.ExampleSplit` is the shared fixture for render tests and benchmarks.
- Include edge cases for extent allocation, frame/content sizing, clipping, and debug output.

## Commit & PR Notes
- Use Conventional Commits (e.g., `feat(keel): add clip constraint` or `fix(keel): handle extent min`).
- For rendering changes, include a short before/after terminal output snippet in the PR description.
