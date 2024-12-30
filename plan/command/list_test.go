package command_test

import (
	"testing"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/command"
	"go-mod.ewintr.nl/planner/plan/storage/memory"
)

func TestList(t *testing.T) {
	t.Parallel()

	taskRepo := memory.NewTask()
	localRepo := memory.NewLocalID()
	e := item.Task{
		ID:   "id",
		Date: item.NewDate(2024, 10, 7),
		TaskBody: item.TaskBody{
			Title: "name",
		},
	}
	if err := taskRepo.Store(e); err != nil {
		t.Errorf("exp nil, got %v", err)
	}
	if err := localRepo.Store(e.ID, 1); err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	for _, tc := range []struct {
		name   string
		main   []string
		expRes bool
		expErr bool
	}{
		{
			name:   "empty",
			main:   []string{},
			expRes: true,
		},
		{
			name:   "list",
			main:   []string{"list"},
			expRes: true,
		},
		{
			name: "empty list",
			main: []string{"list", "recur"},
		},
		{
			name:   "wrong",
			main:   []string{"delete"},
			expErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// parse
			cmd, actErr := command.NewListArgs().Parse(tc.main, nil)
			if tc.expErr != (actErr != nil) {
				t.Errorf("exp %v, got %v", tc.expErr, actErr)
			}
			if tc.expErr {
				return
			}

			// do
			res, err := cmd.Do(command.Dependencies{
				TaskRepo:    taskRepo,
				LocalIDRepo: localRepo,
			})
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}

			// check
			listRes := res.(command.ListResult)
			actRes := len(listRes.Tasks) > 0
			if tc.expRes != actRes {
				t.Errorf("exp %v, got %v", tc.expRes, actRes)
			}
		})
	}
}
