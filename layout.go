package keel

import "errors"

// NodeKind identifies the kind of node in a resolved layout tree.
type NodeKind uint8

const (
	// NodeContainer represents a container with child allocations.
	NodeContainer NodeKind = iota
	// NodeBlock represents a block that renders content.
	NodeBlock
)

// Rect describes an allocated rectangle in the render space.
type Rect struct {
	X, Y          int
	Width, Height int
}

// Resolved holds a layout tree resolved to concrete allocations.
type Resolved[KID KeelID] struct {
	Width, Height int
	Root          ResolvedNode[KID]
}

// ResolvedNode represents a resolved layout node.
type ResolvedNode[KID KeelID] struct {
	Kind     NodeKind
	Axis     Axis
	Rect     Rect
	Block    Block[KID]
	Children []ResolvedNode[KID]
}

// Resolve resolves a layout tree into concrete allocations for the given size.
func Resolve[KID KeelID](ctx Context[KID], renderable Renderable, size Size) (Resolved[KID], error) {
	path := ""
	if ctx.Logger != nil {
		path = "/"
	}
	rect := Rect{X: 0, Y: 0, Width: size.Width, Height: size.Height}
	root, err := resolveWithPath[KID](renderable, rect, path, ctx.Logger)
	if err != nil {
		return Resolved[KID]{}, err
	}
	return Resolved[KID]{
		Width:  size.Width,
		Height: size.Height,
		Root:   root,
	}, nil
}

func resolveWithPath[KID KeelID](renderable Renderable, rect Rect, path string, logger LoggerFunc) (ResolvedNode[KID], error) {
	switch n := renderable.(type) {
	case Container:
		return resolveContainerWithPath[KID](n, rect, path, logger)
	case Block[KID]:
		return ResolvedNode[KID]{
			Kind:  NodeBlock,
			Rect:  rect,
			Block: n,
		}, nil
	default:
		err := &ConfigError{Reason: ErrUnknownRenderable}
		logError(logger, path, "dispatch", err)
		return ResolvedNode[KID]{}, err
	}
}

func resolveContainerWithPath[KID KeelID](container Container, rect Rect, path string, logger LoggerFunc) (ResolvedNode[KID], error) {
	length := container.Len()
	if length <= 0 {
		return ResolvedNode[KID]{
			Kind:     NodeContainer,
			Rect:     rect,
			Children: nil,
		}, nil
	}

	axis := container.GetAxis()
	if axis != AxisHorizontal && axis != AxisVertical {
		err := &ConfigError{Reason: ErrInvalidAxis}
		logError(logger, path, "container.axis", err)
		return ResolvedNode[KID]{}, err
	}

	extents, err := GetContainerExtents(container)
	if err != nil {
		logError(logger, path, "container.slot", err)
		return ResolvedNode[KID]{}, err
	}

	total := rect.Width
	if axis == AxisVertical {
		total = rect.Height
	}

	sizes, required, err := ResolveExtents(total, extents)
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
		logError(logger, path, "container.resolve", err)
		return ResolvedNode[KID]{}, err
	}

	logf(
		logger,
		path,
		LogEventContainerAlloc,
		axis.String(),
		total,
		len(sizes),
		sizes,
		required,
	)

	children := make([]ResolvedNode[KID], length)
	offset := 0
	for i, size := range sizes {
		slot, ok := container.Slot(i)
		if !ok || slot == nil {
			err := &SlotError{Index: i, Reason: ErrNilSlot}
			logError(logger, path, "container.slot", err)
			return ResolvedNode[KID]{}, err
		}

		childRect := rect
		if axis == AxisHorizontal {
			childRect.X += offset
			childRect.Width = size
		} else {
			childRect.Y += offset
			childRect.Height = size
		}

		childPath := path
		if logger != nil {
			childPath = appendPath(path, i)
		}

		child, err := resolveWithPath[KID](slot, childRect, childPath, logger)
		if err != nil {
			logError(logger, path, "container.render", err)
			return ResolvedNode[KID]{}, err
		}

		children[i] = child
		offset += size
	}

	return ResolvedNode[KID]{
		Kind:     NodeContainer,
		Axis:     axis,
		Rect:     rect,
		Children: children,
	}, nil
}
