package command_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/command"
	"go-mod.ewintr.nl/planner/plan/storage"
	"go-mod.ewintr.nl/planner/plan/storage/memory"
	"go-mod.ewintr.nl/planner/sync/client"
)

func TestSyncParse(t *testing.T) {
	t.Parallel()

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
			_, actErr := command.SyncArgs{}.Parse(tc.main, nil)
			if tc.expErr != (actErr != nil) {
				t.Errorf("exp %v, got %v", tc.expErr, actErr)
			}
		})
	}
}

func TestSyncSend(t *testing.T) {
	t.Parallel()

	syncClient := client.NewMemory()
	mems := memory.New()

	it := item.Item{
		ID:   "a",
		Kind: item.KindTask,
		Body: `{
  "title":"title",
  "start":"2024-10-18T08:00:00Z",
  "duration":"1h" 
}`,
	}
	if err := mems.Sync(nil).Store(it); err != nil {
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
			ks:       []item.Kind{item.KindTask},
			expItems: []item.Item{it},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd, err := command.SyncArgs{}.Parse([]string{"sync"}, nil)
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if _, err := cmd.Do(mems, syncClient); err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			actItems, actErr := syncClient.Updated(tc.ks, tc.ts)
			if actErr != nil {
				t.Errorf("exp nil, got %v", actErr)
			}
			if diff := cmp.Diff(tc.expItems, actItems); diff != "" {
				t.Errorf("(exp +, got -)\n%s", diff)
			}

			actLeft, actErr := mems.Sync(nil).FindAll()
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
		present    []item.Task
		updated    []item.Item
		expTask    []item.Task
		expLocalID map[string]int
	}{
		{
			name:       "no new",
			expTask:    []item.Task{},
			expLocalID: map[string]int{},
		},
		{
			name: "new",
			updated: []item.Item{{
				ID:   "a",
				Kind: item.KindTask,
				Body: `{
  "title":"title",
  "start":"2024-10-23T08:00:00Z",
  "duration":"1h" 
}`,
			}},
			expTask: []item.Task{{
				ID:   "a",
				Date: item.NewDate(2024, 10, 23),
				TaskBody: item.TaskBody{
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
			present: []item.Task{{
				ID:   "a",
				Date: item.NewDate(2024, 10, 23),
				TaskBody: item.TaskBody{
					Title:    "title",
					Duration: oneHour,
				},
			}},
			updated: []item.Item{{
				ID:   "a",
				Kind: item.KindTask,
				Body: `{
  "title":"new title",
  "start":"2024-10-23T08:00:00Z",
  "duration":"1h" 
}`,
			}},
			expTask: []item.Task{{
				ID:   "a",
				Date: item.NewDate(2024, 10, 23),
				TaskBody: item.TaskBody{
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
			syncClient := client.NewMemory()
			mems := memory.New()

			// setup
			for i, p := range tc.present {
				if err := mems.Task(nil).Store(p); err != nil {
					t.Errorf("exp nil, got %v", err)
				}
				if err := mems.LocalID(nil).Store(p.ID, i+1); err != nil {
					t.Errorf("exp nil, got %v", err)
				}
			}
			if err := syncClient.Update(tc.updated); err != nil {
				t.Errorf("exp nil, got %v", err)
			}

			// sync
			cmd, err := command.NewSyncArgs().Parse([]string{"sync"}, nil)
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if _, err := cmd.Do(mems, syncClient); err != nil {
				t.Errorf("exp nil, got %v", err)
			}

			// check result
			actTasks, err := mems.Task(nil).FindMany(storage.TaskListParams{})
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if diff := item.TaskDiffs(tc.expTask, actTasks); diff != "" {
				t.Errorf("(exp +, got -)\n%s", diff)
			}
			actLocalIDs, err := mems.LocalID(nil).FindAll()
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if diff := cmp.Diff(tc.expLocalID, actLocalIDs); diff != "" {
				t.Errorf("(exp +, got -)\n%s", diff)
			}
		})
	}
}
