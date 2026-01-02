//go:generate stringer -type=Axis -trimprefix=Axis
package core

// Axis represents a layout axis used by stacks to split space.
type Axis uint8

const (
	// AxisHorizontal lays out content left-to-right.
	AxisHorizontal Axis = 0
	// AxisVertical lays out content top-to-bottom.
	AxisVertical Axis = 1
)
