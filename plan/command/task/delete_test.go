package task_test

import (
	"errors"
	"testing"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/command/task"
	"go-mod.ewintr.nl/planner/plan/storage"
	"go-mod.ewintr.nl/planner/plan/storage/memory"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	e := item.Task{
		ID:   "id",
		Date: item.NewDate(2024, 10, 7),
		TaskBody: item.TaskBody{
			Title: "name",
		},
	}

	for _, tc := range []struct {
		name        string
		main        []string
		flags       map[string]string
		expParseErr bool
		expDoErr    bool
	}{
		{
			name:        "invalid",
			main:        []string{"update"},
			expParseErr: true,
		},
		{
			name:     "not found",
			main:     []string{"delete", "5"},
			expDoErr: true,
		},
		{
			name: "valid",
			main: []string{"delete", "1"},
		},
		{
			name: "done",
			main: []string{"done", "1"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// setup
			mems := memory.New()

			if err := mems.Task(nil).Store(e); err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if err := mems.LocalID(nil).Store(e.ID, 1); err != nil {
				t.Errorf("exp nil, got %v", err)
			}

			// parse
			cmd, actParseErr := task.NewDeleteArgs().Parse(tc.main, tc.flags)
			if tc.expParseErr != (actParseErr != nil) {
				t.Errorf("exp %v, got %v", tc.expParseErr, actParseErr)
			}
			if tc.expParseErr {
				return
			}

			// do
			_, actDoErr := cmd.Do(mems, nil)
			if tc.expDoErr != (actDoErr != nil) {
				t.Errorf("exp false, got %v", actDoErr)
			}
			if tc.expDoErr {
				return
			}

			// check
			_, repoErr := mems.Task(nil).FindOne(e.ID)
			if !errors.Is(repoErr, storage.ErrNotFound) {
				t.Errorf("exp %v, got %v", storage.ErrNotFound, repoErr)
			}
			idMap, idErr := mems.LocalID(nil).FindAll()
			if idErr != nil {
				t.Errorf("exp nil, got %v", idErr)
			}
			if len(idMap) != 0 {
				t.Errorf("exp 0, got %v", len(idMap))
			}
			updated, err := mems.Sync(nil).FindAll()
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if len(updated) != 1 {
				t.Errorf("exp 1, got %v", len(updated))
			}
		})
	}
}
