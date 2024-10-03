package storage_test

import (
	"testing"

	"go-mod.ewintr.nl/planner/plan/storage"
)

func TestNextLocalId(t *testing.T) {
	for _, tc := range []struct {
		name string
		used []int
		exp  int
	}{
		{
			name: "empty",
			used: []int{},
			exp:  1,
		},
		{
			name: "not empty",
			used: []int{5},
			exp:  6,
		},
		{
			name: "multiple",
			used: []int{2, 3, 4},
			exp:  5,
		},
		{
			name: "holes",
			used: []int{1, 5, 8},
			exp:  9,
		},
		{
			name: "expand limit",
			used: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			exp:  11,
		},
		{
			name: "wrap if possible",
			used: []int{8, 9},
			exp:  1,
		},
		{
			name: "find hole",
			used: []int{1, 2, 3, 4, 5, 7, 8, 9},
			exp:  6,
		},
		{
			name: "dont wrap if expanded before",
			used: []int{15, 16},
			exp:  17,
		},
		{
			name: "do wrap if expanded limit is reached",
			used: []int{99},
			exp:  1,
		},
		{
			name: "sync bug",
			used: []int{151, 956, 955, 150, 154, 155, 145, 144,
				136, 152, 148, 146, 934, 149, 937, 135, 140, 139,
				143, 137, 153, 939, 138, 953, 147, 141, 938, 142,
			},
			exp: 957,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			act := storage.NextLocalID(tc.used)
			if tc.exp != act {
				t.Errorf("exp %v, got %v", tc.exp, act)
			}
		})
	}
}
