package command_test

import (
	"fmt"
	"testing"
	"time"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/command"
	"go-mod.ewintr.nl/planner/plan/storage/memory"
)

func TestUpdateExecute(t *testing.T) {
	t.Parallel()

	tskID := "c"
	lid := 3
	oneHour, err := time.ParseDuration("1h")
	if err != nil {
		t.Errorf("exp nil, got %v", err)
	}
	title := "title"
	project := "project"
	aDate := item.NewDate(2024, 10, 6)
	aTime := item.NewTime(10, 0)
	twoHour, err := time.ParseDuration("2h")
	if err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	for _, tc := range []struct {
		name        string
		localID     int
		main        []string
		fields      map[string]string
		expTask     item.Task
		expParseErr bool
		expDoErr    bool
	}{
		{
			name:        "no args",
			expParseErr: true,
		},
		{
			name:     "not found",
			main:     []string{"update", "1"},
			expDoErr: true,
		},
		{
			name:    "name",
			localID: lid,
			main:    []string{"update", fmt.Sprintf("%d", lid), "updated"},
			expTask: item.Task{
				ID:   tskID,
				Date: item.NewDate(2024, 10, 6),
				TaskBody: item.TaskBody{
					Title:    "updated",
					Project:  project,
					Time:     aTime,
					Duration: oneHour,
				},
			},
		},
		{
			name:    "project",
			localID: lid,
			main:    []string{"update", fmt.Sprintf("%d", lid)},
			fields: map[string]string{
				"p": "updated",
			},
			expTask: item.Task{
				ID:   tskID,
				Date: item.NewDate(2024, 10, 2),
				TaskBody: item.TaskBody{
					Title:    title,
					Project:  "updated",
					Time:     aTime,
					Duration: oneHour,
				},
			},
		},
		{
			name:    "invalid date",
			localID: lid,
			main:    []string{"update", fmt.Sprintf("%d", lid)},
			fields: map[string]string{
				"on": "invalid",
			},
			expParseErr: true,
		},
		{
			name:    "date",
			localID: lid,
			main:    []string{"update", fmt.Sprintf("%d", lid)},
			fields: map[string]string{
				"on": "2024-10-02",
			},
			expTask: item.Task{
				ID:   tskID,
				Date: item.NewDate(2024, 10, 2),
				TaskBody: item.TaskBody{
					Title:    title,
					Project:  project,
					Time:     aTime,
					Duration: oneHour,
				},
			},
		},
		{
			name:    "invalid time",
			localID: lid,
			main:    []string{"update", fmt.Sprintf("%d", lid)},
			fields: map[string]string{
				"at": "invalid",
			},
			expParseErr: true,
		},
		{
			name:    "time",
			localID: lid,
			main:    []string{"update", fmt.Sprintf("%d", lid)},
			fields: map[string]string{
				"time": "11:00",
			},
			expTask: item.Task{
				ID:   tskID,
				Date: item.NewDate(2024, 10, 6),
				TaskBody: item.TaskBody{
					Title:    title,
					Project:  project,
					Time:     item.NewTime(11, 0),
					Duration: oneHour,
				},
			},
		},
		{
			name:    "invalid duration",
			localID: lid,
			main:    []string{"update", fmt.Sprintf("%d", lid)},
			fields: map[string]string{
				"for": "invalid",
			},
			expParseErr: true,
		},
		{
			name:    "duration",
			localID: lid,
			main:    []string{"update", fmt.Sprintf("%d", lid)},
			fields: map[string]string{
				"for": "2h",
			},
			expTask: item.Task{
				ID:   tskID,
				Date: item.NewDate(2024, 10, 6),
				TaskBody: item.TaskBody{
					Title:    title,
					Project:  project,
					Time:     aTime,
					Duration: twoHour,
				},
			},
		},
		{
			name: "invalid recurrer",
			main: []string{"update", fmt.Sprintf("%d", lid)},
			fields: map[string]string{
				"rec": "invalud",
			},
			expParseErr: true,
		},
		{
			name: "valid recurrer",
			main: []string{"update", fmt.Sprintf("%d", lid)},
			fields: map[string]string{
				"rec": "2024-12-08, daily",
			},
			expTask: item.Task{
				ID:       tskID,
				Date:     aDate,
				Recurrer: item.NewRecurrer("2024-12-08, daily"),
				TaskBody: item.TaskBody{
					Title:    title,
					Project:  project,
					Time:     aTime,
					Duration: oneHour,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// setup
			taskRepo := memory.NewTask()
			localIDRepo := memory.NewLocalID()
			syncRepo := memory.NewSync()
			if err := taskRepo.Store(item.Task{
				ID:   tskID,
				Date: aDate,
				TaskBody: item.TaskBody{
					Title:    title,
					Project:  project,
					Time:     aTime,
					Duration: oneHour,
				},
			}); err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if err := localIDRepo.Store(tskID, lid); err != nil {
				t.Errorf("exp nil, ,got %v", err)
			}

			// parse
			cmd, actErr := command.NewUpdateArgs().Parse(tc.main, tc.fields)
			if tc.expParseErr != (actErr != nil) {
				t.Errorf("exp %v, got %v", tc.expParseErr, actErr)
			}
			if tc.expParseErr {
				return
			}

			// do
			_, actDoErr := cmd.Do(command.Dependencies{
				TaskRepo:    taskRepo,
				LocalIDRepo: localIDRepo,
				SyncRepo:    syncRepo,
			})
			if tc.expDoErr != (actDoErr != nil) {
				t.Errorf("exp %v, got %v", tc.expDoErr, actDoErr)
			}
			if tc.expDoErr {
				return
			}

			// check
			actTask, err := taskRepo.FindOne(tskID)
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if diff := item.TaskDiff(tc.expTask, actTask); diff != "" {
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
