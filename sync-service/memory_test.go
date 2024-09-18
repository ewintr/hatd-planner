package main

import (
	"testing"
	"time"
)

func TestMemoryItem(t *testing.T) {
	t.Parallel()

	mem := NewMemory()

	t.Log("start empty")
	actItems, actErr := mem.Updated([]Kind{}, time.Time{})
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if len(actItems) != 0 {
		t.Errorf("exp 0, got %d", len(actItems))
	}

	t.Log("add one")
	t1 := NewItem(Kind("kinda"), "test")
	if actErr := mem.Update(t1); actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	actItems, actErr = mem.Updated([]Kind{}, time.Time{})
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
	t2 := NewItem(Kind("kindb"), "test 2")
	if actErr := mem.Update(t2); actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	actItems, actErr = mem.Updated([]Kind{}, time.Time{})
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if len(actItems) != 2 {
		t.Errorf("exp 2, gor %d", len(actItems))
	}
	if actItems[0].ID != t1.ID {
		t.Errorf("exp %v, got %v", actItems[0].ID, t1.ID)
	}
	if actItems[1].ID != t2.ID {
		t.Errorf("exp %v, got %v", actItems[1].ID, t2.ID)
	}

	actItems, actErr = mem.Updated([]Kind{}, before)
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
	t1.Updated = time.Now()
	if actErr := mem.Update(t1); actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	actItems, actErr = mem.Updated([]Kind{}, before)
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if len(actItems) != 2 {
		t.Errorf("exp 2, gor %d", len(actItems))
	}
	if actItems[0].ID != t1.ID {
		t.Errorf("exp %v, got %v", actItems[0].ID, t1.ID)
	}
	if actItems[1].ID != t2.ID {
		t.Errorf("exp %v, got %v", actItems[1].ID, t2.ID)
	}

	t.Log("select kind")
	actItems, actErr = mem.Updated([]Kind{"kinda"}, time.Time{})
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
