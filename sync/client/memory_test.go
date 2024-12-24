package client_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/sync/client"
)

func TestMemory(t *testing.T) {
	t.Parallel()

	mem := client.NewMemory()

	now := time.Now()
	items := []item.Item{
		{ID: "a", Kind: item.KindSchedule, Updated: now.Add(-15 * time.Minute)},
		{ID: "b", Kind: item.KindTask, Updated: now.Add(-10 * time.Minute)},
		{ID: "c", Kind: item.KindSchedule, Updated: now.Add(-5 * time.Minute)},
	}
	if err := mem.Update(items); err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	for _, tc := range []struct {
		name     string
		ks       []item.Kind
		ts       time.Time
		expItems []item.Item
	}{
		{
			name:     "empty",
			ks:       make([]item.Kind, 0),
			expItems: make([]item.Item, 0),
		},
		{
			name:     "kind",
			ks:       []item.Kind{item.KindTask},
			expItems: []item.Item{items[1]},
		},
		{
			name:     "timestamp",
			ks:       []item.Kind{item.KindSchedule, item.KindTask},
			ts:       now.Add(-10 * time.Minute),
			expItems: items[1:],
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actItems, actErr := mem.Updated(tc.ks, tc.ts)
			if actErr != nil {
				t.Errorf("exp nil, got %v", actErr)
			}
			if diff := cmp.Diff(tc.expItems, actItems); diff != "" {
				t.Errorf("(exp +, got -)\n%s", diff)
			}
		})
	}
}
