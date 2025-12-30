# Contributing

Thanks for contributing to keel. This repo uses `mise` for consistent tooling.

## Setup

1) Install Mise: https://mise.jdx.dev/getting-started.html#installing-mise-cli
2) Trust and install tools:

```sh
mise trust
mise install
```

## Common tasks

- `mise run demo` to run the example dashboard
- `mise run test` to run tests (no cache)
- `mise run bench` to run benchmarks
- `mise run precommit` to run fmt, vet, build, and tests
- `mise run bench-report` to update `current_bench_result.txt` and `BENCHMARKS.md`

## Pull requests

- Use Conventional Commits (e.g., `feat: add allocator`, `fix: handle zero sizes`).
- Run `mise run precommit` before opening a PR.
- For rendering changes, include a short before/after snippet in the PR description.
