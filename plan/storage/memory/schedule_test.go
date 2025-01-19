package memory_test

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage/memory"
)

func TestSchedule(t *testing.T) {
	t.Parallel()

	mem := memory.NewSchedule()

	actScheds, actErr := mem.Find(item.NewDateFromString("1900-01-01"), item.NewDateFromString("9999-12-31"))
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if len(actScheds) != 0 {
		t.Errorf("exp 0, got %d", len(actScheds))
	}

	s1 := item.Schedule{
		ID:   "id-1",
		Date: item.NewDateFromString("2025-01-20"),
	}
	if err := mem.Store(s1); err != nil {
		t.Errorf("exp nil, got %v", err)
	}
	s2 := item.Schedule{
		ID:   "id-2",
		Date: item.NewDateFromString("2025-01-21"),
	}
	if err := mem.Store(s2); err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	for _, tc := range []struct {
		name  string
		start string
		end   string
		exp   []string
	}{
		{
			name:  "all",
			start: "1900-01-01",
			end:   "9999-12-31",
			exp:   []string{s1.ID, s2.ID},
		},
		{
			name:  "last",
			start: s2.Date.String(),
			end:   "9999-12-31",
			exp:   []string{s2.ID},
		},
		{
			name:  "first",
			start: "1900-01-01",
			end:   s1.Date.String(),
			exp:   []string{s1.ID},
		},
		{
			name:  "none",
			start: "1900-01-01",
			end:   "2025-01-01",
			exp:   make([]string, 0),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actScheds, actErr = mem.Find(item.NewDateFromString(tc.start), item.NewDateFromString(tc.end))
			if actErr != nil {
				t.Errorf("exp nil, got %v", actErr)
			}
			actIDs := make([]string, 0, len(actScheds))
			for _, s := range actScheds {
				actIDs = append(actIDs, s.ID)
			}
			sort.Strings(actIDs)
			if diff := cmp.Diff(tc.exp, actIDs); diff != "" {
				t.Errorf("(+exp, -got)%s\n", diff)
			}
		})
	}

}
