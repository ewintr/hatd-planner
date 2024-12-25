package item

import (
	"encoding/json"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

type Kind string

const (
	KindSchedule Kind = "schedule"
	KindTask     Kind = "task"
)

var (
	KnownKinds = []Kind{KindSchedule, KindTask}
)

type Item struct {
	ID        string    `json:"id"`
	Kind      Kind      `json:"kind"`
	Updated   time.Time `json:"updated"`
	Deleted   bool      `json:"deleted"`
	Date      Date      `json:"date"`
	Recurrer  Recurrer  `json:"recurrer"`
	RecurNext Date      `json:"recurNext"`
	Body      string    `json:"body"`
}

func (i Item) MarshalJSON() ([]byte, error) {
	var recurStr string
	if i.Recurrer != nil {
		recurStr = i.Recurrer.String()
	}
	type Alias Item
	return json.Marshal(&struct {
		Recurrer string `json:"recurrer"`
		*Alias
	}{
		Recurrer: recurStr,
		Alias:    (*Alias)(&i),
	})
}

func (i *Item) UnmarshalJSON(data []byte) error {
	type Alias Item
	aux := &struct {
		Recurrer string `json:"recurrer"`
		*Alias
	}{
		Alias: (*Alias)(i),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	i.Recurrer = NewRecurrer(aux.Recurrer)

	return nil
}

func NewItem(k Kind, body string) Item {
	return Item{
		ID:      uuid.New().String(),
		Kind:    k,
		Updated: time.Now(),
		Body:    body,
	}
}

func ItemDiff(exp, got Item) string {
	expJSON, _ := json.Marshal(exp)
	actJSON, _ := json.Marshal(got)

	return cmp.Diff(string(expJSON), string(actJSON))
}
