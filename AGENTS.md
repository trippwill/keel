# Keel Package Guidelines

## Scope & Structure
- This repo is a standalone Go module for the layout/render engine `keel` (`go.mod`).
- Core shared types and errors (including debug providers) live in `core/`. Advanced layout primitives and allocation live in `engine/`. Renderer and helpers live in `render.go`, `renderer.go`, `frame.go`, `stack.go`, `extent.go`, and `logging/`.
- Examples and shared fixtures live in the `examples` subpackage; a runnable demo is in `examples/dashboard`.

## Build, Test, and Development Commands
- `go test ./...` runs unit tests for the module.
- `go test ./... -bench='BenchmarkRender|BenchmarkArrange' -benchmem` runs benchmarks.
- `go generate ./...` regenerates stringer output for enums (see `geom.go`).
- `mise run demo` runs the example dashboard.
- `mise run test` runs tests with cache disabled.
- `mise run bench` runs the benchmark task.
- `mise run bench-report` updates `current_bench_result.txt` and `BENCHMARKS.md`.
- `mise run precommit` runs fmt, vet, build, and tests. If it hangs, rerun with elevated permissions.
- prefer using `mise` for consistent environment setup.

## Coding Style & Naming Conventions
- Follow `gofmt` output and standard Go conventions; exported symbols must have doc comments.
- Keep error strings lowercase and concise; use `SpecError` for configuration issues and `ExtentTooSmallError` for sizing failures (both in `keel/err.go`).
- Treat styles from `StyleProvider` as immutable; cached styles are expected.
- `ContentProvider` receives the frame ID plus `FrameInfo`; ensure content respects the content box and the frame's fit mode.
- Prefer small, composable helpers for allocation and rendering steps.

## Testing Guidelines
- Use `_test.go` files in this module with `package keel` or `package keel_test` as needed; `keel_test` is preferred when internals are not required.
- `examples.ExampleSplit` is the shared fixture for render tests and benchmarks.
- Include edge cases for extent allocation, frame/content sizing, clipping, and debug output.

## Commit & PR Notes
- Always run `mise run precommit` before committing.
- Use Conventional Commits (e.g., `feat(keel): add clip constraint` or `fix(keel): handle extent min`).
- For rendering changes, include a short before/after terminal output snippet in the PR description.
