package item

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/go-cmp/cmp"
)

type TaskBody struct {
	Title    string        `json:"title"`
	Time     Time          `json:"time"`
	Duration time.Duration `json:"duration"`
}

func (e TaskBody) MarshalJSON() ([]byte, error) {
	type Alias TaskBody
	return json.Marshal(&struct {
		Duration string `json:"duration"`
		*Alias
	}{
		Duration: e.Duration.String(),
		Alias:    (*Alias)(&e),
	})
}

func (e *TaskBody) UnmarshalJSON(data []byte) error {
	type Alias TaskBody
	aux := &struct {
		Duration string `json:"duration"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var err error
	if e.Duration, err = time.ParseDuration(aux.Duration); err != nil {
		return err
	}

	return nil
}

type Task struct {
	ID        string   `json:"id"`
	Date      Date     `json:"date"`
	Recurrer  Recurrer `json:"recurrer"`
	RecurNext Date     `json:"recurNext"`
	TaskBody
}

func NewTask(i Item) (Task, error) {
	if i.Kind != KindTask {
		return Task{}, fmt.Errorf("item is not an task")
	}

	var t Task
	if err := json.Unmarshal([]byte(i.Body), &t); err != nil {
		return Task{}, fmt.Errorf("could not unmarshal item body: %v", err)
	}

	t.ID = i.ID
	t.Date = i.Date
	t.Recurrer = i.Recurrer
	t.RecurNext = i.RecurNext

	return t, nil
}

func (t Task) Item() (Item, error) {
	body, err := json.Marshal(t.TaskBody)
	if err != nil {
		return Item{}, fmt.Errorf("could not marshal task body to json")
	}

	return Item{
		ID:        t.ID,
		Kind:      KindTask,
		Date:      t.Date,
		Recurrer:  t.Recurrer,
		RecurNext: t.RecurNext,
		Body:      string(body),
	}, nil
}

func (t Task) Valid() bool {
	if t.Title == "" {
		return false
	}

	return true
}

func TaskDiff(a, b Task) string {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)

	return cmp.Diff(string(aJSON), string(bJSON))
}

func TaskDiffs(a, b []Task) string {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)

	return cmp.Diff(string(aJSON), string(bJSON))
}
