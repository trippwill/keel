package keel

// ExtentResolver distributes a total number of cells across slot extents.
// The extentAt callback is invoked for each index as needed.
type ExtentResolver func(total int, count int, extentAt func(index int) (ExtentConstraint, error)) ([]int, int, error)

// RowResolver distributes width across slots for horizontal splits.
func RowResolver(total int, count int, extentAt func(index int) (ExtentConstraint, error)) ([]int, int, error) {
	return resolve(total, count, extentAt)
}

// ColResolver distributes height across slots for vertical splits.
func ColResolver(total int, count int, extentAt func(index int) (ExtentConstraint, error)) ([]int, int, error) {
	return resolve(total, count, extentAt)
}

// resolve distributes a total number of cells across slot extents.
func resolve(total int, count int, extentAt func(index int) (ExtentConstraint, error)) ([]int, int, error) {
	if total < 0 {
		return nil, 0, &ConfigError{Reason: ErrInvalidTotal}
	}
	if count == 0 {
		return nil, 0, &ConfigError{Reason: ErrEmptyExtents}
	}

	sizes := make([]int, count)
	flexUnits := 0
	required := 0
	hasFlex := false

	for i := range count {
		spec, err := extentAt(i)
		if err != nil {
			return nil, required, err
		}
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

	leftover := total - required
	if !hasFlex {
		if leftover > 0 {
			sizes[len(sizes)-1] += leftover
		}
		return sizes, required, nil
	}

	remainder := leftover
	for i := range count {
		spec, err := extentAt(i)
		if err != nil {
			return nil, required, err
		}
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
			spec, err := extentAt(i)
			if err != nil {
				return nil, required, err
			}
			if spec.Kind != ExtentFlex {
				continue
			}
			sizes[i]++
			remainder--
		}
	}

	return sizes, required, nil
}
