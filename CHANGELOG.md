# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

- Initial split from chiplog.
- Breaking: resolve flow renamed (`Resolve` → `Arrange`, `Resolved` → `Layout`), and `Render` now renders a `Layout`; use `RenderSpec`/`RenderStackSpec` for specs.
- Breaking: core interfaces renamed to `Spec`/`FrameSpec`/`StackSpec` with method renames (`Extent`, `ID`, `Fit`, `Axis`, `Slot`).
- Breaking: allocator APIs renamed (`RowResolver`/`ColResolver` → `ArrangeStack`, `ResolveExtents` → `ArrangeExtents`).
- `Size` moved into `geom.go`.
