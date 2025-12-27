// Package keel is a minimal layout engine for terminal applications.
//
// Rendering is top-down: each container splits its allocated space along an
// axis and passes the resulting width/height to its children. Leaf nodes render
// content and optional lipgloss styles inside that allocation.
//
// Box model (used by leaf nodes):
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
//   - ContentProvider receives RenderInfo with allocation and content box sizes,
//     plus any ClipConstraint applied to the block.
//   - ClipConstraint is an optional content-only cap; content is clipped first
//     and the clipped content must fit the content box.
//   - StyleProvider may return cached styles; the renderer copies them and
//     treats them as immutable.
//
// Rendering is strict: if the frame doesn't fit, or content (or clip) exceeds
// the content box, rendering returns ExtentTooSmallError.
package keel
