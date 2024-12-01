package command_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/command"
	"go-mod.ewintr.nl/planner/plan/storage/memory"
)

func TestUpdateExecute(t *testing.T) {
	t.Parallel()

	eid := "c"
	lid := 3
	oneHour, err := time.ParseDuration("1h")
	if err != nil {
		t.Errorf("exp nil, got %v", err)
	}
	title := "title"
	start := time.Date(2024, 10, 6, 10, 0, 0, 0, time.UTC)
	twoHour, err := time.ParseDuration("2h")
	if err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	for _, tc := range []struct {
		name     string
		localID  int
		main     []string
		flags    map[string]string
		expEvent item.Event
		expErr   bool
	}{
		{
			name:   "no args",
			expErr: true,
		},
		{
			name:    "not found",
			localID: 1,
			expErr:  true,
		},
		{
			name:    "name",
			localID: lid,
			main:    []string{"update", fmt.Sprintf("%d", lid), "updated"},
			expEvent: item.Event{
				ID: eid,
				EventBody: item.EventBody{
					Title:    "updated",
					Start:    start,
					Duration: oneHour,
				},
			},
		},
		{
			name:    "invalid on",
			localID: lid,
			main:    []string{"update", fmt.Sprintf("%d", lid)},
			flags: map[string]string{
				"on": "invalid",
			},
			expErr: true,
		},
		{
			name:    "on",
			localID: lid,
			main:    []string{"update", fmt.Sprintf("%d", lid)},
			flags: map[string]string{
				"on": "2024-10-02",
			},
			expEvent: item.Event{
				ID: eid,
				EventBody: item.EventBody{
					Title:    title,
					Start:    time.Date(2024, 10, 2, 10, 0, 0, 0, time.UTC),
					Duration: oneHour,
				},
			},
		},
		{
			name:    "invalid at",
			localID: lid,
			main:    []string{"update", fmt.Sprintf("%d", lid)},
			flags: map[string]string{
				"at": "invalid",
			},
			expErr: true,
		},
		{
			name:    "at",
			localID: lid,
			main:    []string{"update", fmt.Sprintf("%d", lid)},
			flags: map[string]string{
				"at": "11:00",
			},
			expEvent: item.Event{
				ID: eid,
				EventBody: item.EventBody{
					Title:    title,
					Start:    time.Date(2024, 10, 6, 11, 0, 0, 0, time.UTC),
					Duration: oneHour,
				},
			},
		},
		{
			name:    "on and at",
			localID: lid,
			main:    []string{"update", fmt.Sprintf("%d", lid)},
			flags: map[string]string{
				"on": "2024-10-02",
				"at": "11:00",
			},
			expEvent: item.Event{
				ID: eid,
				EventBody: item.EventBody{
					Title:    title,
					Start:    time.Date(2024, 10, 2, 11, 0, 0, 0, time.UTC),
					Duration: oneHour,
				},
			},
		},
		{
			name:    "invalid for",
			localID: lid,
			main:    []string{"update", fmt.Sprintf("%d", lid)},
			flags: map[string]string{
				"for": "invalid",
			},
			expErr: true,
		},
		{
			name:    "for",
			localID: lid,
			main:    []string{"update", fmt.Sprintf("%d", lid)},
			flags: map[string]string{
				"for": "2h",
			},
			expEvent: item.Event{
				ID: eid,
				EventBody: item.EventBody{
					Title:    title,
					Start:    time.Date(2024, 10, 6, 10, 0, 0, 0, time.UTC),
					Duration: twoHour,
				},
			},
		},
		{
			name: "invalid rec start",
			main: []string{"update", fmt.Sprintf("%d", lid)},
			flags: map[string]string{
				"rec-start": "invalud",
			},
			expErr: true,
		},
		{
			name: "valid rec start",
			main: []string{"update", fmt.Sprintf("%d", lid)},
			flags: map[string]string{
				"rec-start": "2024-12-08",
			},
			expEvent: item.Event{
				ID: eid,
				Recurrer: &item.Recur{
					Start: time.Date(2024, 12, 8, 0, 0, 0, 0, time.UTC),
				},
				EventBody: item.EventBody{
					Title:    title,
					Start:    start,
					Duration: oneHour,
				},
			},
		},
		{
			name: "invalid rec period",
			main: []string{"update", fmt.Sprintf("%d", lid)},
			flags: map[string]string{
				"rec-period": "invalid",
			},
			expErr: true,
		},
		{
			name: "valid rec period",
			main: []string{"update", fmt.Sprintf("%d", lid)},
			flags: map[string]string{
				"rec-period": "month",
			},
			expEvent: item.Event{
				ID: eid,
				Recurrer: &item.Recur{
					Period: item.PeriodMonth,
				},
				EventBody: item.EventBody{
					Title:    title,
					Start:    start,
					Duration: oneHour,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			eventRepo := memory.NewEvent()
			localIDRepo := memory.NewLocalID()
			syncRepo := memory.NewSync()
			if err := eventRepo.Store(item.Event{
				ID: eid,
				EventBody: item.EventBody{
					Title:    title,
					Start:    start,
					Duration: oneHour,
				},
			}); err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if err := localIDRepo.Store(eid, lid); err != nil {
				t.Errorf("exp nil, ,got %v", err)
			}

			cmd := command.NewUpdate(localIDRepo, eventRepo, syncRepo)
			actParseErr := cmd.Execute(tc.main, tc.flags) != nil
			if tc.expErr != actParseErr {
				t.Errorf("exp %v, got %v", tc.expErr, actParseErr)
			}
			if tc.expErr {
				return
			}

			actEvent, err := eventRepo.Find(eid)
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if diff := cmp.Diff(tc.expEvent, actEvent); diff != "" {
				t.Errorf("(exp -, got +)\n%s", diff)
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
