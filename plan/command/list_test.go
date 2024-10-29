package command_test

import (
	"testing"
	"time"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/command"
	"go-mod.ewintr.nl/planner/plan/storage/memory"
)

func TestList(t *testing.T) {
	t.Parallel()

	eventRepo := memory.NewEvent()
	localRepo := memory.NewLocalID()
	e := item.Event{
		ID: "id",
		EventBody: item.EventBody{
			Title: "name",
			Start: time.Date(2024, 10, 7, 9, 30, 0, 0, time.UTC),
		},
	}
	if err := eventRepo.Store(e); err != nil {
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
			cmd := command.NewList(localRepo, eventRepo)
			actErr := cmd.Execute(tc.main, nil) != nil
			if tc.expErr != actErr {
				t.Errorf("exp %v, got %v", tc.expErr, actErr)
			}
		})
	}
}
