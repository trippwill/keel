package keel

import "testing"

func TestExtentConstraintHelpers(t *testing.T) {
	cases := []struct {
		name string
		got  ExtentConstraint
		want ExtentConstraint
	}{
		{
			name: "flex unit",
			got:  FlexUnit(),
			want: ExtentConstraint{Kind: ExtentFlex, Units: 1, MinCells: 0, MaxCells: 0},
		},
		{
			name: "fixed",
			got:  Fixed(3),
			want: ExtentConstraint{Kind: ExtentFixed, Units: 3, MinCells: 3, MaxCells: 0},
		},
		{
			name: "flex",
			got:  Flex(2),
			want: ExtentConstraint{Kind: ExtentFlex, Units: 2, MinCells: 0, MaxCells: 0},
		},
		{
			name: "flex min",
			got:  FlexMin(2, 5),
			want: ExtentConstraint{Kind: ExtentFlex, Units: 2, MinCells: 5, MaxCells: 0},
		},
		{
			name: "flex max",
			got:  FlexMax(2, 5),
			want: ExtentConstraint{Kind: ExtentFlex, Units: 2, MinCells: 0, MaxCells: 5},
		},
		{
			name: "flex min max",
			got:  FlexMinMax(2, 3, 5),
			want: ExtentConstraint{Kind: ExtentFlex, Units: 2, MinCells: 3, MaxCells: 5},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Fatalf("expected %v, got %v", tc.want, tc.got)
			}
		})
	}
}

func TestExtentConstraintExtent(t *testing.T) {
	spec := FlexMinMax(2, 1, 4)
	if got := spec.Extent(); got != spec {
		t.Fatalf("expected %v, got %v", spec, got)
	}
}
