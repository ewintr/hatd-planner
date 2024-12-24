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
		expErr bool
	}{
		{
			name: "empty",
			main: []string{},
		},
		{
			name: "list",
			main: []string{"list"},
		},
		{
			name:   "wrong",
			main:   []string{"delete"},
			expErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd := command.NewList(localRepo, taskRepo)
			actErr := cmd.Execute(tc.main, nil) != nil
			if tc.expErr != actErr {
				t.Errorf("exp %v, got %v", tc.expErr, actErr)
			}
		})
	}
}
