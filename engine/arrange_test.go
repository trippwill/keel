package engine

import (
	"errors"
	"reflect"
	"testing"

	"github.com/trippwill/keel/core"
)

func TestAllocateValidation(t *testing.T) {
	cases := []struct {
		name  string
		total int
		specs []core.ExtentConstraint
		err   error
	}{
		{
			name:  "negative total",
			total: -1,
			specs: []core.ExtentConstraint{{Kind: core.ExtentFixed, Units: 1}},
			err:   core.ErrConfigurationInvalid,
		},
		{
			name:  "units must be positive",
			total: 1,
			specs: []core.ExtentConstraint{{Kind: core.ExtentFixed, Units: 0}},
			err:   core.ErrConfigurationInvalid,
		},
		{
			name:  "content min non-negative",
			total: 1,
			specs: []core.ExtentConstraint{{Kind: core.ExtentFlex, Units: 1, MinCells: -1}},
			err:   core.ErrConfigurationInvalid,
		},
		{
			name:  "content max non-negative",
			total: 1,
			specs: []core.ExtentConstraint{{Kind: core.ExtentFlex, Units: 1, MaxCells: -1}},
			err:   core.ErrConfigurationInvalid,
		},
		{
			name:  "fixed must cover content min",
			total: 10,
			specs: []core.ExtentConstraint{{Kind: core.ExtentFixed, Units: 1, MinCells: 2}},
			err:   core.ErrConfigurationInvalid,
		},
		{
			name:  "max must cover min",
			total: 10,
			specs: []core.ExtentConstraint{{Kind: core.ExtentFlex, Units: 1, MinCells: 3, MaxCells: 2}},
			err:   core.ErrConfigurationInvalid,
		},
		{
			name:  "invalid size kind",
			total: 10,
			specs: []core.ExtentConstraint{{Kind: core.ExtentKind(99), Units: 1}},
			err:   core.ErrConfigurationInvalid,
		},
		{
			name:  "required exceeds total",
			total: 2,
			specs: []core.ExtentConstraint{{Kind: core.ExtentFlex, Units: 1, MinCells: 3}},
			err:   core.ErrExtentTooSmall,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := ArrangeExtents(tc.total, tc.specs)
			if !errors.Is(err, tc.err) {
				t.Fatalf("expected %v, got %v", tc.err, err)
			}
		})
	}
}

func TestArrangeEmptyExtents(t *testing.T) {
	cases := []struct {
		name  string
		total int
	}{
		{name: "zero total", total: 0},
		{name: "positive total", total: 10},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, required, err := ArrangeExtents(tc.total, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != 0 {
				t.Fatalf("expected empty sizes, got %v", got)
			}
			if required != 0 {
				t.Fatalf("expected required 0, got %d", required)
			}
		})
	}
}

func TestAllocateDistribution(t *testing.T) {
	cases := []struct {
		name  string
		total int
		specs []core.ExtentConstraint
		want  []int
	}{
		{
			name:  "flex distribution with remainder",
			total: 10,
			specs: []core.ExtentConstraint{
				{Kind: core.ExtentFlex, Units: 1, MinCells: 0},
				{Kind: core.ExtentFlex, Units: 3, MinCells: 0},
			},
			want: []int{3, 7},
		},
		{
			name:  "leftover goes to last when no flex",
			total: 5,
			specs: []core.ExtentConstraint{
				{Kind: core.ExtentFixed, Units: 2, MinCells: 0},
				{Kind: core.ExtentFixed, Units: 1, MinCells: 0},
			},
			want: []int{2, 3},
		},
		{
			name:  "mix of fixed and flex with content min",
			total: 10,
			specs: []core.ExtentConstraint{
				{Kind: core.ExtentFixed, Units: 2, MinCells: 1},
				{Kind: core.ExtentFlex, Units: 2, MinCells: 3},
				{Kind: core.ExtentFlex, Units: 1, MinCells: 1},
			},
			want: []int{2, 6, 2},
		},
		{
			name:  "max caps flex distribution",
			total: 10,
			specs: []core.ExtentConstraint{
				{Kind: core.ExtentFlex, Units: 1, MaxCells: 3},
				{Kind: core.ExtentFlex, Units: 1},
			},
			want: []int{3, 7},
		},
		{
			name:  "soft max releases when needed",
			total: 10,
			specs: []core.ExtentConstraint{
				{Kind: core.ExtentFlex, Units: 1, MaxCells: 3},
				{Kind: core.ExtentFlex, Units: 1, MaxCells: 3},
			},
			want: []int{5, 5},
		},
		{
			name:  "fixed ignores max",
			total: 5,
			specs: []core.ExtentConstraint{
				{Kind: core.ExtentFixed, Units: 2, MaxCells: 1},
				{Kind: core.ExtentFlex, Units: 1},
			},
			want: []int{2, 3},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, _, err := ArrangeExtents(tc.total, tc.specs)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
}

func TestArrangeInvariants(t *testing.T) {
	specs := []core.ExtentConstraint{
		{Kind: core.ExtentFixed, Units: 2, MinCells: 2},
		{Kind: core.ExtentFlex, Units: 1, MinCells: 1},
		{Kind: core.ExtentFixed, Units: 3, MinCells: 3},
		{Kind: core.ExtentFlex, Units: 2, MinCells: 0},
	}
	total := 12

	sizes, _, err := ArrangeExtents(total, specs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sum := 0
	for i, size := range sizes {
		sum += size
		if size < specs[i].MinCells {
			t.Fatalf("slot %d expected >= %d, got %d", i, specs[i].MinCells, size)
		}
		if specs[i].Kind == core.ExtentFixed && size != specs[i].Units {
			t.Fatalf("slot %d expected fixed %d, got %d", i, specs[i].Units, size)
		}
	}
	if sum != total {
		t.Fatalf("expected total %d, got %d", total, sum)
	}
}

func TestArrangeSoftMaxExceedsWhenNeeded(t *testing.T) {
	specs := []core.ExtentConstraint{
		{Kind: core.ExtentFlex, Units: 1, MaxCells: 2},
		{Kind: core.ExtentFlex, Units: 1, MaxCells: 2},
	}

	sizes, _, err := ArrangeExtents(7, specs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sizes[0] <= specs[0].MaxCells || sizes[1] <= specs[1].MaxCells {
		t.Fatalf("expected soft max to be exceeded, got %v", sizes)
	}
	if sizes[0]+sizes[1] != 7 {
		t.Fatalf("expected total 7, got %d", sizes[0]+sizes[1])
	}
}

func TestArrangeMaxDistributesRemainder(t *testing.T) {
	specs := []core.ExtentConstraint{
		{Kind: core.ExtentFlex, Units: 1, MaxCells: 4},
		{Kind: core.ExtentFlex, Units: 1, MaxCells: 4},
	}

	sizes, _, err := ArrangeExtents(5, specs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []int{3, 2}
	if !reflect.DeepEqual(sizes, want) {
		t.Fatalf("expected %v, got %v", want, sizes)
	}
}

func TestArrangeFlexMinMaxHonorsBounds(t *testing.T) {
	specs := []core.ExtentConstraint{
		{Kind: core.ExtentFlex, Units: 1, MinCells: 2, MaxCells: 4},
		{Kind: core.ExtentFlex, Units: 1},
	}

	sizes, _, err := ArrangeExtents(6, specs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sizes[0] != 4 {
		t.Fatalf("expected max 4 for first slot, got %d", sizes[0])
	}
	if sizes[1] != 2 {
		t.Fatalf("expected remaining 2 for second slot, got %d", sizes[1])
	}
}
