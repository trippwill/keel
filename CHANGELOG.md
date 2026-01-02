# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

- Initial split from chiplog.
- Breaking: rendering now goes through `Renderer` with stored specs and cached layouts; `Arrange`/`RenderSpec`/`RenderStackSpec` were removed from the public API.
- Breaking: core interfaces renamed to `Spec`/`FrameSpec`/`StackSpec` with method renames (`Extent`, `ID`, `Fit`, `Axis`, `Slot`).
- Breaking: allocator APIs renamed (`RowResolver`/`ColResolver` → `ArrangeStack`, `ResolveExtents` → `ArrangeExtents`).
- `Size` moved into `geom.go`.
- Breaking: simplified error surface; config issues now return `SpecError` (wrapping `ErrConfigurationInvalid`) and size failures return `ExtentTooSmallError` with string axes.
