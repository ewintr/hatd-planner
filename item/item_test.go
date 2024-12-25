package item_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"go-mod.ewintr.nl/planner/item"
)

func TestItemJSON(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		item    item.Item
		expJSON string
	}{
		{
			name: "minimal",
			item: item.Item{
				ID:   "a",
				Kind: item.KindTask,
				Body: `{"title":"title"}`,
			},
			expJSON: `{
  "recurrer": "",
  "id": "a",
  "kind": "task",
  "updated": "0001-01-01T00:00:00Z",
  "deleted": false,
  "date": "",
  "recurNext": "",
  "body": "{\"title\":\"title\"}"
}`,
		},
		{
			name: "full",
			item: item.Item{
				ID:        "a",
				Kind:      item.KindTask,
				Updated:   time.Date(2024, 12, 25, 11, 9, 0, 0, time.UTC),
				Deleted:   true,
				Date:      item.NewDate(2024, 12, 26),
				Recurrer:  item.NewRecurrer("2024-12-25, daily"),
				RecurNext: item.NewDateFromString("2024-12-30"),
				Body:      `{"title":"title"}`,
			},
			expJSON: `{
  "recurrer": "2024-12-25, daily",
  "id": "a",
  "kind": "task",
  "updated": "2024-12-25T11:09:00Z",
  "deleted": true,
  "date": "2024-12-26",
  "recurNext": "2024-12-30",
  "body": "{\"title\":\"title\"}"
}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actJSON, err := json.Marshal(tc.item)
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			expJSON := bytes.NewBuffer([]byte(``))
			if err := json.Compact(expJSON, []byte(tc.expJSON)); err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if expJSON.String() != string(actJSON) {
				t.Errorf("exp %v, got %v", expJSON.String(), string(actJSON))
			}
			var actItem item.Item
			if err := json.Unmarshal(actJSON, &actItem); err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			if diff := item.ItemDiff(tc.item, actItem); diff != "" {
				t.Errorf("(+exp, -got)%s\n", diff)
			}
		})
	}
}
