package keel

import (
	"fmt"
	"strconv"

	gloss "github.com/charmbracelet/lipgloss"
)

// Render arranges and renders the stored spec at the given size.
func (r *Renderer[KID]) Render(size Size) (string, error) {
	if r == nil {
		return "", ErrRendererMissing
	}
	if r.spec == nil {
		return "", ErrSpecMissing
	}
	layout, err := r.ensureLayout(size)
	if err != nil {
		return "", err
	}
	return r.renderLayout(layout)
}

func (r *Renderer[KID]) ensureLayout(size Size) (Layout[KID], error) {
	if r.hasLayout && r.last == size {
		return r.layout, nil
	}
	layout, err := arrange(r, r.spec, size)
	if err != nil {
		return Layout[KID]{}, err
	}
	r.layout = layout
	r.last = size
	r.hasLayout = true
	return layout, nil
}

func (r *Renderer[KID]) renderLayout(layout Layout[KID]) (string, error) {
	path := ""
	logger := rendererLogger(r)
	if logger != nil {
		path = "/"
	}
	return renderLayoutWithPath(layout.Root, r, path)
}

func renderLayoutWithPath[KID KeelID](node LayoutNode[KID], r *Renderer[KID], path string) (string, error) {
	logger := rendererLogger(r)
	switch node.Kind {
	case NodeStack:
		if len(node.Slots) == 0 {
			return "", nil
		}
		axis := node.Axis
		if axis != AxisHorizontal && axis != AxisVertical {
			err := &ConfigError{Reason: ErrInvalidAxis}
			logError(logger, path, "stack.axis", err)
			return "", err
		}

		rendered := make([]string, len(node.Slots))
		for i, slot := range node.Slots {
			slotPath := path
			if logger != nil {
				slotPath = appendPath(path, i)
			}
			out, err := renderLayoutWithPath(slot, r, slotPath)
			if err != nil {
				logError(logger, path, "stack.render", err)
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
			logError(logger, path, "dispatch", err)
			return "", err
		}
		size := Size{Width: node.Rect.Width, Height: node.Rect.Height}
		return renderFrameWithPath(node.Frame, r, size, path)
	default:
		err := &ConfigError{Reason: ErrUnknownSpec}
		logError(logger, path, "dispatch", err)
		return "", err
	}
}

func renderFrameWithPath[KID KeelID](frame FrameSpec[KID], r *Renderer[KID], size Size, path string) (string, error) {
	logger := rendererLogger(r)
	providedStyle := styleFor(r, frame)

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
		logError(logger, path, "frame.frame", err)
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
		logError(logger, path, "frame.frame", err)
		return "", err
	}

	availableWidth := size.Width - frameWidth
	availableHeight := size.Height - frameHeight

	info := FrameInfo{
		Width:         size.Width,
		Height:        size.Height,
		ContentWidth:  availableWidth,
		ContentHeight: availableHeight,
		FrameWidth:    frameWidth,
		FrameHeight:   frameHeight,
		Fit:           frame.Fit(),
	}

	logf(
		logger,
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

	content, err := contentFor(r, frame.ID(), info)
	if err != nil {
		logError(logger, path, "frame.content", err)
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
			logError(logger, path, "frame.content", err)
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
			logError(logger, path, "frame.content", err)
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
			logError(logger, path, "frame.content", err)
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
			logError(logger, path, "frame.content", err)
			return "", err
		}
	case FitOverflow:
		// No fitting or validation; let lipgloss render freely.
	default:
		err := &ConfigError{}
		logError(logger, path, "frame.fit", err)
		return "", err
	}

	outerWidth := size.Width - marginWidth - borderWidth
	outerHeight := size.Height - marginHeight - borderHeight
	style = style.
		Width(outerWidth).
		Height(outerHeight)

	return style.Render(contentToRender), nil
}

func styleFor[KID KeelID](r *Renderer[KID], frame FrameSpec[KID]) *gloss.Style {
	if r == nil || r.style == nil {
		return nil
	}
	return r.style(frame.ID())
}

func contentFor[KID KeelID](r *Renderer[KID], id KID, info FrameInfo) (string, error) {
	ecp := effectiveContentProvider(r)
	if ecp == nil {
		return "", &ContentProviderMissingError{ID: id}
	}

	return ecp(id, info)
}

func effectiveContentProvider[KID KeelID](r *Renderer[KID]) ContentProvider[KID] {
	if r == nil {
		return nil
	}
	if r.config != nil && r.config.debug {
		return DefaultDebugProvider[KID]
	}
	return r.content
}

func rendererLogger[KID KeelID](r *Renderer[KID]) LoggerFunc {
	if r == nil || r.config == nil {
		return nil
	}
	return r.config.logger
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
