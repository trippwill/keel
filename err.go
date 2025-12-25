package keel

import (
	"errors"
	"fmt"
)

type TargetTooSmallError struct {
	Axis       Axis
	Need, Have int
}

func (e *TargetTooSmallError) Error() string {
	return fmt.Sprintf(
		"target too small on %s axis: need %d, have %d",
		e.Axis,
		e.Need,
		e.Have,
	)
}

var ErrConfigurationInvalid = fmt.Errorf("configuration invalid")

var ErrAllocatorUnimplemented = errors.New("allocator unimplemented")

var ErrTargetTooSmall = errors.New("target too small")
