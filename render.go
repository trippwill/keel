package keel

import (
	"fmt"
	"strconv"

	gloss "github.com/charmbracelet/lipgloss"
)

// RenderSpec renders a [Spec] given the rendering context and size.
// It arranges allocations and renders the resulting tree.
func RenderSpec[KID KeelID](ctx Context[KID], spec Spec, size Size) (string, error) {
	layout, err := Arrange(ctx, spec, size)
	if err != nil {
		return "", err
	}
	return Render(ctx, layout)
}

// RenderStackSpec renders a [StackSpec] given the rendering context and size.
// The stack's extents are arranged before rendering.
func RenderStackSpec[KID KeelID](ctx Context[KID], stack StackSpec, size Size) (string, error) {
	layout, err := Arrange(ctx, stack, size)
	if err != nil {
		return "", err
	}
	return Render(ctx, layout)
}

// Render renders an arranged [Layout] tree with the given context.
func Render[KID KeelID](ctx Context[KID], layout Layout[KID]) (string, error) {
	path := ""
	if ctx.Logger != nil {
		path = "/"
	}
	return renderLayoutWithPath(layout.Root, ctx, path)
}

func renderLayoutWithPath[KID KeelID](node LayoutNode[KID], ctx Context[KID], path string) (string, error) {
	switch node.Kind {
	case NodeStack:
		if len(node.Slots) == 0 {
			return "", nil
		}
		axis := node.Axis
		if axis != AxisHorizontal && axis != AxisVertical {
			err := &ConfigError{Reason: ErrInvalidAxis}
			logError(ctx.Logger, path, "stack.axis", err)
			return "", err
		}

		rendered := make([]string, len(node.Slots))
		for i, slot := range node.Slots {
			slotPath := path
			if ctx.Logger != nil {
				slotPath = appendPath(path, i)
			}
			out, err := renderLayoutWithPath(slot, ctx, slotPath)
			if err != nil {
				logError(ctx.Logger, path, "stack.render", err)
				return "", err
			}
			rendered[i] = out
		}

		if axis == AxisHorizontal {
			return gloss.JoinHorizontal(gloss.Top, rendered...), nil
		}
		return gloss.JoinVertical(gloss.Left, rendered...), nil

	case NodeFrame:
		if node.Frame == nil {
			err := &ConfigError{Reason: ErrUnknownSpec}
			logError(ctx.Logger, path, "dispatch", err)
			return "", err
		}
		size := Size{Width: node.Rect.Width, Height: node.Rect.Height}
		return renderFrameWithPath(node.Frame, ctx, size, path)
	default:
		err := &ConfigError{Reason: ErrUnknownSpec}
		logError(ctx.Logger, path, "dispatch", err)
		return "", err
	}
}

// RenderFrame renders a [FrameSpec] given the rendering context and size.
// It validates that frame and content (after clipping) fit in the allocation.
// Styles are copied before mutation, so cached styles are safe to reuse.
func RenderFrame[KID KeelID](ctx Context[KID], frame FrameSpec[KID], size Size) (string, error) {
	path := ""
	if ctx.Logger != nil {
		path = "/"
	}
	return renderFrameWithPath(frame, ctx, size, path)
}

func renderFrameWithPath[KID KeelID](frame FrameSpec[KID], ctx Context[KID], size Size, path string) (string, error) {
	providedStyle := styleFor(ctx, frame)

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
			Source: sourceFor(frame),
			Reason: "frame",
		}
		logError(ctx.Logger, path, "frame.frame", err)
		return "", err
	}
	if frameHeight > size.Height {
		err := &ExtentTooSmallError{
			Axis:   AxisVertical,
			Need:   frameHeight,
			Have:   size.Height,
			Source: sourceFor(frame),
			Reason: "frame",
		}
		logError(ctx.Logger, path, "frame.frame", err)
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
		Fit:           frame.Fit(),
	}

	logf(
		ctx.Logger,
		path,
		LogEventFrameRender,
		frame.ID(),
		info.Width,
		info.Height,
		info.FrameWidth,
		info.FrameHeight,
		info.ContentWidth,
		info.ContentHeight,
		info.Fit,
	)

	content, err := contentFor(ctx, frame.ID(), info)
	if err != nil {
		logError(ctx.Logger, path, "frame.content", err)
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
				Source: sourceFor(frame),
				Reason: "content",
			}
			logError(ctx.Logger, path, "frame.content", err)
			return "", err
		}
		if contentHeight > availableHeight {
			err := &ExtentTooSmallError{
				Axis:   AxisVertical,
				Need:   frameHeight + contentHeight,
				Have:   size.Height,
				Source: sourceFor(frame),
				Reason: "content",
			}
			logError(ctx.Logger, path, "frame.content", err)
			return "", err
		}
	case FitExact:
		contentWidth, contentHeight := gloss.Size(contentToRender)
		if contentWidth > availableWidth {
			err := &ExtentTooSmallError{
				Axis:   AxisHorizontal,
				Need:   frameWidth + contentWidth,
				Have:   size.Width,
				Source: sourceFor(frame),
				Reason: "content",
			}
			logError(ctx.Logger, path, "frame.content", err)
			return "", err
		}
		if contentHeight > availableHeight {
			err := &ExtentTooSmallError{
				Axis:   AxisVertical,
				Need:   frameHeight + contentHeight,
				Have:   size.Height,
				Source: sourceFor(frame),
				Reason: "content",
			}
			logError(ctx.Logger, path, "frame.content", err)
			return "", err
		}
	case FitOverflow:
		// No fitting or validation; let lipgloss render freely.
	default:
		err := &ConfigError{}
		logError(ctx.Logger, path, "frame.fit", err)
		return "", err
	}

	outerWidth := size.Width - marginWidth - borderWidth
	outerHeight := size.Height - marginHeight - borderHeight
	style = style.
		Width(outerWidth).
		Height(outerHeight)

	return style.Render(contentToRender), nil
}

func styleFor[KID KeelID](ctx Context[KID], frame FrameSpec[KID]) *gloss.Style {
	if ctx.StyleProvider == nil {
		return nil
	}
	return ctx.StyleProvider(frame.ID())
}

func contentFor[KID KeelID](ctx Context[KID], id KID, info RenderInfo) (string, error) {
	if ctx.ContentProvider == nil {
		return "", &ContentProviderMissingError{ID: id}
	}

	return ctx.ContentProvider(id, info)
}

func sourceFor[KID KeelID](frame FrameSpec[KID]) string {
	return fmt.Sprintf("frame %v", frame.ID())
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
