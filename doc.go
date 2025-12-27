// Package keel provides deterministic spatial layout for discrete character buffers.
//
// Rendering is top-down: each container splits its allocated space along an
// axis and passes the resulting width/height to its slots. Blocks render
// content and optional lipgloss styles inside that allocation. Rendering is
// strict: if frames or content (after clipping) do not fit, rendering fails
// with ExtentTooSmallError. Keel does not perform intrinsic measurement,
// constraint solving, or stateful rendering.
//
// Box model (used by blocks):
//
//	+---------------------------------+
//	|           Margin                |
//	|  +---------------------------+  |
//	|  |        Border             |  |
//	|  |  +---------------------+  |  |
//	|  |  |     Padding         |  |  |
//	|  |  |  +---------------+  |  |  |
//	|  |  |  |  Content      |  |  |  |
//	|  |  |  +---------------+  |  |  |
//	|  |  +---------------------+  |  |
//	|  +---------------------------+  |
//	+---------------------------------+
//
// Sizing rules:
//   - ExtentConstraint and Context.Width/Height describe the total allocation
//     (Content + Padding + Border + Margin) for a renderable along an axis.
//   - lipgloss.Style.Width/Height describe the inner box (Content + Padding),
//     excluding border and margins.
//   - lipgloss.Style.GetFrameSize returns Margin + Padding + Border.
//   - ContentProvider receives the block ID and RenderInfo with allocation and
//     content box sizes, plus any ClipConstraint applied to the block.
//   - ClipConstraint is an optional content-only cap; content is clipped first
//     and the clipped content must fit the content box.
//   - StyleProvider may return cached styles; the renderer copies them and
//     treats them as immutable.
package keel
