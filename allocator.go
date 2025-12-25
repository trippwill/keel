package keel

// Allocator distributes a total number of cells across slot SizeSpecs.
type Allocator func(total int, specs []SizeSpec) ([]int, error)

// RowAllocator distributes width across slots for horizontal splits.
func RowAllocator(total int, specs []SizeSpec) ([]int, error) {
	return Allocate(total, specs)
}

// ColAllocator distributes height across slots for vertical splits.
func ColAllocator(total int, specs []SizeSpec) ([]int, error) {
	return Allocate(total, specs)
}

// Allocate distributes a total number of cells across SizeSpecs.
func Allocate(total int, specs []SizeSpec) ([]int, error) {
	if total < 0 {
		return nil, ErrConfigurationInvalid
	}
	if len(specs) == 0 {
		return nil, ErrConfigurationInvalid
	}

	sizes := make([]int, len(specs))
	flexIdx := make([]int, 0, len(specs))
	flexUnits := 0
	required := 0

	for i, spec := range specs {
		if spec.Units <= 0 {
			return nil, ErrConfigurationInvalid
		}
		if spec.ContentMin < 0 {
			return nil, ErrConfigurationInvalid
		}

		switch spec.Kind {
		case SizeFixed:
			if spec.Units < spec.ContentMin {
				return nil, ErrConfigurationInvalid
			}
			sizes[i] = spec.Units
		case SizeFlex:
			sizes[i] = spec.ContentMin
			flexIdx = append(flexIdx, i)
			flexUnits += spec.Units
		default:
			return nil, ErrConfigurationInvalid
		}

		required += sizes[i]
	}

	if required > total {
		return nil, ErrTargetTooSmall
	}

	leftover := total - required
	if len(flexIdx) == 0 {
		if leftover > 0 {
			sizes[len(sizes)-1] += leftover
		}
		return sizes, nil
	}

	remainder := leftover
	for _, i := range flexIdx {
		add := 0
		if flexUnits > 0 {
			add = leftover * specs[i].Units / flexUnits
		}
		sizes[i] += add
		remainder -= add
	}

	for i := 0; i < remainder; i++ {
		idx := flexIdx[i%len(flexIdx)]
		sizes[idx]++
	}

	return sizes, nil
}
