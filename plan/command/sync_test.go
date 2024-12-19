package command_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/command"
	"go-mod.ewintr.nl/planner/plan/storage/memory"
	"go-mod.ewintr.nl/planner/sync/client"
)

func TestSyncParse(t *testing.T) {
	t.Parallel()

	syncClient := client.NewMemory()
	syncRepo := memory.NewSync()
	localIDRepo := memory.NewLocalID()
	eventRepo := memory.NewEvent()

	for _, tc := range []struct {
		name   string
		main   []string
		expErr bool
	}{
		{
			name:   "empty",
			expErr: true,
		},
		{
			name:   "wrong",
			main:   []string{"wrong"},
			expErr: true,
		},
		{
			name: "valid",
			main: []string{"sync"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd := command.NewSync(syncClient, syncRepo, localIDRepo, eventRepo)
			actErr := cmd.Execute(tc.main, nil) != nil
			if tc.expErr != actErr {
				t.Errorf("exp %v, got %v", tc.expErr, actErr)
			}
		})
	}
}

func TestSyncSend(t *testing.T) {
	t.Parallel()

	syncClient := client.NewMemory()
	syncRepo := memory.NewSync()
	localIDRepo := memory.NewLocalID()
	eventRepo := memory.NewEvent()

	it := item.Item{
		ID:   "a",
		Kind: item.KindEvent,
		Body: `{
  "title":"title",
  "start":"2024-10-18T08:00:00Z",
  "duration":"1h" 
}`,
	}
	if err := syncRepo.Store(it); err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	for _, tc := range []struct {
		name     string
		ks       []item.Kind
		ts       time.Time
		expItems []item.Item
	}{
		{
			name:     "single",
			ks:       []item.Kind{item.KindEvent},
			expItems: []item.Item{it},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd := command.NewSync(syncClient, syncRepo, localIDRepo, eventRepo)
			if err := cmd.Execute([]string{"sync"}, nil); err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			actItems, actErr := syncClient.Updated(tc.ks, tc.ts)
			if actErr != nil {
				t.Errorf("exp nil, got %v", actErr)
			}
			if diff := cmp.Diff(tc.expItems, actItems); diff != "" {
				t.Errorf("(exp +, got -)\n%s", diff)
			}

			actLeft, actErr := syncRepo.FindAll()
			if actErr != nil {
				t.Errorf("exp nil, got %v", actErr)
			}
			if len(actLeft) != 0 {
				t.Errorf("exp 0, got %v", actLeft)
			}
		})
	}
}

func TestSyncReceive(t *testing.T) {
	t.Parallel()

	oneHour, err := time.ParseDuration("1h")
	if err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	for _, tc := range []struct {
		name       string
		present    []item.Event
		updated    []item.Item
		expEvent   []item.Event
		expLocalID map[string]int
	}{
		{
			name:       "no new",
			expEvent:   []item.Event{},
			expLocalID: map[string]int{},
		},
		{
			name: "new",
			updated: []item.Item{{
				ID:   "a",
				Kind: item.KindEvent,
				Body: `{
  "title":"title",
  "start":"2024-10-23T08:00:00Z",
  "duration":"1h" 
}`,
			}},
			expEvent: []item.Event{{
				ID:   "a",
				Date: item.NewDate(2024, 10, 23),
				EventBody: item.EventBody{
					Title:    "title",
					Duration: oneHour,
				},
			}},
			expLocalID: map[string]int{
				"a": 1,
			},
		},
		{
			name: "update existing",
			present: []item.Event{{
				ID:   "a",
				Date: item.NewDate(2024, 10, 23),
				EventBody: item.EventBody{
					Title:    "title",
					Duration: oneHour,
				},
			}},
			updated: []item.Item{{
				ID:   "a",
				Kind: item.KindEvent,
				Body: `{
  "title":"new title",
  "start":"2024-10-23T08:00:00Z",
  "duration":"1h" 
}`,
			}},
			expEvent: []item.Event{{
				ID:   "a",
				Date: item.NewDate(2024, 10, 23),
				EventBody: item.EventBody{
					Title:    "new title",
					Duration: oneHour,
				},
			}},
			expLocalID: map[string]int{
				"a": 1,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// setup
			syncClient := client.NewMemory()
			syncRepo := memory.NewSync()
			localIDRepo := memory.NewLocalID()
			eventRepo := memory.NewEvent()

			for i, p := range tc.present {
				if err := eventRepo.Store(p); err != nil {
					t.Errorf("exp nil, got %v", err)
				}
				if err := localIDRepo.Store(p.ID, i+1); err != nil {
					t.Errorf("exp nil, got %v", err)
				}
			}
			if err := syncClient.Update(tc.updated); err != nil {
				t.Errorf("exp nil, got %v", err)
			}

			// sync
			cmd := command.NewSync(syncClient, syncRepo, localIDRepo, eventRepo)
			if err := cmd.Execute([]string{"sync"}, nil); err != nil {
				t.Errorf("exp nil, got %v", err)
			}

			// check result
			actEvents, err := eventRepo.FindAll()
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if diff := item.EventDiffs(tc.expEvent, actEvents); diff != "" {
				t.Errorf("(exp +, got -)\n%s", diff)
			}
			actLocalIDs, err := localIDRepo.FindAll()
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if diff := cmp.Diff(tc.expLocalID, actLocalIDs); diff != "" {
				t.Errorf("(exp +, got -)\n%s", diff)
			}
		})
	}
}
