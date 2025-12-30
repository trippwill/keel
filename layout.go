package keel

import "errors"

// NodeKind identifies the kind of node in an arranged layout tree.
type NodeKind uint8

const (
	// NodeStack represents a stack with slot allocations.
	NodeStack NodeKind = iota
	// NodeFrame represents a frame that renders content.
	NodeFrame
)

// Rect describes an allocated rectangle in the render space.
type Rect struct {
	X, Y          int
	Width, Height int
}

// Layout holds a layout tree arranged to concrete allocations.
type Layout[KID KeelID] struct {
	Width, Height int
	Root          LayoutNode[KID]
}

// LayoutNode represents an arranged layout node.
type LayoutNode[KID KeelID] struct {
	Kind  NodeKind
	Axis  Axis
	Rect  Rect
	Frame FrameSpec[KID]
	Slots []LayoutNode[KID]
}

// Arrange arranges a [Spec] tree into concrete allocations for the given size.
func Arrange[KID KeelID](ctx Context[KID], spec Spec, size Size) (Layout[KID], error) {
	path := ""
	if ctx.Logger != nil {
		path = "/"
	}
	rect := Rect{X: 0, Y: 0, Width: size.Width, Height: size.Height}
	root, err := arrangeWithPath[KID](spec, rect, path, ctx.Logger)
	if err != nil {
		return Layout[KID]{}, err
	}
	return Layout[KID]{
		Width:  size.Width,
		Height: size.Height,
		Root:   root,
	}, nil
}

func arrangeWithPath[KID KeelID](spec Spec, rect Rect, path string, logger LoggerFunc) (LayoutNode[KID], error) {
	switch n := spec.(type) {
	case StackSpec:
		return arrangeStackWithPath[KID](n, rect, path, logger)
	case FrameSpec[KID]:
		return LayoutNode[KID]{
			Kind:  NodeFrame,
			Rect:  rect,
			Frame: n,
		}, nil
	default:
		err := &ConfigError{Reason: ErrUnknownSpec}
		logError(logger, path, "dispatch", err)
		return LayoutNode[KID]{}, err
	}
}

func arrangeStackWithPath[KID KeelID](stack StackSpec, rect Rect, path string, logger LoggerFunc) (LayoutNode[KID], error) {
	length := stack.Len()
	if length <= 0 {
		return LayoutNode[KID]{
			Kind:  NodeStack,
			Rect:  rect,
			Slots: nil,
		}, nil
	}

	axis := stack.Axis()
	if axis != AxisHorizontal && axis != AxisVertical {
		err := &ConfigError{Reason: ErrInvalidAxis}
		logError(logger, path, "stack.axis", err)
		return LayoutNode[KID]{}, err
	}

	extents, err := GetStackExtents(stack)
	if err != nil {
		logError(logger, path, "stack.slot", err)
		return LayoutNode[KID]{}, err
	}

	total := rect.Width
	if axis == AxisVertical {
		total = rect.Height
	}

	sizes, required, err := ArrangeExtents(total, extents)
	if err != nil {
		if errors.Is(err, ErrExtentTooSmall) {
			source := "horizontal split"
			if axis == AxisVertical {
				source = "vertical split"
			}
			err = &ExtentTooSmallError{
				Axis:   axis,
				Need:   required,
				Have:   total,
				Source: source,
				Reason: "allocation",
			}
		}
		logError(logger, path, "stack.arrange", err)
		return LayoutNode[KID]{}, err
	}

	logf(
		logger,
		path,
		LogEventStackAlloc,
		axis.String(),
		total,
		len(sizes),
		sizes,
		required,
	)

	slots := make([]LayoutNode[KID], length)
	offset := 0
	for i, size := range sizes {
		slot, ok := stack.Slot(i)
		if !ok || slot == nil {
			err := &SlotError{Index: i, Reason: ErrNilSlot}
			logError(logger, path, "stack.slot", err)
			return LayoutNode[KID]{}, err
		}

		slotRect := rect
		if axis == AxisHorizontal {
			slotRect.X += offset
			slotRect.Width = size
		} else {
			slotRect.Y += offset
			slotRect.Height = size
		}

		slotPath := path
		if logger != nil {
			slotPath = appendPath(path, i)
		}

		slotNode, err := arrangeWithPath[KID](slot, slotRect, slotPath, logger)
		if err != nil {
			logError(logger, path, "stack.render", err)
			return LayoutNode[KID]{}, err
		}

		slots[i] = slotNode
		offset += size
	}

	return LayoutNode[KID]{
		Kind:  NodeStack,
		Axis:  axis,
		Rect:  rect,
		Slots: slots,
	}, nil
}
