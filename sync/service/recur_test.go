package main

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"go-mod.ewintr.nl/planner/item"
)

func TestRecur(t *testing.T) {
	t.Parallel()

	now := time.Now()
	today := item.NewDate(2024, 1, 1)
	until := today.Add(3)
	mem := NewMemory()
	rec := NewRecur(mem, mem, slog.New(slog.NewTextHandler(io.Discard, nil)))

	testItem := item.Item{
		ID:        "test-1",
		Kind:      item.KindTask,
		Updated:   now,
		Deleted:   false,
		Recurrer:  item.NewRecurrer("2024-01-01, daily"),
		RecurNext: today,
		Body:      `{"title":"Test task","start":"2024-01-01T10:00:00Z","duration":"30m"}`,
	}

	// Store the item
	if err := mem.Update(testItem, testItem.Updated); err != nil {
		t.Fatalf("failed to store test item: %v", err)
	}

	// Run recurrence
	if err := rec.Recur(until); err != nil {
		t.Errorf("Recur failed: %v", err)
	}

	// Verify results
	items, err := mem.Updated([]item.Kind{item.KindTask}, now)
	if err != nil {
		t.Errorf("failed to get updated items: %v", err)
	}

	if len(items) != 5 { // Original + 4 new instances
		t.Errorf("expected 5 items, got %d", len(items))
	}

	// Check that RecurNext was updated
	recurItems, err := mem.ShouldRecur(until.Add(1))
	if err != nil {
		t.Errorf("exp nil, got %v", err)
	}
	if len(recurItems) != 1 {
		t.Errorf("expected 1 recur item, got %d", len(recurItems))
	}
	if !recurItems[0].RecurNext.After(today) {
		t.Errorf("RecurNext was not updated, still %v", recurItems[0].RecurNext)
	}
}
