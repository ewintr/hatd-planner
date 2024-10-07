package command_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/command"
	"go-mod.ewintr.nl/planner/plan/storage/memory"
)

func TestUpdate(t *testing.T) {
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
		args     map[string]string
		expEvent item.Event
		expErr   bool
	}{
		{
			name:    "no args",
			localID: lid,
			expEvent: item.Event{
				ID: eid,
				EventBody: item.EventBody{
					Title:    title,
					Start:    start,
					Duration: oneHour,
				},
			},
		},
		{
			name:    "not found",
			localID: 1,
			expErr:  true,
		},
		{
			name:    "name",
			localID: lid,
			args: map[string]string{
				"name": "updated",
			},
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
			args: map[string]string{
				"on": "invalid",
			},
			expErr: true,
		},
		{
			name:    "on",
			localID: lid,
			args: map[string]string{
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
			args: map[string]string{
				"at": "invalid",
			},
			expErr: true,
		},
		{
			name:    "at",
			localID: lid,
			args: map[string]string{
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
			args: map[string]string{
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
			args: map[string]string{
				"for": "invalid",
			},
			expErr: true,
		},
		{
			name:    "for",
			localID: lid,
			args: map[string]string{
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

			actErr := command.Update(localIDRepo, eventRepo, syncRepo, tc.localID, tc.args["name"], tc.args["on"], tc.args["at"], tc.args["for"]) != nil
			if tc.expErr != actErr {
				t.Errorf("exp %v, got %v", tc.expErr, actErr)
			}
			if tc.expErr {
				return
			}

			actEvent, err := eventRepo.Find(eid)
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if diff := cmp.Diff(tc.expEvent, actEvent); diff != "" {
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
