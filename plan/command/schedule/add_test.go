package schedule_test

import (
	"testing"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/command/schedule"
	"go-mod.ewintr.nl/planner/plan/storage/memory"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	aDate := item.NewDate(2024, 11, 2)

	for _, tc := range []struct {
		name        string
		main        []string
		fields      map[string]string
		expErr      bool
		expSchedule item.Schedule
	}{
		{
			name:   "empty",
			expErr: true,
		},
		{
			name: "title missing",
			main: []string{"sched", "add"},
			fields: map[string]string{
				"date": aDate.String(),
			},
			expErr: true,
		},
		{
			name: "all",
			main: []string{"sched", "add", "title"},
			fields: map[string]string{
				"date": aDate.String(),
			},
			expSchedule: item.Schedule{
				ID:   "title",
				Date: aDate,
				ScheduleBody: item.ScheduleBody{
					Title: "title",
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// setup
			mems := memory.New()

			// parse
			cmd, actParseErr := schedule.NewAddArgs().Parse(tc.main, tc.fields)
			if tc.expErr != (actParseErr != nil) {
				t.Errorf("exp %v, got %v", tc.expErr, actParseErr)
			}
			if tc.expErr {
				return
			}

			// do
			if _, err := cmd.Do(mems, nil); err != nil {
				t.Errorf("exp nil, got %v", err)
			}

			// check
			actSchedules, err := mems.Schedule(nil).Find(aDate.Add(-1), aDate.Add(1))
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if len(actSchedules) != 1 {
				t.Errorf("exp 1, got %d", len(actSchedules))
			}

			actLocalIDs, err := mems.LocalID(nil).FindAll()
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if len(actLocalIDs) != 1 {
				t.Errorf("exp 1, got %v", len(actLocalIDs))
			}
			if _, ok := actLocalIDs[actSchedules[0].ID]; !ok {
				t.Errorf("exp true, got %v", ok)
			}

			if actSchedules[0].ID == "" {
				t.Errorf("exp string not te be empty")
			}
			tc.expSchedule.ID = actSchedules[0].ID
			if diff := item.ScheduleDiff(tc.expSchedule, actSchedules[0]); diff != "" {
				t.Errorf("(exp -, got +)\n%s", diff)
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
