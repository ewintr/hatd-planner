package memory_test

import (
	"fmt"
	"testing"
	"time"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage/memory"
)

func TestSync(t *testing.T) {
	t.Parallel()

	mem := memory.NewSync()

	t.Log("store")
	now := time.Now()
	ts := now
	count := 3
	for i := 0; i < count; i++ {
		mem.Store(item.Item{
			ID:      fmt.Sprintf("id-%d", i),
			Updated: ts,
		})
		ts = ts.Add(-1 * time.Minute)
	}

	t.Log("find all")
	actItems, actErr := mem.FindAll()
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if len(actItems) != count {
		t.Errorf("exp %v, got %v", count, len(actItems))
	}

	t.Log("last update")
	actLU, actErr := mem.LastUpdate()
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if !actLU.Equal(now) {
		t.Errorf("exp %v, got %v", now, actLU)
	}

	t.Log("delete all")
	if err := mem.DeleteAll(); err != nil {
		t.Errorf("exp nil, got %v", err)
	}
	actItems, actErr = mem.FindAll()
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if len(actItems) != 0 {
		t.Errorf("exp 0, got %v", len(actItems))
	}

}
