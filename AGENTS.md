# Keel Package Guidelines

## Scope & Structure
- `keel/` is a standalone Go module for the layout/render engine (`keel/go.mod`).
- Core APIs live in `renderable.go`, `panel.go`, `split.go`, `geom.go`, `resolver.go`, `render.go`, `debug.go`, and `err.go`.
- Examples and shared fixtures live in the `keel/examples` subpackage; a runnable demo is in `keel/examples/dashboard`.

## Build, Test, and Development Commands
- `go test ./keel/...` runs unit tests for the module.
- `go test ./keel/... -bench=BenchmarkRenderExampleSplit -benchmem` runs the benchmark.
- `go generate ./keel/...` regenerates stringer output for enums (see `keel/geom.go`).
- `mise run test-keel` runs module tests with repo-configured caches.
- `mise run precommit-keel` runs `go work sync`, generate, fmt, vet, build, and tests for `keel`.
- `mise run bench-keel` runs the benchmark using the shared task.
- `mise run bench-report` updates `keel/BENCHMARKS.md` with a single benchmark run.

## Coding Style & Naming Conventions
- Follow `gofmt` output and standard Go conventions; exported symbols must have doc comments.
- Keep error strings lowercase and concise; wrap with typed errors from `keel/err.go`.
- Treat styles from `StyleProvider` as immutable; cached styles are expected.
- `ContentProvider` receives `RenderInfo`; ensure content respects the content box and clip.
- Prefer small, composable helpers for allocation and rendering steps.

## Testing Guidelines
- Use `_test.go` files in `keel/` or `keel_test` as needed; `keel_test` is preferred when internals are not required.
- `keel/examples.ExampleSplit` is the shared fixture for render tests and benchmarks.
- Include edge cases for extent allocation, frame/content sizing, clipping, and debug output.

## Commit & PR Notes
- Use Conventional Commits (e.g., `feat(keel): add clip constraint` or `fix(keel): handle extent min`).
- For rendering changes, include a short before/after terminal output snippet in the PR description.
