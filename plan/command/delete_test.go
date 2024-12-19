package command_test

import (
	"errors"
	"testing"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/command"
	"go-mod.ewintr.nl/planner/plan/storage"
	"go-mod.ewintr.nl/planner/plan/storage/memory"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	e := item.Event{
		ID:   "id",
		Date: item.NewDate(2024, 10, 7),
		EventBody: item.EventBody{
			Title: "name",
		},
	}

	for _, tc := range []struct {
		name   string
		main   []string
		flags  map[string]string
		expErr bool
	}{
		{
			name:   "invalid",
			main:   []string{"update"},
			expErr: true,
		},
		{
			name:   "not found",
			main:   []string{"delete", "5"},
			expErr: true,
		},
		{
			name: "valid",
			main: []string{"delete", "1"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			eventRepo := memory.NewEvent()
			syncRepo := memory.NewSync()
			if err := eventRepo.Store(e); err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			localRepo := memory.NewLocalID()
			if err := localRepo.Store(e.ID, 1); err != nil {
				t.Errorf("exp nil, got %v", err)
			}

			cmd := command.NewDelete(localRepo, eventRepo, syncRepo)

			actErr := cmd.Execute(tc.main, tc.flags) != nil
			if tc.expErr != actErr {
				t.Errorf("exp %v, got %v", tc.expErr, actErr)
			}
			if tc.expErr {
				return
			}

			_, repoErr := eventRepo.Find(e.ID)
			if !errors.Is(repoErr, storage.ErrNotFound) {
				t.Errorf("exp %v, got %v", storage.ErrNotFound, actErr)
			}
			idMap, idErr := localRepo.FindAll()
			if idErr != nil {
				t.Errorf("exp nil, got %v", idErr)
			}
			if len(idMap) != 0 {
				t.Errorf("exp 0, got %v", len(idMap))
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
