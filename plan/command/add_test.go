package command_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/command"
	"go-mod.ewintr.nl/planner/plan/storage/memory"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	oneHour, err := time.ParseDuration("1h")
	if err != nil {
		t.Errorf("exp nil, got %v", err)
	}
	oneDay, err := time.ParseDuration("24h")
	if err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	for _, tc := range []struct {
		name     string
		args     map[string]string
		expEvent item.Event
		expErr   bool
	}{
		{
			name: "no name",
			args: map[string]string{
				"on":  "2024-10-01",
				"at":  "9:00",
				"for": "1h",
			},
			expErr: true,
		},
		{
			name: "no date",
			args: map[string]string{
				"name": "event",
				"at":   "9:00",
				"for":  "1h",
			},
			expErr: true,
		},
		{
			name: "duration, but no time",
			args: map[string]string{
				"name": "event",
				"on":   "2024-10-01",
				"for":  "1h",
			},
			expErr: true,
		},
		{
			name: "time, but no duration",
			args: map[string]string{
				"name": "event",
				"on":   "2024-10-01",
				"at":   "9:00",
			},
			expEvent: item.Event{
				ID: "a",
				EventBody: item.EventBody{
					Title: "event",
					Start: time.Date(2024, 10, 1, 9, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "no time, no duration",
			args: map[string]string{
				"name": "event",
				"on":   "2024-10-01",
			},
			expEvent: item.Event{
				ID: "a",
				EventBody: item.EventBody{
					Title:    "event",
					Start:    time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
					Duration: oneDay,
				},
			},
		},
		{
			name: "full",
			args: map[string]string{
				"name": "event",
				"on":   "2024-10-01",
				"at":   "9:00",
				"for":  "1h",
			},
			expEvent: item.Event{
				ID: "a",
				EventBody: item.EventBody{
					Title:    "event",
					Start:    time.Date(2024, 10, 1, 9, 0, 0, 0, time.UTC),
					Duration: oneHour,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			eventRepo := memory.NewEvent()
			localRepo := memory.NewLocalID()
			actErr := command.Add(localRepo, eventRepo, tc.args["name"], tc.args["on"], tc.args["at"], tc.args["for"]) != nil
			if tc.expErr != actErr {
				t.Errorf("exp %v, got %v", tc.expErr, actErr)
			}
			if tc.expErr {
				return
			}
			actEvents, err := eventRepo.FindAll()
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if len(actEvents) != 1 {
				t.Errorf("exp 1, got %d", len(actEvents))
			}

			actLocalIDs, err := localRepo.FindAll()
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if len(actLocalIDs) != 1 {
				t.Errorf("exp 1, got %v", len(actLocalIDs))
			}
			if _, ok := actLocalIDs[actEvents[0].ID]; !ok {
				t.Errorf("exp true, got %v", ok)
			}

			if actEvents[0].ID == "" {
				t.Errorf("exp string not te be empty")
			}
			tc.expEvent.ID = actEvents[0].ID
			if diff := cmp.Diff(tc.expEvent, actEvents[0]); diff != "" {
				t.Errorf("(exp +, got -)\n%s", diff)
			}
		})
	}
}
