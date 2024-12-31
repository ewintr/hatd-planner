package item_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go-mod.ewintr.nl/planner/item"
)

func TestNewTask(t *testing.T) {
	t.Parallel()

	oneHour, err := time.ParseDuration("1h")
	if err != nil {
		t.Errorf("exp nil, got %v", err)
	}
	for _, tc := range []struct {
		name    string
		it      item.Item
		expTask item.Task
		expErr  bool
	}{
		{
			name: "wrong kind",
			it: item.Item{
				ID:   "a",
				Date: item.NewDate(2024, 9, 20),
				Kind: item.KindSchedule,
				Body: `{
  "title":"title",
  "time":"08:00",
  "duration":"1h"
}`,
			},
			expErr: true,
		},
		{
			name: "invalid json",
			it: item.Item{
				ID:   "a",
				Kind: item.KindTask,
				Body: `{"id":"a"`,
			},
			expErr: true,
		},
		{
			name: "valid",
			it: item.Item{
				ID:       "a",
				Kind:     item.KindTask,
				Date:     item.NewDate(2024, 9, 20),
				Recurrer: item.NewRecurrer("2024-12-08, daily"),
				Body: `{
  "title":"title",
  "project":"project",
  "time":"08:00",
  "duration":"1h"
}`,
			},
			expTask: item.Task{
				ID:       "a",
				Date:     item.NewDate(2024, 9, 20),
				Recurrer: item.NewRecurrer("2024-12-08, daily"),
				TaskBody: item.TaskBody{
					Title:    "title",
					Project:  "project",
					Time:     item.NewTime(8, 0),
					Duration: oneHour,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actTask, actErr := item.NewTask(tc.it)
			if tc.expErr != (actErr != nil) {
				t.Errorf("exp nil, got %v", actErr)
			}
			if tc.expErr {
				return
			}
			if diff := item.TaskDiff(tc.expTask, actTask); diff != "" {
				t.Errorf("(+exp, -got)\n%s", diff)
			}
		})
	}
}

func TestTaskItem(t *testing.T) {
	t.Parallel()

	oneHour, err := time.ParseDuration("1h")
	if err != nil {
		t.Errorf("exp nil, got %v", err)
	}
	for _, tc := range []struct {
		name    string
		tsk     item.Task
		expItem item.Item
		expErr  bool
	}{
		{
			name: "empty",
			expItem: item.Item{
				Kind:    item.KindTask,
				Updated: time.Time{},
				Body:    `{"duration":"0s","title":"","project":"","time":""}`,
			},
		},
		{
			name: "normal",
			tsk: item.Task{
				ID:   "a",
				Date: item.NewDate(2024, 9, 23),
				TaskBody: item.TaskBody{
					Title:    "title",
					Project:  "project",
					Time:     item.NewTime(8, 0),
					Duration: oneHour,
				},
			},
			expItem: item.Item{
				ID:      "a",
				Kind:    item.KindTask,
				Updated: time.Time{},
				Date:    item.NewDate(2024, 9, 23),
				Body:    `{"duration":"1h0m0s","title":"title","project":"project","time":"08:00"}`,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actItem, actErr := tc.tsk.Item()
			if tc.expErr != (actErr != nil) {
				t.Errorf("exp nil, got %v", actErr)
			}
			if tc.expErr {
				return
			}
			if diff := cmp.Diff(tc.expItem, actItem); diff != "" {
				t.Errorf("(exp+, got -)\n%s", diff)
			}
		})
	}
}

func TestTaskValidate(t *testing.T) {
	t.Parallel()

	oneHour, err := time.ParseDuration("1h")
	if err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	for _, tc := range []struct {
		name string
		tsk  item.Task
		exp  bool
	}{
		{
			name: "empty",
		},
		{
			name: "missing title",
			tsk: item.Task{
				ID:   "a",
				Date: item.NewDate(2024, 9, 20),
				TaskBody: item.TaskBody{
					Time:     item.NewTime(8, 0),
					Duration: oneHour,
				},
			},
		},
		{
			name: "valid",
			tsk: item.Task{
				ID:   "a",
				Date: item.NewDate(2024, 9, 20),
				TaskBody: item.TaskBody{
					Title:    "title",
					Time:     item.NewTime(8, 0),
					Duration: oneHour,
				},
			},
			exp: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if act := tc.tsk.Valid(); tc.exp != act {
				t.Errorf("exp %v, got %v", tc.exp, act)
			}

		})
	}
}
