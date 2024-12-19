package item

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/go-cmp/cmp"
)

type EventBody struct {
	Title    string        `json:"title"`
	Time     Time          `json:"time"`
	Duration time.Duration `json:"duration"`
}

func (e EventBody) MarshalJSON() ([]byte, error) {
	type Alias EventBody
	return json.Marshal(&struct {
		Duration string `json:"duration"`
		*Alias
	}{
		Duration: e.Duration.String(),
		Alias:    (*Alias)(&e),
	})
}

func (e *EventBody) UnmarshalJSON(data []byte) error {
	type Alias EventBody
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

type Event struct {
	ID        string   `json:"id"`
	Date      Date     `json:"date"`
	Recurrer  Recurrer `json:"recurrer"`
	RecurNext Date     `json:"recurNext"`
	EventBody
}

func NewEvent(i Item) (Event, error) {
	if i.Kind != KindEvent {
		return Event{}, fmt.Errorf("item is not an event")
	}

	var e Event
	if err := json.Unmarshal([]byte(i.Body), &e); err != nil {
		return Event{}, fmt.Errorf("could not unmarshal item body: %v", err)
	}

	e.ID = i.ID
	e.Date = i.Date
	e.Recurrer = i.Recurrer
	e.RecurNext = i.RecurNext

	return e, nil
}

func (e Event) Item() (Item, error) {
	body, err := json.Marshal(e.EventBody)
	if err != nil {
		return Item{}, fmt.Errorf("could not marshal event body to json")
	}

	return Item{
		ID:        e.ID,
		Kind:      KindEvent,
		Date:      e.Date,
		Recurrer:  e.Recurrer,
		RecurNext: e.RecurNext,
		Body:      string(body),
	}, nil
}

func (e Event) Valid() bool {
	if e.Title == "" {
		return false
	}
	if e.Date.IsZero() {
		return false
	}
	if e.Duration.Seconds() < 1 {
		return false
	}

	return true
}

func EventDiff(a, b Event) string {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)

	return cmp.Diff(string(aJSON), string(bJSON))
}

func EventDiffs(a, b []Event) string {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)

	return cmp.Diff(string(aJSON), string(bJSON))
}
