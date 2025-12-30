// Package keel provides deterministic spatial layout for discrete character buffers.
//
// Rendering is top-down: each container splits its allocated space along an
// axis and passes the resulting width/height to its slots. Blocks render
// content and optional lipgloss styles inside that allocation. Rendering is
// strict by default: if frames or content do not fit, rendering fails with
// ExtentTooSmallError unless a [FitMode] permits fitting. Keel does not perform
// intrinsic measurement,
// constraint solving, or stateful rendering.
//
// For repeated renders, you can resolve a layout once with [Resolve] and
// render it later with [RenderResolved] to avoid re-running allocation.
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
//   - ExtentConstraint and the resolved size describe the total allocation
//     (Content + Padding + Border + Margin) for a renderable along an axis.
//   - Flex max caps are soft: if all flex slots hit their max and space remains,
//     the remainder is distributed ignoring max caps.
//   - lipgloss.Style.Width/Height describe the inner box (Content + Padding),
//     excluding border and margins.
//   - lipgloss.Style.GetFrameSize returns Margin + Padding + Border.
//   - ContentProvider receives the block ID and RenderInfo with allocation and
//     content box sizes, plus the block's FitMode.
//   - FitMode controls whether content is wrapped, clipped, or allowed to
//     overflow the content box before rendering.
//   - StyleProvider may return cached styles; the renderer copies them and
//     treats them as immutable.
package keel
