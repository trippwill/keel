package keel

import "github.com/trippwill/keel/engine"

// FlexUnit returns a single flexible unit of space with no minimum.
func FlexUnit() engine.ExtentConstraint {
	return engine.ExtentConstraint{Kind: engine.ExtentFlex, Units: 1, MinCells: 0, MaxCells: 0}
}

// Fixed creates a fixed [engine.ExtentConstraint] with the given units in cells.
// Fixed extents must be at least their MinCells value.
func Fixed(units int) engine.ExtentConstraint {
	return engine.ExtentConstraint{Kind: engine.ExtentFixed, Units: units, MinCells: units, MaxCells: 0}
}

// Flex creates a flexible [engine.ExtentConstraint] with the given units in flex space.
// Flex extents receive space after fixed and minimum allocations are satisfied.
func Flex(units int) engine.ExtentConstraint {
	return engine.ExtentConstraint{Kind: engine.ExtentFlex, Units: units, MinCells: 0, MaxCells: 0}
}

// FlexMin creates a flexible [engine.ExtentConstraint] with the given units in flex space,
// and reserves at least minReserved total cells along the stack axis.
func FlexMin(units int, minReserved int) engine.ExtentConstraint {
	return engine.ExtentConstraint{Kind: engine.ExtentFlex, Units: units, MinCells: minReserved, MaxCells: 0}
}

// FlexMax creates a flexible [engine.ExtentConstraint] with the given units in flex space,
// and caps at maxCells total cells along the stack axis.
func FlexMax(units int, maxCells int) engine.ExtentConstraint {
	return engine.ExtentConstraint{Kind: engine.ExtentFlex, Units: units, MinCells: 0, MaxCells: maxCells}
}

// FlexMinMax creates a flexible [engine.ExtentConstraint] with the given units in flex space,
// reserving at least minReserved and capping at maxCells total cells along the axis.
func FlexMinMax(units int, minReserved int, maxCells int) engine.ExtentConstraint {
	return engine.ExtentConstraint{Kind: engine.ExtentFlex, Units: units, MinCells: minReserved, MaxCells: maxCells}
}
