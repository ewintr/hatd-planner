package task_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/command/task"
	"go-mod.ewintr.nl/planner/plan/storage/memory"
)

func TestListParse(t *testing.T) {
	t.Parallel()
	now := time.Now()
	today := item.NewDate(now.Year(), int(now.Month()), now.Day())

	for _, tc := range []struct {
		name    string
		main    []string
		fields  map[string]string
		expArgs task.ListArgs
		expErr  bool
	}{
		{
			name:    "empty",
			main:    []string{},
			fields:  map[string]string{},
			expArgs: task.ListArgs{},
		},
		{
			name:   "today",
			main:   []string{"tod"},
			fields: map[string]string{},
			expArgs: task.ListArgs{
				To: today,
			},
		},
		{
			name:   "tomorrow",
			main:   []string{"tom"},
			fields: map[string]string{},
			expArgs: task.ListArgs{
				From: today.Add(1),
				To:   today.Add(1),
			},
		},
		{
			name:   "week",
			main:   []string{"week"},
			fields: map[string]string{},
			expArgs: task.ListArgs{
				From: today,
				To:   today.Add(7),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			nla := task.NewListArgs()
			cmd, actErr := nla.Parse(tc.main, tc.fields)
			if tc.expErr != (actErr != nil) {
				t.Errorf("exp %v, got %v", tc.expErr, actErr != nil)
			}
			if tc.expErr {
				return
			}
			listCmd, ok := cmd.(task.List)
			if !ok {
				t.Errorf("exp true, got false")
			}
			if diff := cmp.Diff(tc.expArgs, listCmd.Args, cmpopts.IgnoreTypes(map[string][]string{})); diff != "" {
				t.Errorf("(+exp, -got)\n%s\n", diff)
			}
		})
	}
}

func TestList(t *testing.T) {
	t.Parallel()

	mems := memory.New()

	e := item.Task{
		ID:   "id",
		Date: item.NewDate(2024, 10, 7),
		TaskBody: item.TaskBody{
			Title: "name",
		},
	}
	if err := mems.Task(nil).Store(e); err != nil {
		t.Errorf("exp nil, got %v", err)
	}
	if err := mems.LocalID(nil).Store(e.ID, 1); err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	for _, tc := range []struct {
		name   string
		cmd    task.List
		expRes bool
		expErr bool
	}{
		{
			name:   "empty",
			expRes: true,
		},
		{
			name: "empty list",
			cmd: task.List{
				Args: task.ListArgs{
					HasRecurrer: true,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.cmd.Do(mems, nil)
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}

			listRes := res.(task.ListResult)
			actRes := len(listRes.Tasks) > 0
			if tc.expRes != actRes {
				t.Errorf("exp %v, got %v", tc.expRes, actRes)
			}
		})
	}
}
