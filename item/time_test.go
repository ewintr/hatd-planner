package item_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"go-mod.ewintr.nl/planner/item"
)

func TestTime(t *testing.T) {
	t.Parallel()

	h, m := 11, 18
	tm := item.NewTime(h, m)
	expStr := "11:18"
	if expStr != tm.String() {
		t.Errorf("exp %v, got %v", expStr, tm.String())
	}
	actJSON, err := json.Marshal(tm)
	if err != nil {
		t.Errorf("exp nil, got %v", err)
	}
	expJSON := fmt.Sprintf("%q", expStr)
	if expJSON != string(actJSON) {
		t.Errorf("exp %v, got %v", expJSON, string(actJSON))
	}
	var actTM item.Time
	if err := json.Unmarshal(actJSON, &actTM); err != nil {
		t.Errorf("exp nil, got %v", err)
	}
	if expStr != actTM.String() {
		t.Errorf("ecp %v, got %v", expStr, actTM.String())
	}
}

func TestTimeFromString(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		str  string
		exp  string
	}{
		{
			name: "empty",
			exp:  "00:00",
		},
		{
			name: "invalid",
			str:  "invalid",
			exp:  "00:00",
		},
		{
			name: "valid",
			str:  "11:42",
			exp:  "11:42",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			act := item.NewTimeFromString(tc.str)
			if tc.exp != act.String() {
				t.Errorf("exp %v, got %v", tc.exp, act.String())
			}
		})
	}
}
