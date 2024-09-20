package item

import (
	"encoding/json"
	"fmt"
	"time"
)

type EventBody struct {
	Title string    `json:"title"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func (e EventBody) MarshalJSON() ([]byte, error) {
	type Alias EventBody
	return json.Marshal(&struct {
		Start string `json:"start"`
		End   string `json:"end"`
		*Alias
	}{
		Start: e.Start.UTC().Format(time.RFC3339),
		End:   e.End.UTC().Format(time.RFC3339),
		Alias: (*Alias)(&e),
	})
}

type Event struct {
	ID string `json:"id"`
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

	return e, nil
}

func (e Event) Item() (Item, error) {
	body, err := json.Marshal(EventBody{
		Title: e.Title,
		Start: e.Start,
		End:   e.End,
	})
	if err != nil {
		return Item{}, fmt.Errorf("could not marshal event to json")
	}

	return Item{
		ID:   e.ID,
		Kind: KindEvent,
		Body: string(body),
	}, nil
}
