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
		specs []SizeSpec
		err   error
	}{
		{
			name:  "negative total",
			total: -1,
			specs: []SizeSpec{Flex(1)},
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
			specs: []SizeSpec{{Kind: SizeFixed, Units: 0}},
			err:   ErrConfigurationInvalid,
		},
		{
			name:  "content min non-negative",
			total: 1,
			specs: []SizeSpec{{Kind: SizeFlex, Units: 1, ContentMin: -1}},
			err:   ErrConfigurationInvalid,
		},
		{
			name:  "fixed must cover content min",
			total: 10,
			specs: []SizeSpec{{Kind: SizeFixed, Units: 1, ContentMin: 2}},
			err:   ErrConfigurationInvalid,
		},
		{
			name:  "invalid size kind",
			total: 10,
			specs: []SizeSpec{{Kind: SizeKind(99), Units: 1}},
			err:   ErrConfigurationInvalid,
		},
		{
			name:  "required exceeds total",
			total: 2,
			specs: []SizeSpec{{Kind: SizeFlex, Units: 1, ContentMin: 3}},
			err:   ErrTargetTooSmall,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Allocate(tc.total, tc.specs)
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
		specs []SizeSpec
		want  []int
	}{
		{
			name:  "flex distribution with remainder",
			total: 10,
			specs: []SizeSpec{
				{Kind: SizeFlex, Units: 1, ContentMin: 0},
				{Kind: SizeFlex, Units: 3, ContentMin: 0},
			},
			want: []int{3, 7},
		},
		{
			name:  "leftover goes to last when no flex",
			total: 5,
			specs: []SizeSpec{
				{Kind: SizeFixed, Units: 2, ContentMin: 0},
				{Kind: SizeFixed, Units: 1, ContentMin: 0},
			},
			want: []int{2, 3},
		},
		{
			name:  "mix of fixed and flex with content min",
			total: 10,
			specs: []SizeSpec{
				{Kind: SizeFixed, Units: 2, ContentMin: 1},
				{Kind: SizeFlex, Units: 2, ContentMin: 3},
				{Kind: SizeFlex, Units: 1, ContentMin: 1},
			},
			want: []int{2, 6, 2},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Allocate(tc.total, tc.specs)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
}
