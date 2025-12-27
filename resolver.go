package keel

// RowResolver distributes width across slots for horizontal splits.
//
// Arguments:
//
//	total:    The total width to distribute.
//	container: The Container whose slots will be allocated.
//
// Returns:
//   - Per-slot widths ([]int)
//   - Minimum required total (int)
//   - Error, if allocation fails
//
// This function mirrors the allocation rules used by containers. The slot extents are determined by calling Slot(i).GetExtent() on the container.
func RowResolver(total int, container Container) ([]int, int, error) {
	extents, err := GetContainerExtents(container)
	if err != nil {
		return nil, 0, err
	}
	return ResolveExtents(total, extents)
}

// ColResolver distributes height across slots for vertical splits.
//
// Arguments:
//
//	total:    The total height to distribute.
//	container: The Container whose slots will be allocated.
//
// Returns:
//   - Per-slot heights ([]int)
//   - Minimum required total (int)
//   - Error, if allocation fails
//
// This function mirrors the allocation rules used by containers. The slot extents are determined by calling Slot(i).GetExtent() on the container.
func ColResolver(total int, container Container) ([]int, int, error) {
	extents, err := GetContainerExtents(container)
	if err != nil {
		return nil, 0, err
	}
	return ResolveExtents(total, extents)
}

// ResolveExtents distributes a total number of cells across slot extents.
//
// Arguments:
//
//	total:    The total number of cells to distribute.
//	extents: The ExtentConstraints for each slot.
//
// Returns:
//   - Per-slot sizes ([]int)
//   - Minimum required total (int)
//   - Error, if allocation fails
func ResolveExtents(total int, extents []ExtentConstraint) ([]int, int, error) {
	if total < 0 {
		return nil, 0, &ConfigError{Reason: ErrInvalidTotal}
	}

	count := len(extents)
	if count <= 0 {
		return nil, 0, &ConfigError{Reason: ErrEmptyExtents}
	}

	sizes := make([]int, count)
	flexUnits := 0
	required := 0
	hasFlex := false

	// Pass 1: validate and allocate fixed extents, accumulate flex units and min sizes.
	for i := range count {
		spec := extents[i]
		if spec.Units <= 0 {
			return nil, required, &ExtentError{Index: i, Reason: ErrInvalidExtentUnits}
		}
		if spec.MinCells < 0 {
			return nil, required, &ExtentError{Index: i, Reason: ErrInvalidExtentMinCells}
		}

		switch spec.Kind {
		case ExtentFixed:
			if spec.Units < spec.MinCells {
				return nil, required, &ExtentError{Index: i, Reason: ErrInvalidExtentMin}
			}
			sizes[i] = spec.Units
		case ExtentFlex:
			sizes[i] = spec.MinCells
			flexUnits += spec.Units
			hasFlex = true
		default:
			return nil, required, &ExtentError{Index: i, Reason: ErrInvalidExtentKind}
		}

		required += sizes[i]
	}

	if required > total {
		return nil, required, ErrExtentTooSmall
	}

	// Pass 2: distribute leftover space to flex extents.
	leftover := total - required
	if !hasFlex {
		if leftover > 0 {
			sizes[len(sizes)-1] += leftover
		}
		return sizes, required, nil
	}

	remainder := leftover
	for i := range count {
		spec := extents[i]
		if spec.Kind != ExtentFlex {
			continue
		}
		add := 0
		if flexUnits > 0 {
			add = leftover * spec.Units / flexUnits
		}
		sizes[i] += add
		remainder -= add
	}

	if remainder > 0 {
		for i := 0; i < count && remainder > 0; i++ {
			spec := extents[i]
			if spec.Kind != ExtentFlex {
				continue
			}
			sizes[i]++
			remainder--
		}
	}

	return sizes, required, nil
}

// GetContainerExtents retrieves the extents of all slots in a container.
//
// Returns:
// - Slice of ExtentConstraint for each slot
// - Error, if any slot is nil
func GetContainerExtents(container Container) ([]ExtentConstraint, error) {
	extents := make([]ExtentConstraint, container.Len())
	for i := range extents {
		slot, ok := container.Slot(i)
		if !ok || slot == nil {
			return nil, &SlotError{Index: i, Reason: ErrNilSlot}
		}

		extents[i] = slot.GetExtent()
	}

	return extents, nil
}
