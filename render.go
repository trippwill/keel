package keel

import gloss "github.com/charmbracelet/lipgloss"

func styleFor[KID KeelID](ctx Context[KID], id KID, kind NodeKind) *gloss.Style {
	if ctx.StyleProvider == nil {
		return nil
	}
	return ctx.StyleProvider(id, kind)
}

func nodeKind[KID KeelID](node Renderable[KID]) NodeKind {
	if _, ok := node.(Container[KID]); ok {
		return NodeContainer
	}
	return NodeContent
}

func renderSplit[KID KeelID](s *SplitSpec[KID], ctx Context[KID]) (string, error) {
	if s == nil {
		return "", ErrConfigurationInvalid
	}

	axis := s.GetAxis()
	if axis != AxisHorizontal && axis != AxisVertical {
		return "", ErrConfigurationInvalid
	}

	containerStyle := styleFor(ctx, s.GetID(), NodeContainer)
	frameWidth, frameHeight := 0, 0
	if containerStyle != nil {
		frameWidth, frameHeight = containerStyle.GetFrameSize()
	}
	if frameWidth > ctx.Width {
		return "", &TargetTooSmallError{Axis: AxisHorizontal, Need: frameWidth, Have: ctx.Width}
	}
	if frameHeight > ctx.Height {
		return "", &TargetTooSmallError{Axis: AxisVertical, Need: frameHeight, Have: ctx.Height}
	}

	contentWidth := ctx.Width - frameWidth
	contentHeight := ctx.Height - frameHeight

	count := s.Len()
	if count == 0 {
		return "", ErrConfigurationInvalid
	}

	specs := make([]SizeSpec, count)
	nodes := make([]Renderable[KID], count)
	for i := range count {
		spec, node, ok := s.Slot(i)
		if !ok || node == nil {
			return "", ErrConfigurationInvalid
		}
		specs[i] = spec
		nodes[i] = node
	}

	var alloc Allocator
	switch axis {
	case AxisHorizontal:
		alloc = RowAllocator
	case AxisVertical:
		alloc = ColAllocator
	}

	if alloc == nil {
		return "", ErrAllocatorUnimplemented
	}

	total := contentWidth
	if axis == AxisVertical {
		total = contentHeight
	}

	sizes, err := alloc(total, specs)
	if err != nil {
		if err == ErrTargetTooSmall {
			return "", &TargetTooSmallError{Axis: axis, Need: requiredMin(specs), Have: total}
		}
		return "", err
	}
	if len(sizes) != count {
		return "", ErrConfigurationInvalid
	}

	rendered := make([]string, count)
	for i, node := range nodes {
		kind := nodeKind(node)
		style := styleFor(ctx, node.GetID(), kind)
		childFrameWidth, childFrameHeight := 0, 0
		if style != nil {
			childFrameWidth, childFrameHeight = style.GetFrameSize()
		}

		if axis == AxisHorizontal {
			if childFrameHeight > contentHeight {
				return "", &TargetTooSmallError{Axis: AxisVertical, Need: childFrameHeight, Have: contentHeight}
			}
			need := childFrameWidth + specs[i].ContentMin
			if sizes[i] < need {
				return "", &TargetTooSmallError{Axis: AxisHorizontal, Need: need, Have: sizes[i]}
			}
		} else {
			if childFrameWidth > contentWidth {
				return "", &TargetTooSmallError{Axis: AxisHorizontal, Need: childFrameWidth, Have: contentWidth}
			}
			need := childFrameHeight + specs[i].ContentMin
			if sizes[i] < need {
				return "", &TargetTooSmallError{Axis: AxisVertical, Need: need, Have: sizes[i]}
			}
		}

		var childCtx Context[KID]
		if axis == AxisHorizontal {
			childCtx = ctx.WithSize(sizes[i], contentHeight)
		} else {
			childCtx = ctx.WithSize(contentWidth, sizes[i])
		}

		out, err := node.Render(childCtx)
		if err != nil {
			return "", err
		}
		rendered[i] = out
	}

	var composed string
	if axis == AxisHorizontal {
		composed = gloss.JoinHorizontal(gloss.Top, rendered...)
	} else {
		composed = gloss.JoinVertical(gloss.Left, rendered...)
	}

	if containerStyle != nil {
		composed = containerStyle.Render(composed)
	}

	return composed, nil
}

func renderPanel[KID KeelID](p *PanelSpec[KID], ctx Context[KID]) (string, error) {
	if p == nil {
		return "", ErrConfigurationInvalid
	}
	if ctx.ContentProvider == nil {
		return "", ErrConfigurationInvalid
	}

	style := styleFor(ctx, p.GetID(), NodeContent)
	if style != nil {
		frameWidth, frameHeight := style.GetFrameSize()
		if frameWidth > ctx.Width {
			return "", &TargetTooSmallError{Axis: AxisHorizontal, Need: frameWidth, Have: ctx.Width}
		}
		if frameHeight > ctx.Height {
			return "", &TargetTooSmallError{Axis: AxisVertical, Need: frameHeight, Have: ctx.Height}
		}
	}

	content, err := ctx.ContentProvider(p.GetID())
	if err != nil {
		return "", err
	}

	if style != nil {
		content = style.Render(content)
	}

	return content, nil
}

func requiredMin(specs []SizeSpec) int {
	required := 0
	for _, spec := range specs {
		switch spec.Kind {
		case SizeFixed:
			required += spec.Units
		case SizeFlex:
			required += spec.ContentMin
		}
	}
	return required
}
