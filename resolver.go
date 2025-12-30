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
	required, flexUnits, hasFlex, hasFlexMax, err := seedSizes(sizes, extents)
	if err != nil {
		return nil, required, err
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
			flexSpecs := collectFlexSpecs(extents)
			remaining := distributeFlexWithMax(sizes, flexSpecs, leftover)
			if remaining > 0 {
				distributeFlexIgnoringMax(sizes, extents, flexUnits, remaining)
			}
		} else {
			distributeFlexIgnoringMax(sizes, extents, flexUnits, leftover)
		}
	}

	return sizes, required, nil
}

type flexSpec struct {
	index int
	units int
	max   int
}

func seedSizes(sizes []int, extents []ExtentConstraint) (int, int, bool, bool, error) {
	required := 0
	flexUnits := 0
	hasFlex := false
	hasFlexMax := false

	for i := range extents {
		spec := extents[i]
		if spec.Units <= 0 {
			return required, flexUnits, hasFlex, hasFlexMax, &ExtentError{Index: i, Reason: ErrInvalidExtentUnits}
		}
		if spec.MinCells < 0 {
			return required, flexUnits, hasFlex, hasFlexMax, &ExtentError{Index: i, Reason: ErrInvalidExtentMinCells}
		}
		if spec.MaxCells < 0 {
			return required, flexUnits, hasFlex, hasFlexMax, &ExtentError{Index: i, Reason: ErrInvalidExtentMaxCells}
		}

		switch spec.Kind {
		case ExtentFixed:
			if spec.Units < spec.MinCells {
				return required, flexUnits, hasFlex, hasFlexMax, &ExtentError{Index: i, Reason: ErrInvalidExtentMin}
			}
			sizes[i] = spec.Units
		case ExtentFlex:
			if spec.MaxCells > 0 && spec.MaxCells < spec.MinCells {
				return required, flexUnits, hasFlex, hasFlexMax, &ExtentError{Index: i, Reason: ErrInvalidExtentMax}
			}
			sizes[i] = spec.MinCells
			flexUnits += spec.Units
			hasFlex = true
			if spec.MaxCells > 0 {
				hasFlexMax = true
			}
		default:
			return required, flexUnits, hasFlex, hasFlexMax, &ExtentError{Index: i, Reason: ErrInvalidExtentKind}
		}

		required += sizes[i]
	}

	return required, flexUnits, hasFlex, hasFlexMax, nil
}

func collectFlexSpecs(extents []ExtentConstraint) []flexSpec {
	flexSpecs := make([]flexSpec, 0, len(extents))
	for i, spec := range extents {
		if spec.Kind != ExtentFlex {
			continue
		}
		flexSpecs = append(flexSpecs, flexSpec{
			index: i,
			units: spec.Units,
			max:   spec.MaxCells,
		})
	}
	return flexSpecs
}

func distributeFlexWithMax(sizes []int, flexSpecs []flexSpec, leftover int) int {
	if leftover <= 0 {
		return 0
	}

	unlimited := false
	capacity := 0
	for _, spec := range flexSpecs {
		if spec.max == 0 {
			unlimited = true
			continue
		}
		if sizes[spec.index] < spec.max {
			capacity += spec.max - sizes[spec.index]
		}
	}

	if len(flexSpecs) == 0 {
		return leftover
	}

	amount := leftover
	remaining := 0
	if !unlimited && capacity < leftover {
		amount = capacity
		remaining = leftover - capacity
	}

	if amount > 0 {
		remaining += distributeFlexCapped(sizes, flexSpecs, amount)
	}

	return remaining
}

func distributeFlexCapped(sizes []int, flexSpecs []flexSpec, amount int) int {
	remaining := amount
	active := make([]int, 0, len(flexSpecs))
	for i, spec := range flexSpecs {
		if spec.max == 0 || sizes[spec.index] < spec.max {
			active = append(active, i)
		}
	}

	for remaining > 0 && len(active) > 0 {
		totalUnits := 0
		for _, specIndex := range active {
			totalUnits += flexSpecs[specIndex].units
		}

		distributed := 0
		for _, specIndex := range active {
			spec := flexSpecs[specIndex]
			add := 0
			if totalUnits > 0 {
				add = remaining * spec.units / totalUnits
			}
			cap := maxFlexAdd(spec.max, sizes[spec.index], remaining)
			if add > cap {
				add = cap
			}
			if add > 0 {
				sizes[spec.index] += add
				distributed += add
			}
		}
		remaining -= distributed

		if remaining > 0 {
			for _, specIndex := range active {
				if remaining == 0 {
					break
				}
				spec := flexSpecs[specIndex]
				cap := maxFlexAdd(spec.max, sizes[spec.index], remaining)
				if cap <= 0 {
					continue
				}
				sizes[spec.index]++
				remaining--
			}
		}

		next := active[:0]
		for _, specIndex := range active {
			spec := flexSpecs[specIndex]
			if spec.max == 0 || sizes[spec.index] < spec.max {
				next = append(next, specIndex)
			}
		}
		active = next
	}

	return remaining
}

func maxFlexAdd(max int, size int, remaining int) int {
	if max == 0 {
		return remaining
	}
	if size >= max {
		return 0
	}
	return max - size
}

func distributeFlexIgnoringMax(sizes []int, extents []ExtentConstraint, flexUnits int, leftover int) {
	if leftover <= 0 {
		return
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
