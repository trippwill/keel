package keel

import (
	"errors"
	"fmt"

	gloss "github.com/charmbracelet/lipgloss"
)

// Render renders a [Renderable] node given the rendering context.
func Render[KID KeelID](node Renderable, ctx Context[KID]) (string, error) {
	switch n := node.(type) {
	case Block[KID]:
		return RenderBlock(n, ctx)
	case Container:
		return RenderContainer(n, ctx)
	default:
		return "", &ConfigError{Reason: ErrUnknownRenderable}
	}
}

// RenderContainer renders a [Container] node given the rendering context.
func RenderContainer[KID KeelID](container Container, ctx Context[KID]) (string, error) {
	length := container.Len()
	axis := container.GetAxis()
	if axis != AxisHorizontal && axis != AxisVertical {
		return "", &ConfigError{Reason: ErrInvalidAxis}
	}

	extentAt := func(index int) (ExtentConstraint, error) {
		child, ok := container.Slot(index)
		if !ok || child == nil {
			return ExtentConstraint{}, &SlotError{Index: index, Reason: ErrNilSlot}
		}
		return child.GetExtent(), nil
	}

	rendered := make([]string, length)
	switch axis {
	case AxisHorizontal:
		widths, required, err := RowResolver(ctx.Width, length, extentAt)
		if err != nil {
			if errors.Is(err, ErrExtentTooSmall) {
				return "", &ExtentTooSmallError{
					Axis:   AxisHorizontal,
					Need:   required,
					Have:   ctx.Width,
					Source: "horizontal split",
					Reason: "allocation",
				}
			}
			return "", err
		}

		for i, width := range widths {
			child, ok := container.Slot(i)
			if !ok || child == nil {
				return "", &SlotError{Index: i, Reason: ErrNilSlot}
			}
			_ctx := ctx.WithWidth(width)
			rendered[i], err = Render(child, _ctx)
			if err != nil {
				return "", err
			}
		}
		return gloss.JoinHorizontal(gloss.Top, rendered...), nil

	case AxisVertical:
		heights, required, err := ColResolver(ctx.Height, length, extentAt)
		if err != nil {
			if errors.Is(err, ErrExtentTooSmall) {
				return "", &ExtentTooSmallError{
					Axis:   AxisVertical,
					Need:   required,
					Have:   ctx.Height,
					Source: "vertical split",
					Reason: "allocation",
				}
			}
			return "", err
		}

		for i, height := range heights {
			child, ok := container.Slot(i)
			if !ok || child == nil {
				return "", &SlotError{Index: i, Reason: ErrNilSlot}
			}
			_ctx := ctx.WithHeight(height)
			rendered[i], err = Render(child, _ctx)
			if err != nil {
				return "", err
			}
		}
		return gloss.JoinVertical(gloss.Left, rendered...), nil
	}

	return "", &ConfigError{Reason: ErrInvalidAxis}
}

// RenderBlock renders a [Block] node given the rendering context.
// It validates that frame and content (or clip) fit in the allocation.
// Styles are copied before mutation, so cached styles are safe to reuse.
func RenderBlock[KID KeelID](node Block[KID], ctx Context[KID]) (string, error) {
	style := styleFor(ctx, node)
	frameWidth, frameHeight := 0, 0
	marginWidth, marginHeight := 0, 0
	borderWidth, borderHeight := 0, 0
	var transform func(string) string
	if style != nil {
		frameWidth, frameHeight = style.GetFrameSize()
		marginWidth = style.GetHorizontalMargins()
		marginHeight = style.GetVerticalMargins()
		borderWidth = style.GetHorizontalBorderSize()
		borderHeight = style.GetVerticalBorderSize()
		transform = style.GetTransform()
	}

	if frameWidth > ctx.Width {
		return "", &ExtentTooSmallError{
			Axis:   AxisHorizontal,
			Need:   frameWidth,
			Have:   ctx.Width,
			Source: sourceFor(node),
			Reason: "frame",
		}
	}
	if frameHeight > ctx.Height {
		return "", &ExtentTooSmallError{
			Axis:   AxisVertical,
			Need:   frameHeight,
			Have:   ctx.Height,
			Source: sourceFor(node),
			Reason: "frame",
		}
	}

	var outerStyle gloss.Style
	if style == nil {
		outerStyle = gloss.NewStyle()
	} else {
		outerStyle = (*style)
	}

	availableWidth := ctx.Width - frameWidth
	availableHeight := ctx.Height - frameHeight

	info := RenderInfo[KID]{
		ID:            node.GetID(),
		Width:         ctx.Width,
		Height:        ctx.Height,
		ContentWidth:  availableWidth,
		ContentHeight: availableHeight,
		FrameWidth:    frameWidth,
		FrameHeight:   frameHeight,
		Clip:          node.GetClip(),
	}

	content, err := contentFor(ctx, info)
	if err != nil {
		return "", err
	}

	contentForMeasure := content
	if transform != nil {
		contentForMeasure = transform(contentForMeasure)
		outerStyle = outerStyle.UnsetTransform()
	}

	// Clip only affects content, not the frame.
	clip := info.Clip
	contentToRender := contentForMeasure
	if clip.Width > 0 || clip.Height > 0 {
		clipStyle := gloss.NewStyle()
		if clip.Width > 0 {
			clipStyle = clipStyle.MaxWidth(clip.Width)
		}
		if clip.Height > 0 {
			clipStyle = clipStyle.MaxHeight(clip.Height)
		}
		contentToRender = clipStyle.Render(contentForMeasure)
	}

	contentWidth, contentHeight := gloss.Size(contentToRender)
	if contentWidth > availableWidth {
		return "", &ExtentTooSmallError{
			Axis:   AxisHorizontal,
			Need:   frameWidth + contentWidth,
			Have:   ctx.Width,
			Source: sourceFor(node),
			Reason: "content",
		}
	}
	if contentHeight > availableHeight {
		return "", &ExtentTooSmallError{
			Axis:   AxisVertical,
			Need:   frameHeight + contentHeight,
			Have:   ctx.Height,
			Source: sourceFor(node),
			Reason: "content",
		}
	}

	outerWidth := ctx.Width - marginWidth - borderWidth
	outerHeight := ctx.Height - marginHeight - borderHeight
	outerStyle = outerStyle.
		Width(outerWidth).
		Height(outerHeight)

	return outerStyle.Render(contentToRender), nil
}

func styleFor[KID KeelID](ctx Context[KID], node Block[KID]) *gloss.Style {
	if ctx.StyleProvider == nil {
		return nil
	}
	return ctx.StyleProvider(node.GetID())
}

func contentFor[KID KeelID](ctx Context[KID], info RenderInfo[KID]) (string, error) {
	if ctx.ContentProvider == nil {
		return "", &ContentProviderMissingError{ID: info.ID}
	}

	return ctx.ContentProvider(info)
}

func sourceFor[KID KeelID](node Block[KID]) string {
	return fmt.Sprintf("block %v", node.GetID())
}
