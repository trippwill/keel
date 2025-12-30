package keel

import (
	"fmt"
	"strconv"

	gloss "github.com/charmbracelet/lipgloss"
)

// Render renders a [Renderable] given the rendering context and size.
// It resolves allocations and renders the resulting tree.
func Render[KID KeelID](ctx Context[KID], renderable Renderable, size Size) (string, error) {
	resolved, err := Resolve[KID](ctx, renderable, size)
	if err != nil {
		return "", err
	}
	return RenderResolved(ctx, resolved)
}

// RenderContainer renders a [Container] given the rendering context and size.
// The container's extents are resolved before rendering.
func RenderContainer[KID KeelID](ctx Context[KID], container Container, size Size) (string, error) {
	resolved, err := Resolve[KID](ctx, container, size)
	if err != nil {
		return "", err
	}
	return RenderResolved(ctx, resolved)
}

// RenderResolved renders a resolved layout tree with the given context.
func RenderResolved[KID KeelID](ctx Context[KID], resolved Resolved[KID]) (string, error) {
	path := ""
	if ctx.Logger != nil {
		path = "/"
	}
	return renderResolvedWithPath(resolved.Root, ctx, path)
}

func renderResolvedWithPath[KID KeelID](node ResolvedNode[KID], ctx Context[KID], path string) (string, error) {
	switch node.Kind {
	case NodeContainer:
		if len(node.Children) == 0 {
			return "", nil
		}
		axis := node.Axis
		if axis != AxisHorizontal && axis != AxisVertical {
			err := &ConfigError{Reason: ErrInvalidAxis}
			logError(ctx.Logger, path, "container.axis", err)
			return "", err
		}

		rendered := make([]string, len(node.Children))
		for i, child := range node.Children {
			childPath := path
			if ctx.Logger != nil {
				childPath = appendPath(path, i)
			}
			out, err := renderResolvedWithPath(child, ctx, childPath)
			if err != nil {
				logError(ctx.Logger, path, "container.render", err)
				return "", err
			}
			rendered[i] = out
		}

		if axis == AxisHorizontal {
			return gloss.JoinHorizontal(gloss.Top, rendered...), nil
		}
		return gloss.JoinVertical(gloss.Left, rendered...), nil

	case NodeBlock:
		if node.Block == nil {
			err := &ConfigError{Reason: ErrUnknownRenderable}
			logError(ctx.Logger, path, "dispatch", err)
			return "", err
		}
		size := Size{Width: node.Rect.Width, Height: node.Rect.Height}
		return renderBlockWithPath(node.Block, ctx, size, path)
	default:
		err := &ConfigError{Reason: ErrUnknownRenderable}
		logError(ctx.Logger, path, "dispatch", err)
		return "", err
	}
}

// RenderBlock renders a [Block] given the rendering context and size.
// It validates that frame and content (after clipping) fit in the allocation.
// Styles are copied before mutation, so cached styles are safe to reuse.
func RenderBlock[KID KeelID](ctx Context[KID], block Block[KID], size Size) (string, error) {
	path := ""
	if ctx.Logger != nil {
		path = "/"
	}
	return renderBlockWithPath(block, ctx, size, path)
}

