package storage

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go-mod.ewintr.nl/planner/item"
)

func TestMemory(t *testing.T) {
	t.Parallel()

	mem := NewMemory()

	t.Log("empty")
	actEvents, actErr := mem.FindAll()
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if len(actEvents) != 0 {
		t.Errorf("exp 0, got %d", len(actEvents))
	}

	t.Log("store")
	e1 := item.Event{
		ID: "id-1",
	}
	if err := mem.Store(e1); err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	e2 := item.Event{
		ID: "id-2",
	}
	if err := mem.Store(e2); err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	t.Log("find one")
	actEvent, actErr := mem.Find(e1.ID)
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if actEvent.ID != e1.ID {
		t.Errorf("exp %v, got %v", e1.ID, actEvent.ID)
	}

	t.Log("find all")
	actEvents, actErr = mem.FindAll()
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if diff := cmp.Diff([]item.Event{e1, e2}, actEvents); diff != "" {
		t.Errorf("(exp -, got +)\n%s", diff)
	}
}
