package main

import (
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go-mod.ewintr.nl/planner/item"
)

func TestMemoryUpdate(t *testing.T) {
	t.Parallel()

	mem := NewMemory()

	t.Log("start empty")
	actItems, actErr := mem.Updated([]item.Kind{}, time.Time{})
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if len(actItems) != 0 {
		t.Errorf("exp 0, got %d", len(actItems))
	}

	t.Log("add one")
	t1 := item.NewItem(item.Kind("kinda"), "test")
	if actErr := mem.Update(t1, t1.Updated); actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	actItems, actErr = mem.Updated([]item.Kind{}, time.Time{})
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if len(actItems) != 1 {
		t.Errorf("exp 1, gor %d", len(actItems))
	}
	if actItems[0].ID != t1.ID {
		t.Errorf("exp %v, got %v", actItems[0].ID, t1.ID)
	}

	before := time.Now()

	t.Log("add second")
	t2 := item.NewItem(item.Kind("kindb"), "test 2")
	if actErr := mem.Update(t2, t2.Updated); actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	actItems, actErr = mem.Updated([]item.Kind{}, time.Time{})
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if diff := cmp.Diff([]item.Item{t1, t2}, actItems, cmpopts.SortSlices(func(i, j item.Item) bool {
		return i.ID < j.ID
	})); diff != "" {
		t.Errorf("(exp +, got -)\n%s", diff)
	}

	actItems, actErr = mem.Updated([]item.Kind{}, before)
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if len(actItems) != 1 {
		t.Errorf("exp 1, gor %d", len(actItems))
	}
	if actItems[0].ID != t2.ID {
		t.Errorf("exp %v, got %v", actItems[0].ID, t2.ID)
	}

	t.Log("update first")
	if actErr := mem.Update(t1, time.Now()); actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	actItems, actErr = mem.Updated([]item.Kind{}, before)
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if len(actItems) != 2 {
		t.Errorf("exp 2, gor %d", len(actItems))
	}
	sort.Slice(actItems, func(i, j int) bool {
		return actItems[i].ID < actItems[j].ID
	})
	expItems := []item.Item{t1, t2}
	sort.Slice(expItems, func(i, j int) bool {
		return expItems[i].ID < expItems[j].ID
	})

	if actItems[0].ID != expItems[0].ID {
		t.Errorf("exp %v, got %v", actItems[0].ID, expItems[0].ID)
	}
	if actItems[1].ID != expItems[1].ID {
		t.Errorf("exp %v, got %v", actItems[1].ID, expItems[1].ID)
	}

	t.Log("select kind")
	actItems, actErr = mem.Updated([]item.Kind{"kinda"}, time.Time{})
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if len(actItems) != 1 {
		t.Errorf("exp 1, got %d", len(actItems))
	}
	if actItems[0].ID != t1.ID {
		t.Errorf("exp %v, got %v", t1.ID, actItems[0].ID)
	}
}

func TestMemoryRecur(t *testing.T) {
	t.Parallel()

	mem := NewMemory()
	now := time.Now()
	earlier := now.Add(-5 * time.Minute)
	today := item.NewDate(2024, 12, 1)
	yesterday := item.NewDate(2024, 11, 30)

	t.Log("start")
	i1 := item.Item{
		ID:        "a",
		Updated:   earlier,
		Recurrer:  item.NewRecurrer("2024-11-30, daily"),
		RecurNext: yesterday,
	}
	i2 := item.Item{
		ID:      "b",
		Updated: earlier,
	}

	for _, i := range []item.Item{i1, i2} {
		if err := mem.Update(i, i.Updated); err != nil {
			t.Errorf("exp nil, ot %v", err)
		}
	}

	t.Log("get recurrers")
	rs, err := mem.ShouldRecur(today)
	if err != nil {
		t.Errorf("exp nil, gt %v", err)
	}
	if diff := cmp.Diff([]item.Item{i1}, rs); diff != "" {
		t.Errorf("(exp +, got -)\n%s", diff)
	}

}