func renderBlockWithPath[KID KeelID](block Block[KID], ctx Context[KID], size Size, path string) (string, error) {
	providedStyle := styleFor(ctx, block)

	// Initialize to default values
	var (
		style                     gloss.Style
		frameWidth, frameHeight   int
		marginWidth, marginHeight int
		borderWidth, borderHeight int
		transform                 func(string) string
	)

	if providedStyle == nil {
		style = gloss.NewStyle()
	} else {
		style = (*providedStyle)
		frameWidth, frameHeight = style.GetFrameSize()
		marginWidth = style.GetHorizontalMargins()
		marginHeight = style.GetVerticalMargins()
		borderWidth = style.GetHorizontalBorderSize()
		borderHeight = style.GetVerticalBorderSize()
		transform = style.GetTransform()

	}

	if frameWidth > size.Width {
		err := &ExtentTooSmallError{
			Axis:   AxisHorizontal,
			Need:   frameWidth,
			Have:   size.Width,
			Source: sourceFor(block),
			Reason: "frame",
		}
		logError(ctx.Logger, path, "block.frame", err)
		return "", err
	}
	if frameHeight > size.Height {
		err := &ExtentTooSmallError{
			Axis:   AxisVertical,
			Need:   frameHeight,
			Have:   size.Height,
			Source: sourceFor(block),
			Reason: "frame",
		}
		logError(ctx.Logger, path, "block.frame", err)
		return "", err
	}

	availableWidth := size.Width - frameWidth
	availableHeight := size.Height - frameHeight

	info := RenderInfo{
		Width:         size.Width,
		Height:        size.Height,
		ContentWidth:  availableWidth,
		ContentHeight: availableHeight,
		FrameWidth:    frameWidth,
		FrameHeight:   frameHeight,
		Fit:           block.GetFit(),
	}

	logf(
		ctx.Logger,
		path,
		LogEventBlockRender,
		block.GetID(),
		info.Width,
		info.Height,
		info.FrameWidth,
		info.FrameHeight,
		info.ContentWidth,
		info.ContentHeight,
		info.Fit,
	)

	content, err := contentFor(ctx, block.GetID(), info)
	if err != nil {
		logError(ctx.Logger, path, "block.content", err)
		return "", err
	}

	contentForMeasure := content
	if transform != nil {
		contentForMeasure = transform(contentForMeasure)
		style = style.UnsetTransform()
	}

	contentToRender := contentForMeasure
	switch info.Fit {
	case FitClip:
		if availableWidth <= 0 || availableHeight <= 0 {
			contentToRender = ""
			break
		}
		contentToRender = gloss.NewStyle().
			MaxWidth(availableWidth).
			MaxHeight(availableHeight).
			Render(contentToRender)
	case FitWrapClip:
		if availableWidth <= 0 || availableHeight <= 0 {
			contentToRender = ""
			break
		}
		contentToRender = gloss.NewStyle().
			Width(availableWidth).
			MaxWidth(availableWidth).
			MaxHeight(availableHeight).
			Render(contentToRender)
	case FitWrapStrict:
		if availableWidth > 0 {
			contentToRender = gloss.NewStyle().
				Width(availableWidth).
				Render(contentToRender)
		}
		contentWidth, contentHeight := gloss.Size(contentToRender)
		if contentWidth > availableWidth {
			err := &ExtentTooSmallError{
				Axis:   AxisHorizontal,
				Need:   frameWidth + contentWidth,
				Have:   size.Width,
				Source: sourceFor(block),
				Reason: "content",
			}
			logError(ctx.Logger, path, "block.content", err)
			return "", err
		}
		if contentHeight > availableHeight {
			err := &ExtentTooSmallError{
				Axis:   AxisVertical,
				Need:   frameHeight + contentHeight,
				Have:   size.Height,
				Source: sourceFor(block),
				Reason: "content",
			}
			logError(ctx.Logger, path, "block.content", err)
			return "", err
		}
	case FitExact:
		contentWidth, contentHeight := gloss.Size(contentToRender)
		if contentWidth > availableWidth {
			err := &ExtentTooSmallError{
				Axis:   AxisHorizontal,
				Need:   frameWidth + contentWidth,
				Have:   size.Width,
				Source: sourceFor(block),
				Reason: "content",
			}
			logError(ctx.Logger, path, "block.content", err)
			return "", err
		}
		if contentHeight > availableHeight {
			err := &ExtentTooSmallError{
				Axis:   AxisVertical,
				Need:   frameHeight + contentHeight,
				Have:   size.Height,
				Source: sourceFor(block),
				Reason: "content",
			}
			logError(ctx.Logger, path, "block.content", err)
			return "", err
		}
	case FitOverflow:
		// No fitting or validation; let lipgloss render freely.
	default:
		err := &ConfigError{}
		logError(ctx.Logger, path, "block.fit", err)
		return "", err
	}

	outerWidth := size.Width - marginWidth - borderWidth
	outerHeight := size.Height - marginHeight - borderHeight
	style = style.
		Width(outerWidth).
		Height(outerHeight)

	return style.Render(contentToRender), nil
}

func styleFor[KID KeelID](ctx Context[KID], block Block[KID]) *gloss.Style {
	if ctx.StyleProvider == nil {
		return nil
	}
	return ctx.StyleProvider(block.GetID())
}

func contentFor[KID KeelID](ctx Context[KID], id KID, info RenderInfo) (string, error) {
	if ctx.ContentProvider == nil {
		return "", &ContentProviderMissingError{ID: id}
	}

	return ctx.ContentProvider(id, info)
}

func sourceFor[KID KeelID](block Block[KID]) string {
	return fmt.Sprintf("block %v", block.GetID())
}

func logf(logger LoggerFunc, path string, event LogEvent, args ...any) {
	if logger == nil {
		return
	}
	msgFormat, ok := LogEventFormats[event]
	if !ok {
		msgFormat = "event=%v"
		args = []any{event}
	}
	logger(event, path, fmt.Sprintf(msgFormat, args...))
}

func logError(logger LoggerFunc, path string, stage string, err error) {
	if logger == nil || err == nil {
		return
	}
	logf(logger, path, LogEventRenderError, stage, err)
}

func appendPath(path string, index int) string {
	if path == "/" {
		return "/" + strconv.Itoa(index)
	}
	return path + "/" + strconv.Itoa(index)
}
