package keel

import (
	"errors"
	"fmt"

	gloss "github.com/charmbracelet/lipgloss"
)

// Render renders a [Renderable] given the rendering context.
// It dispatches to RenderContainer or RenderBlock based on the concrete type.
func Render[KID KeelID](renderable Renderable, ctx Context[KID]) (string, error) {
	switch n := renderable.(type) {
	case Container:
		return RenderContainer(n, ctx)
	case Block[KID]:
		return RenderBlock(n, ctx)
	default:
		return "", &ConfigError{Reason: ErrUnknownRenderable}
	}
}

// RenderContainer renders a [Container] given the rendering context.
// The container's extent is split across slots using the resolver rules.
// Allocation failures are reported as ExtentTooSmallError.
func RenderContainer[KID KeelID](container Container, ctx Context[KID]) (string, error) {
	length := container.Len()
	if length <= 0 {
		return "", nil
	}

	axis := container.GetAxis()
	if axis != AxisHorizontal && axis != AxisVertical {
		return "", &ConfigError{Reason: ErrInvalidAxis}
	}

	rendered := make([]string, length)

	switch axis {
	case AxisHorizontal:
		widths, required, err := RowResolver(ctx.Width, container)
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
			slot, ok := container.Slot(i)
			if !ok || slot == nil {
				return "", &SlotError{Index: i, Reason: ErrNilSlot}
			}
			_ctx := ctx.WithWidth(width)
			rendered[i], err = Render(slot, _ctx)
			if err != nil {
				return "", err
			}
		}
		return gloss.JoinHorizontal(gloss.Top, rendered...), nil

	case AxisVertical:
		heights, required, err := ColResolver(ctx.Height, container)
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
			slot, ok := container.Slot(i)
			if !ok || slot == nil {
				return "", &SlotError{Index: i, Reason: ErrNilSlot}
			}
			_ctx := ctx.WithHeight(height)
			rendered[i], err = Render(slot, _ctx)
			if err != nil {
				return "", err
			}
		}
		return gloss.JoinVertical(gloss.Left, rendered...), nil
	}

	return "", &ConfigError{Reason: ErrInvalidAxis}
}

// RenderBlock renders a [Block] given the rendering context.
// It validates that frame and content (after clipping) fit in the allocation.
// Styles are copied before mutation, so cached styles are safe to reuse.
func RenderBlock[KID KeelID](block Block[KID], ctx Context[KID]) (string, error) {
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
		return "", &ExtentTooSmallError{
			Axis:   AxisHorizontal,
			Need:   frameWidth,
			Have:   ctx.Width,
			Source: sourceFor(block),
			Reason: "frame",
		}
	}
	if frameHeight > ctx.Height {
		return "", &ExtentTooSmallError{
			Axis:   AxisVertical,
			Need:   frameHeight,
			Have:   ctx.Height,
			Source: sourceFor(block),
			Reason: "frame",
		}
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
		Clip:          block.GetClip(),
	}

	content, err := contentFor(ctx, block.GetID(), info)
	if err != nil {
		return "", err
	}

	contentForMeasure := content
	if transform != nil {
		contentForMeasure = transform(contentForMeasure)
		style = style.UnsetTransform()
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
			Source: sourceFor(block),
			Reason: "content",
		}
	}
	if contentHeight > availableHeight {
		return "", &ExtentTooSmallError{
			Axis:   AxisVertical,
			Need:   frameHeight + contentHeight,
			Have:   ctx.Height,
			Source: sourceFor(block),
			Reason: "content",
		}
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
