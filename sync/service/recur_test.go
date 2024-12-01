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

	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	mem := NewMemory()
	rec := NewRecur(mem, mem, slog.New(slog.NewTextHandler(io.Discard, nil)))

	// Create a recurring item
	recur := &item.Recur{
		Start:  now,
		Period: item.PeriodDay,
		Count:  1,
	}
	testItem := item.Item{
		ID:        "test-1",
		Kind:      item.KindEvent,
		Updated:   now,
		Deleted:   false,
		Recurrer:  recur,
		RecurNext: now,
		Body:      `{"title":"Test Event","start":"2024-01-01T10:00:00Z","duration":"30m"}`,
	}

	// Store the item
	if err := mem.Update(testItem, testItem.Updated); err != nil {
		t.Fatalf("failed to store test item: %v", err)
	}

	// Run recurrence
	if err := rec.Recur(); err != nil {
		t.Errorf("Recur failed: %v", err)
	}

	// Verify results
	items, err := mem.Updated([]item.Kind{item.KindEvent}, now)
	if err != nil {
		t.Errorf("failed to get updated items: %v", err)
	}

	if len(items) != 2 { // Original + new instance
		t.Errorf("expected 2 items, got %d", len(items))
	}

	// Check that RecurNext was updated
	recurItems, err := mem.RecursBefore(now.Add(48 * time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if len(recurItems) != 1 {
		t.Errorf("expected 1 recur item, got %d", len(recurItems))
	}
	if !recurItems[0].RecurNext.After(now) {
		t.Errorf("RecurNext was not updated, still %v", recurItems[0].RecurNext)
	}
}
