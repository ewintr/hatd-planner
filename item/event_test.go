package item_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go-mod.ewintr.nl/planner/item"
)

func TestNewEvent(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		it       item.Item
		expEvent item.Event
		expErr   bool
	}{
		{
			name: "wrong kind",
			it: item.Item{
				ID:   "a",
				Kind: item.KindTask,
				Body: `{
  "title":"title",
  "start":"2024-09-20T08:00:00Z",
  "end":"2024-09-20T10:00:00Z" 
}`,
			},
			expErr: true,
		},
		{
			name: "invalid json",
			it: item.Item{
				ID:   "a",
				Kind: item.KindEvent,
				Body: `{"id":"a"`,
			},
			expErr: true,
		},
		{
			name: "valid",
			it: item.Item{
				ID:   "a",
				Kind: item.KindEvent,
				Body: `{
  "title":"title",
  "start":"2024-09-20T08:00:00Z",
  "end":"2024-09-20T10:00:00Z" 
}`,
			},
			expEvent: item.Event{
				ID: "a",
				EventBody: item.EventBody{
					Title: "title",
					Start: time.Date(2024, 9, 20, 8, 0, 0, 0, time.UTC),
					End:   time.Date(2024, 9, 20, 10, 0, 0, 0, time.UTC),
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actEvent, actErr := item.NewEvent(tc.it)
			if tc.expErr != (actErr != nil) {
				t.Errorf("exp nil, got %v", actErr)
			}
			if tc.expErr {
				return
			}
			if diff := cmp.Diff(tc.expEvent, actEvent); diff != "" {
				t.Errorf("(exp +, got -)\n%s", diff)
			}
		})
	}
}

func TestEventItem(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		event   item.Event
		expItem item.Item
		expErr  bool
	}{
		{
			name: "empty",
			expItem: item.Item{
				Kind:    item.KindEvent,
				Updated: time.Time{},
				Body:    `{"start":"0001-01-01T00:00:00Z","end":"0001-01-01T00:00:00Z","title":""}`,
			},
		},
		{
			name: "normal",
			event: item.Event{
				ID: "a",
				EventBody: item.EventBody{
					Title: "title",
					Start: time.Date(2024, 9, 23, 8, 0, 0, 0, time.UTC),
					End:   time.Date(2024, 9, 23, 10, 0, 0, 0, time.UTC),
				},
			},
			expItem: item.Item{
				ID:      "a",
				Kind:    item.KindEvent,
				Updated: time.Time{},
				Body:    `{"start":"2024-09-23T08:00:00Z","end":"2024-09-23T10:00:00Z","title":"title"}`,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actItem, actErr := tc.event.Item()
			if tc.expErr != (actErr != nil) {
				t.Errorf("exp nil, got %v", actErr)
			}
			if tc.expErr {
				return
			}
			if diff := cmp.Diff(tc.expItem, actItem); diff != "" {
				t.Errorf("(exp+, got -)\n%s", diff)
			}
		})
	}
}
