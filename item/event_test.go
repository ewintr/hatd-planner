package item_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go-mod.ewintr.nl/planner/item"
)

func TestNewEvent(t *testing.T) {
	t.Parallel()

	oneHour, err := time.ParseDuration("1h")
	if err != nil {
		t.Errorf("exp nil, got %v", err)
	}
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
  "duration":"1h" 
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
				Recurrer: &item.Recur{
					Start:  time.Date(2024, 12, 8, 9, 0, 0, 0, time.UTC),
					Period: item.PeriodDay,
					Count:  1,
				},
				Body: `{
  "title":"title",
  "start":"2024-09-20T08:00:00Z",
  "duration":"1h" 
}`,
			},
			expEvent: item.Event{
				ID: "a",
				Recurrer: &item.Recur{
					Start:  time.Date(2024, 12, 8, 9, 0, 0, 0, time.UTC),
					Period: item.PeriodDay,
					Count:  1,
				},
				EventBody: item.EventBody{
					Title:    "title",
					Start:    time.Date(2024, 9, 20, 8, 0, 0, 0, time.UTC),
					Duration: oneHour,
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

	oneHour, err := time.ParseDuration("1h")
	if err != nil {
		t.Errorf("exp nil, got %v", err)
	}
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
				Body:    `{"start":"0001-01-01T00:00:00Z","duration":"0s","title":""}`,
			},
		},
		{
			name: "normal",
			event: item.Event{
				ID: "a",
				EventBody: item.EventBody{
					Title:    "title",
					Start:    time.Date(2024, 9, 23, 8, 0, 0, 0, time.UTC),
					Duration: oneHour,
				},
			},
			expItem: item.Item{
				ID:      "a",
				Kind:    item.KindEvent,
				Updated: time.Time{},
				Body:    `{"start":"2024-09-23T08:00:00Z","duration":"1h0m0s","title":"title"}`,
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

func TestEventValidate(t *testing.T) {
	t.Parallel()

	oneHour, err := time.ParseDuration("1h")
	if err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	for _, tc := range []struct {
		name  string
		event item.Event
		exp   bool
	}{
		{
			name: "empty",
		},
		{
			name: "missing title",
			event: item.Event{
				ID: "a",
				EventBody: item.EventBody{
					Start:    time.Date(2024, 9, 20, 8, 0, 0, 0, time.UTC),
					Duration: oneHour,
				},
			},
		},
		{
			name: "no date",
			event: item.Event{
				ID: "a",
				EventBody: item.EventBody{
					Title:    "title",
					Start:    time.Date(0, 0, 0, 8, 0, 0, 0, time.UTC),
					Duration: oneHour,
				},
			},
		},
		{
			name: "no duration",
			event: item.Event{
				ID: "a",
				EventBody: item.EventBody{
					Title: "title",
					Start: time.Date(2024, 9, 20, 8, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "valid",
			event: item.Event{
				ID: "a",
				EventBody: item.EventBody{
					Title:    "title",
					Start:    time.Date(2024, 9, 20, 8, 0, 0, 0, time.UTC),
					Duration: oneHour,
				},
			},
			exp: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if act := tc.event.Valid(); tc.exp != act {
				t.Errorf("exp %v, got %v", tc.exp, act)
			}

		})
	}
}
