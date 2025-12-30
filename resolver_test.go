package keel

import (
	"errors"
	"reflect"
	"testing"
)

func TestAllocateValidation(t *testing.T) {
	cases := []struct {
		name  string
		total int
		specs []ExtentConstraint
		err   error
	}{
		{
			name:  "negative total",
			total: -1,
			specs: []ExtentConstraint{Flex(1)},
			err:   ErrConfigurationInvalid,
		},
		{
			name:  "units must be positive",
			total: 1,
			specs: []ExtentConstraint{{Kind: ExtentFixed, Units: 0}},
			err:   ErrConfigurationInvalid,
		},
		{
			name:  "content min non-negative",
			total: 1,
			specs: []ExtentConstraint{{Kind: ExtentFlex, Units: 1, MinCells: -1}},
			err:   ErrConfigurationInvalid,
		},
		{
			name:  "content max non-negative",
			total: 1,
			specs: []ExtentConstraint{{Kind: ExtentFlex, Units: 1, MaxCells: -1}},
			err:   ErrConfigurationInvalid,
		},
		{
			name:  "fixed must cover content min",
			total: 10,
			specs: []ExtentConstraint{{Kind: ExtentFixed, Units: 1, MinCells: 2}},
			err:   ErrConfigurationInvalid,
		},
		{
			name:  "max must cover min",
			total: 10,
			specs: []ExtentConstraint{{Kind: ExtentFlex, Units: 1, MinCells: 3, MaxCells: 2}},
			err:   ErrConfigurationInvalid,
		},
		{
			name:  "invalid size kind",
			total: 10,
			specs: []ExtentConstraint{{Kind: ExtentKind(99), Units: 1}},
			err:   ErrConfigurationInvalid,
		},
		{
			name:  "required exceeds total",
			total: 2,
			specs: []ExtentConstraint{{Kind: ExtentFlex, Units: 1, MinCells: 3}},
			err:   ErrExtentTooSmall,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := ResolveExtents(tc.total, tc.specs)
			if !errors.Is(err, tc.err) {
				t.Fatalf("expected %v, got %v", tc.err, err)
			}
		})
	}
}

func TestResolveEmptyExtents(t *testing.T) {
	cases := []struct {
		name  string
		total int
	}{
		{name: "zero total", total: 0},
		{name: "positive total", total: 10},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, required, err := ResolveExtents(tc.total, nil)
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
		specs []ExtentConstraint
		want  []int
	}{
		{
			name:  "flex distribution with remainder",
			total: 10,
			specs: []ExtentConstraint{
				{Kind: ExtentFlex, Units: 1, MinCells: 0},
				{Kind: ExtentFlex, Units: 3, MinCells: 0},
			},
			want: []int{3, 7},
		},
		{
			name:  "leftover goes to last when no flex",
			total: 5,
			specs: []ExtentConstraint{
				{Kind: ExtentFixed, Units: 2, MinCells: 0},
				{Kind: ExtentFixed, Units: 1, MinCells: 0},
			},
			want: []int{2, 3},
		},
		{
			name:  "mix of fixed and flex with content min",
			total: 10,
			specs: []ExtentConstraint{
				{Kind: ExtentFixed, Units: 2, MinCells: 1},
				{Kind: ExtentFlex, Units: 2, MinCells: 3},
				{Kind: ExtentFlex, Units: 1, MinCells: 1},
			},
			want: []int{2, 6, 2},
		},
		{
			name:  "max caps flex distribution",
			total: 10,
			specs: []ExtentConstraint{
				{Kind: ExtentFlex, Units: 1, MaxCells: 3},
				{Kind: ExtentFlex, Units: 1},
			},
			want: []int{3, 7},
		},
		{
			name:  "soft max releases when needed",
			total: 10,
			specs: []ExtentConstraint{
				{Kind: ExtentFlex, Units: 1, MaxCells: 3},
				{Kind: ExtentFlex, Units: 1, MaxCells: 3},
			},
			want: []int{5, 5},
		},
		{
			name:  "fixed ignores max",
			total: 5,
			specs: []ExtentConstraint{
				{Kind: ExtentFixed, Units: 2, MaxCells: 1},
				{Kind: ExtentFlex, Units: 1},
			},
			want: []int{2, 3},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, _, err := ResolveExtents(tc.total, tc.specs)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
}

func TestResolveInvariants(t *testing.T) {
	specs := []ExtentConstraint{
		{Kind: ExtentFixed, Units: 2, MinCells: 2},
		{Kind: ExtentFlex, Units: 1, MinCells: 1},
		{Kind: ExtentFixed, Units: 3, MinCells: 3},
		{Kind: ExtentFlex, Units: 2, MinCells: 0},
	}
	total := 12

	sizes, _, err := ResolveExtents(total, specs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sum := 0
	for i, size := range sizes {
		sum += size
		if size < specs[i].MinCells {
			t.Fatalf("slot %d expected >= %d, got %d", i, specs[i].MinCells, size)
		}
		if specs[i].Kind == ExtentFixed && size != specs[i].Units {
			t.Fatalf("slot %d expected fixed %d, got %d", i, specs[i].Units, size)
		}
	}
	if sum != total {
		t.Fatalf("expected total %d, got %d", total, sum)
	}
}

func TestResolveSoftMaxExceedsWhenNeeded(t *testing.T) {
	specs := []ExtentConstraint{
		{Kind: ExtentFlex, Units: 1, MaxCells: 2},
		{Kind: ExtentFlex, Units: 1, MaxCells: 2},
	}

	sizes, _, err := ResolveExtents(7, specs)
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
