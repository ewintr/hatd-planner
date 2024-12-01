package item

import (
	"encoding/json"
	"fmt"
	"time"
)

type EventBody struct {
	Title    string        `json:"title"`
	Start    time.Time     `json:"start"`
	Duration time.Duration `json:"duration"`
}

func (e EventBody) MarshalJSON() ([]byte, error) {
	type Alias EventBody
	return json.Marshal(&struct {
		Start    string `json:"start"`
		Duration string `json:"duration"`
		*Alias
	}{
		Start:    e.Start.UTC().Format(time.RFC3339),
		Duration: e.Duration.String(),
		Alias:    (*Alias)(&e),
	})
}

func (e *EventBody) UnmarshalJSON(data []byte) error {
	type Alias EventBody
	aux := &struct {
		Start    string `json:"start"`
		Duration string `json:"duration"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var err error
	if e.Start, err = time.Parse(time.RFC3339, aux.Start); err != nil {
		return err
	}

	if e.Duration, err = time.ParseDuration(aux.Duration); err != nil {
		return err
	}

	return nil
}

type Event struct {
	ID        string    `json:"id"`
	Recurrer  *Recur    `json:"recurrer"`
	RecurNext time.Time `json:"recurNext"`
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
	e.Recurrer = i.Recurrer
	e.RecurNext = i.RecurNext

	return e, nil
}

func (e Event) Item() (Item, error) {
	body, err := json.Marshal(EventBody{
		Title:    e.Title,
		Start:    e.Start,
		Duration: e.Duration,
	})
	if err != nil {
		return Item{}, fmt.Errorf("could not marshal event to json")
	}

	return Item{
		ID:        e.ID,
		Kind:      KindEvent,
		Recurrer:  e.Recurrer,
		RecurNext: e.RecurNext,
		Body:      string(body),
	}, nil
}

func (e Event) Valid() bool {
	if e.Title == "" {
		return false
	}
	if e.Start.IsZero() || e.Start.Year() < 2024 {
		return false
	}
	if e.Duration.Seconds() < 1 {
		return false
	}
	if e.Recurrer != nil && !e.Recurrer.Valid() {
		return false
	}

	return true
}
