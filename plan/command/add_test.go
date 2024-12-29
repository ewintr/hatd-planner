package command_test

import (
	"testing"
	"time"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/command"
	"go-mod.ewintr.nl/planner/plan/storage/memory"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	aDate := item.NewDate(2024, 11, 2)
	aTime := item.NewTime(12, 0)
	anHourStr := "1h"
	anHour := time.Hour

	for _, tc := range []struct {
		name    string
		main    []string
		fields  map[string]string
		expErr  bool
		expTask item.Task
	}{
		{
			name:   "empty",
			expErr: true,
		},
		{
			name: "title missing",
			main: []string{"add"},
			fields: map[string]string{
				"date": aDate.String(),
			},
			expErr: true,
		},
		{
			name: "date time duration",
			main: []string{"add", "title"},
			fields: map[string]string{
				"date":     aDate.String(),
				"time":     aTime.String(),
				"duration": anHourStr,
			},
			expTask: item.Task{
				ID:   "title",
				Date: aDate,
				TaskBody: item.TaskBody{
					Title:    "title",
					Time:     aTime,
					Duration: anHour,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			taskRepo := memory.NewTask()
			localIDRepo := memory.NewLocalID()
			syncRepo := memory.NewSync()
			cmd, actParseErr := command.NewAddArgs().Parse(tc.main, tc.fields)
			if tc.expErr != (actParseErr != nil) {
				t.Errorf("exp %v, got %v", tc.expErr, actParseErr)
			}
			if tc.expErr {
				return
			}
			if _, err := cmd.Do(command.Dependencies{
				TaskRepo:    taskRepo,
				LocalIDRepo: localIDRepo,
				SyncRepo:    syncRepo,
			}); err != nil {
				t.Errorf("exp nil, got %v", err)
			}

			actTasks, err := taskRepo.FindAll()
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if len(actTasks) != 1 {
				t.Errorf("exp 1, got %d", len(actTasks))
			}

			actLocalIDs, err := localIDRepo.FindAll()
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if len(actLocalIDs) != 1 {
				t.Errorf("exp 1, got %v", len(actLocalIDs))
			}
			if _, ok := actLocalIDs[actTasks[0].ID]; !ok {
				t.Errorf("exp true, got %v", ok)
			}

			if actTasks[0].ID == "" {
				t.Errorf("exp string not te be empty")
			}
			tc.expTask.ID = actTasks[0].ID
			if diff := item.TaskDiff(tc.expTask, actTasks[0]); diff != "" {
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
