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

	aDateStr := "2024-11-02"
	aDate := time.Date(2024, 11, 2, 0, 0, 0, 0, time.UTC)
	aTimeStr := "12:00"
	aDay := time.Duration(24) * time.Hour
	anHourStr := "1h"
	anHour := time.Hour
	aDateAndTime := time.Date(2024, 11, 2, 12, 0, 0, 0, time.UTC)

	for _, tc := range []struct {
		name     string
		main     []string
		flags    map[string]string
		expErr   bool
		expEvent item.Event
	}{
		{
			name:   "empty",
			expErr: true,
		},
		{
			name: "title missing",
			main: []string{"add"},
			flags: map[string]string{
				command.FlagOn: aDateStr,
			},
			expErr: true,
		},
		{
			name:   "date missing",
			main:   []string{"add", "some", "title"},
			expErr: true,
		},
		{
			name: "only date",
			main: []string{"add", "title"},
			flags: map[string]string{
				command.FlagOn: aDateStr,
			},
			expEvent: item.Event{
				ID: "title",
				EventBody: item.EventBody{
					Title:    "title",
					Start:    aDate,
					Duration: aDay,
				},
			},
		},
		{
			name: "date and time",
			main: []string{"add", "title"},
			flags: map[string]string{
				command.FlagOn: aDateStr,
				command.FlagAt: aTimeStr,
			},
			expEvent: item.Event{
				ID: "title",
				EventBody: item.EventBody{
					Title:    "title",
					Start:    aDateAndTime,
					Duration: anHour,
				},
			},
		},
		{
			name: "date, time and duration",
			main: []string{"add", "title"},
			flags: map[string]string{
				command.FlagOn:  aDateStr,
				command.FlagAt:  aTimeStr,
				command.FlagFor: anHourStr,
			},
			expEvent: item.Event{
				ID: "title",
				EventBody: item.EventBody{
					Title:    "title",
					Start:    aDateAndTime,
					Duration: anHour,
				},
			},
		},
		{
			name: "date and duration",
			main: []string{"add", "title"},
			flags: map[string]string{
				command.FlagOn:  aDateStr,
				command.FlagFor: anHourStr,
			},
			expErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			eventRepo := memory.NewEvent()
			localRepo := memory.NewLocalID()
			syncRepo := memory.NewSync()
			cmd := command.NewAdd(localRepo, eventRepo, syncRepo)
			actParseErr := cmd.Execute(tc.main, tc.flags) != nil
			if tc.expErr != actParseErr {
				t.Errorf("exp %v, got %v", tc.expErr, actParseErr)
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

			updated, err := syncRepo.FindAll()
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if len(updated) != 1 {
				t.Errorf("exp 1, got %v", len(updated))
			}
		})
	}
}
