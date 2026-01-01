package keel

import (
	"github.com/trippwill/keel/core"
	"github.com/trippwill/keel/engine"
)

// Row creates a new horizontal split.
// Slots are stored as references; mutating slots after creation affects the stack.
func Row(size ExtentConstraint, slots ...Spec) StackSpec {
	return engine.NewSplitSpec(core.AxisHorizontal, size, slots...)
}

// Col creates a new vertical split.
// Slots are stored as references; mutating slots after creation affects the stack.
func Col(size ExtentConstraint, slots ...Spec) StackSpec {
	return engine.NewSplitSpec(core.AxisVertical, size, slots...)
}
