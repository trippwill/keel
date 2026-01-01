// Package keel provides deterministic spatial layout for discrete character buffers.
//
// Rendering is top-down: each stack splits its allocated space along an
// axis and passes the resulting width/height to its slots. Frames render
// content and optional lipgloss styles inside that allocation. Rendering is
// strict by default: if frames or content do not fit, rendering fails with
// [ExtentTooSmallError] unless a [FitMode] permits fitting. Keel does not perform
// intrinsic measurement,
// constraint solving, or stateful rendering.
//
// For repeated renders, store a spec on a [Renderer] and call [Renderer.Render].
// The renderer caches the arranged layout for the last size; call [Renderer.Invalidate]
// after mutating a spec.
//
// Box model (used by frames):
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
//   - [ExtentConstraint] and the arranged size describe the total allocation
//     (Content + Padding + Border + Margin) for a [Spec] along an axis.
//   - Flex max caps are soft: if all flex slots hit their max and space remains,
//     the remainder is distributed ignoring max caps.
//   - lipgloss.Style.Width/Height describe the inner box (Content + Padding),
//     excluding border and margins.
//   - lipgloss.Style.GetFrameSize returns Margin + Padding + Border.
//   - [ContentProvider] receives the frame ID and [FrameInfo] with allocation and
//     content box sizes, plus the frame's [FitMode].
//   - FitMode controls whether content is wrapped, clipped, or allowed to
//     overflow the content box before rendering.
//   - [StyleProvider] may return cached styles; the renderer copies them and
//     treats them as immutable.
package keel
