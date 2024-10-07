package command_test

import (
	"errors"
	"testing"
	"time"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/command"
	"go-mod.ewintr.nl/planner/plan/storage"
	"go-mod.ewintr.nl/planner/plan/storage/memory"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	e := item.Event{
		ID: "id",
		EventBody: item.EventBody{
			Title: "name",
			Start: time.Date(2024, 10, 7, 9, 30, 0, 0, time.UTC),
		},
	}

	for _, tc := range []struct {
		name    string
		localID int
		expErr  bool
	}{
		{
			name:    "not found",
			localID: 5,
			expErr:  true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			eventRepo := memory.NewEvent()
			if err := eventRepo.Store(e); err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			localRepo := memory.NewLocalID()
			if err := localRepo.Store(e.ID, 1); err != nil {
				t.Errorf("exp nil, got %v", err)
			}

			actErr := command.Delete(localRepo, eventRepo, tc.localID) != nil
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
		})
	}
}
