//go:generate stringer -type=Axis
//go:generate stringer -type=SizeKind
package keel

// Axis represents a layout axis: horizontal or vertical.
type Axis uint8

const (
	AxisHorizontal Axis = 0
	AxisVertical   Axis = 1
)

// SizeKind represents whether a SizeSpec is fixed or flexible.
type SizeKind uint8

const (
	SizeFixed SizeKind = iota
	SizeFlex
)

// SizeSpec defines how much space a panel should take along an axis.
type SizeSpec struct {
	Kind       SizeKind
	Units      int
	ContentMin int // Minimum cells to reserve for content on this axis
}

func Fixed(units int) SizeSpec {
	return SizeSpec{SizeFixed, units, 0}
}

func FixedMin(units int, minReserved int) SizeSpec {
	return SizeSpec{SizeFixed, units, minReserved}
}

func Flex(units int) SizeSpec {
	return SizeSpec{SizeFlex, units, 0}
}

func FlexMin(units int, minReserved int) SizeSpec {
	return SizeSpec{SizeFlex, units, minReserved}
}
