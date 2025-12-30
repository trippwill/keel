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
		return []int{}, 0, nil
	}

	sizes := make([]int, count)
	required := 0
	hasFlex := false
	hasFlexMax := false

	// Pass 1: validate and allocate fixed extents, accumulate flex units and min sizes.
	for i := range count {
		spec := extents[i]
		if spec.Units <= 0 {
			return nil, required, &ExtentError{Index: i, Reason: ErrInvalidExtentUnits}
		}
		if spec.MinCells < 0 {
			return nil, required, &ExtentError{Index: i, Reason: ErrInvalidExtentMinCells}
		}
		if spec.MaxCells < 0 {
			return nil, required, &ExtentError{Index: i, Reason: ErrInvalidExtentMaxCells}
		}

		switch spec.Kind {
		case ExtentFixed:
			if spec.Units < spec.MinCells {
				return nil, required, &ExtentError{Index: i, Reason: ErrInvalidExtentMin}
			}
			sizes[i] = spec.Units
		case ExtentFlex:
			if spec.MaxCells > 0 && spec.MaxCells < spec.MinCells {
				return nil, required, &ExtentError{Index: i, Reason: ErrInvalidExtentMax}
			}
			sizes[i] = spec.MinCells
			hasFlex = true
			if spec.MaxCells > 0 {
				hasFlexMax = true
			}
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

	if leftover > 0 {
		if hasFlexMax {
			remaining := distributeFlexWithMax(sizes, extents, leftover)
			if remaining > 0 {
				distributeFlexIgnoringMax(sizes, extents, remaining)
			}
		} else {
			distributeFlexIgnoringMax(sizes, extents, leftover)
		}
	}

	return sizes, required, nil
}

func distributeFlexWithMax(sizes []int, extents []ExtentConstraint, leftover int) int {
	if leftover <= 0 {
		return 0
	}

	flexIndices := make([]int, 0, len(extents))
	unlimited := false
	capacity := 0
	for i, spec := range extents {
		if spec.Kind != ExtentFlex {
			continue
		}
		flexIndices = append(flexIndices, i)
		if spec.MaxCells == 0 {
			unlimited = true
			continue
		}
		if sizes[i] < spec.MaxCells {
			capacity += spec.MaxCells - sizes[i]
		}
	}

	if len(flexIndices) == 0 {
		return leftover
	}

	amount := leftover
	remaining := 0
	if !unlimited && capacity < leftover {
		amount = capacity
		remaining = leftover - capacity
	}

	if amount > 0 {
		remaining += distributeFlexCapped(sizes, extents, flexIndices, amount)
	}

	return remaining
}

func distributeFlexCapped(sizes []int, extents []ExtentConstraint, flexIndices []int, amount int) int {
	remaining := amount
	active := make([]int, 0, len(flexIndices))
	for _, i := range flexIndices {
		spec := extents[i]
		if spec.MaxCells == 0 || sizes[i] < spec.MaxCells {
			active = append(active, i)
		}
	}

	for remaining > 0 && len(active) > 0 {
		totalUnits := 0
		for _, i := range active {
			totalUnits += extents[i].Units
		}

		distributed := 0
		for _, i := range active {
			add := 0
			if totalUnits > 0 {
				add = remaining * extents[i].Units / totalUnits
			}
			cap := maxFlexAdd(extents[i], sizes[i], remaining)
			if add > cap {
				add = cap
			}
			if add > 0 {
				sizes[i] += add
				distributed += add
			}
		}
		remaining -= distributed

		if remaining > 0 {
			for _, i := range active {
				if remaining == 0 {
					break
				}
				cap := maxFlexAdd(extents[i], sizes[i], remaining)
				if cap <= 0 {
					continue
				}
				sizes[i]++
				remaining--
			}
		}

		next := active[:0]
		for _, i := range active {
			spec := extents[i]
			if spec.MaxCells == 0 || sizes[i] < spec.MaxCells {
				next = append(next, i)
			}
		}
		active = next
	}

	return remaining
}

func maxFlexAdd(spec ExtentConstraint, size int, remaining int) int {
	if spec.MaxCells == 0 {
		return remaining
	}
	if size >= spec.MaxCells {
		return 0
	}
	return spec.MaxCells - size
}

func distributeFlexIgnoringMax(sizes []int, extents []ExtentConstraint, leftover int) {
	if leftover <= 0 {
		return
	}

	flexUnits := 0
	for _, spec := range extents {
		if spec.Kind == ExtentFlex {
			flexUnits += spec.Units
		}
	}
	if flexUnits == 0 {
		return
	}

	remainder := leftover
	for i, spec := range extents {
		if spec.Kind != ExtentFlex {
			continue
		}
		add := leftover * spec.Units / flexUnits
		sizes[i] += add
		remainder -= add
	}

	if remainder > 0 {
		for i := 0; i < len(extents) && remainder > 0; i++ {
			spec := extents[i]
			if spec.Kind != ExtentFlex {
				continue
			}
			sizes[i]++
			remainder--
		}
	}
}
