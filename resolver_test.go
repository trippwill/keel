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
			name:  "empty specs",
			total: 1,
			specs: nil,
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
			name:  "fixed must cover content min",
			total: 10,
			specs: []ExtentConstraint{{Kind: ExtentFixed, Units: 1, MinCells: 2}},
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
