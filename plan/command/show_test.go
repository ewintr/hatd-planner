package command_test

import (
	"fmt"
	"testing"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/command"
	"go-mod.ewintr.nl/planner/plan/storage/memory"
)

func TestShow(t *testing.T) {
	t.Parallel()

	mems := memory.New()

	tsk := item.Task{
		ID:   "id",
		Date: item.NewDate(2024, 10, 7),
		TaskBody: item.TaskBody{
			Title: "name",
		},
	}
	if err := mems.Task(nil).Store(tsk); err != nil {
		t.Errorf("exp nil, got %v", err)
	}
	if err := mems.LocalID(nil).Store(tsk.ID, 1); err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	for _, tc := range []struct {
		name        string
		main        []string
		expData     [][]string
		expParseErr bool
		expDoErr    bool
	}{
		{
			name:        "empty",
			main:        []string{},
			expParseErr: true,
		},
		{
			name:        "wrong",
			main:        []string{"delete"},
			expParseErr: true,
		},
		{
			name: "local id",
			main: []string{"1"},
			expData: [][]string{
				{"title", tsk.Title},
				{"local id", fmt.Sprintf("%d", 1)},
				{"date", tsk.Date.String()},
				{"time", tsk.Time.String()},
				{"duration", tsk.Duration.String()},
				{"recur", ""},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// parse
			cmd, actParseErr := command.NewShowArgs().Parse(tc.main, nil)
			if tc.expParseErr != (actParseErr != nil) {
				t.Errorf("exp %v, got %v", tc.expParseErr, actParseErr != nil)
			}
			if tc.expParseErr {
				return
			}

			// do
			_, actDoErr := cmd.Do(mems, nil)
			if tc.expDoErr != (actDoErr != nil) {
				t.Errorf("exp %v, got %v", tc.expDoErr, actDoErr != nil)
			}
			if tc.expDoErr {
				return
			}
			// if diff := cmp.Diff(tc.expData, actData); diff != "" {
			// 	t.Errorf("(+exp, -got)%s\n", diff)
			// }
		})
	}

}
