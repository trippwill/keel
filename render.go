package keel

import (
	"errors"
	"fmt"
	"strconv"

	gloss "github.com/charmbracelet/lipgloss"
)

// Render renders a [Renderable] given the rendering context.
// It dispatches to RenderContainer or RenderBlock based on the concrete type.
func Render[KID KeelID](renderable Renderable, ctx Context[KID]) (string, error) {
	path := ""
	if ctx.Logger != nil {
		path = "/"
	}
	return renderWithPath(renderable, ctx, path)
}

func renderWithPath[KID KeelID](renderable Renderable, ctx Context[KID], path string) (string, error) {
	switch n := renderable.(type) {
	case Container:
		return renderContainerWithPath(n, ctx, path)
	case Block[KID]:
		return renderBlockWithPath(n, ctx, path)
	default:
		err := &ConfigError{Reason: ErrUnknownRenderable}
		logError(ctx.Logger, path, "dispatch", err)
		return "", err
	}
}

// RenderContainer renders a [Container] given the rendering context.
// The container's extent is split across slots using the resolver rules.
// Allocation failures are reported as ExtentTooSmallError.
func RenderContainer[KID KeelID](container Container, ctx Context[KID]) (string, error) {
	path := ""
	if ctx.Logger != nil {
		path = "/"
	}
	return renderContainerWithPath(container, ctx, path)
}

func renderContainerWithPath[KID KeelID](container Container, ctx Context[KID], path string) (string, error) {
	length := container.Len()
	if length <= 0 {
		return "", nil
	}

	axis := container.GetAxis()
	if axis != AxisHorizontal && axis != AxisVertical {
		err := &ConfigError{Reason: ErrInvalidAxis}
		logError(ctx.Logger, path, "container.axis", err)
		return "", err
	}

	rendered := make([]string, length)

	switch axis {
	case AxisHorizontal:
		widths, required, err := RowResolver(ctx.Width, container)
		if err != nil {
			if errors.Is(err, ErrExtentTooSmall) {
				err = &ExtentTooSmallError{
					Axis:   AxisHorizontal,
					Need:   required,
					Have:   ctx.Width,
					Source: "horizontal split",
					Reason: "allocation",
				}
			}
			logError(ctx.Logger, path, "container.resolve", err)
			return "", err
		}
		logf(
			ctx.Logger,
			path,
			LogEventContainerAlloc,
			axis.String(),
			ctx.Width,
			len(widths),
			widths,
			required,
		)

		for i, width := range widths {
			slot, ok := container.Slot(i)
			if !ok || slot == nil {
				err := &SlotError{Index: i, Reason: ErrNilSlot}
				logError(ctx.Logger, path, "container.slot", err)
				return "", err
			}
			_ctx := ctx.WithWidth(width)
			childPath := path
			if ctx.Logger != nil {
				childPath = appendPath(path, i)
			}
			rendered[i], err = renderWithPath(slot, _ctx, childPath)
			if err != nil {
				logError(ctx.Logger, path, "container.render", err)
				return "", err
			}
		}
		return gloss.JoinHorizontal(gloss.Top, rendered...), nil

	case AxisVertical:
		heights, required, err := ColResolver(ctx.Height, container)
		if err != nil {
			if errors.Is(err, ErrExtentTooSmall) {
				err = &ExtentTooSmallError{
					Axis:   AxisVertical,
					Need:   required,
					Have:   ctx.Height,
					Source: "vertical split",
					Reason: "allocation",
				}
			}
			logError(ctx.Logger, path, "container.resolve", err)
			return "", err
		}
		logf(
			ctx.Logger,
			path,
			LogEventContainerAlloc,
			axis.String(),
			ctx.Height,
			len(heights),
			heights,
			required,
		)

		for i, height := range heights {
			slot, ok := container.Slot(i)
			if !ok || slot == nil {
				err := &SlotError{Index: i, Reason: ErrNilSlot}
				logError(ctx.Logger, path, "container.slot", err)
				return "", err
			}
			_ctx := ctx.WithHeight(height)
			childPath := path
			if ctx.Logger != nil {
				childPath = appendPath(path, i)
			}
			rendered[i], err = renderWithPath(slot, _ctx, childPath)
			if err != nil {
				logError(ctx.Logger, path, "container.render", err)
				return "", err
			}
		}
		return gloss.JoinVertical(gloss.Left, rendered...), nil
	}

	err := &ConfigError{Reason: ErrInvalidAxis}
	logError(ctx.Logger, path, "container.axis", err)
	return "", err
}

// RenderBlock renders a [Block] given the rendering context.
// It validates that frame and content (after clipping) fit in the allocation.
// Styles are copied before mutation, so cached styles are safe to reuse.
func RenderBlock[KID KeelID](block Block[KID], ctx Context[KID]) (string, error) {
	path := ""
	if ctx.Logger != nil {
		path = "/"
	}
	return renderBlockWithPath(block, ctx, path)
}

func renderBlockWithPath[KID KeelID](block Block[KID], ctx Context[KID], path string) (string, error) {
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

	if frameWidth > ctx.Width {
		err := &ExtentTooSmallError{
			Axis:   AxisHorizontal,
			Need:   frameWidth,
			Have:   ctx.Width,
			Source: sourceFor(block),
			Reason: "frame",
		}
		logError(ctx.Logger, path, "block.frame", err)
		return "", err
	}
	if frameHeight > ctx.Height {
		err := &ExtentTooSmallError{
			Axis:   AxisVertical,
			Need:   frameHeight,
			Have:   ctx.Height,
			Source: sourceFor(block),
			Reason: "frame",
		}
		logError(ctx.Logger, path, "block.frame", err)
		return "", err
	}

	availableWidth := ctx.Width - frameWidth
	availableHeight := ctx.Height - frameHeight

	info := RenderInfo{
		Width:         ctx.Width,
		Height:        ctx.Height,
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
				Have:   ctx.Width,
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
				Have:   ctx.Height,
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
				Have:   ctx.Width,
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
				Have:   ctx.Height,
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

	outerWidth := ctx.Width - marginWidth - borderWidth
	outerHeight := ctx.Height - marginHeight - borderHeight
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
