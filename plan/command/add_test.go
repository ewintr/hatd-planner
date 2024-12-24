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
	aDay := time.Duration(24) * time.Hour
	anHourStr := "1h"
	anHour := time.Hour

	for _, tc := range []struct {
		name    string
		main    []string
		flags   map[string]string
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
			flags: map[string]string{
				command.FlagOn: aDate.String(),
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
				command.FlagOn: aDate.String(),
			},
			expTask: item.Task{
				ID:   "title",
				Date: aDate,
				TaskBody: item.TaskBody{
					Title:    "title",
					Duration: aDay,
				},
			},
		},
		{
			name: "date, time and duration",
			main: []string{"add", "title"},
			flags: map[string]string{
				command.FlagOn:  aDate.String(),
				command.FlagAt:  aTime.String(),
				command.FlagFor: anHourStr,
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
		{
			name: "date and duration",
			main: []string{"add", "title"},
			flags: map[string]string{
				command.FlagOn:  aDate.String(),
				command.FlagFor: anHourStr,
			},
			expErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			taskRepo := memory.NewTask()
			localRepo := memory.NewLocalID()
			syncRepo := memory.NewSync()
			cmd := command.NewAdd(localRepo, taskRepo, syncRepo)
			actParseErr := cmd.Execute(tc.main, tc.flags) != nil
			if tc.expErr != actParseErr {
				t.Errorf("exp %v, got %v", tc.expErr, actParseErr)
			}
			if tc.expErr {
				return
			}

			actTasks, err := taskRepo.FindAll()
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if len(actTasks) != 1 {
				t.Errorf("exp 1, got %d", len(actTasks))
			}

			actLocalIDs, err := localRepo.FindAll()
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
