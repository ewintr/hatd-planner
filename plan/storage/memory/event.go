package memory

import (
	"sort"
	"sync"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage"
)

type Event struct {
	events map[string]item.Event
	mutex  sync.RWMutex
}

func NewEvent() *Event {
	return &Event{
		events: make(map[string]item.Event),
	}
}

func (r *Event) Find(id string) (item.Event, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	event, exists := r.events[id]
	if !exists {
		return item.Event{}, storage.ErrNotFound
	}
	return event, nil
}

func (r *Event) FindAll() ([]item.Event, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	events := make([]item.Event, 0, len(r.events))
	for _, event := range r.events {
		events = append(events, event)
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].ID < events[j].ID
	})

	return events, nil
}

func (r *Event) Store(e item.Event) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.events[e.ID] = e

	return nil
}

func (r *Event) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.events[id]; !exists {
		return storage.ErrNotFound
	}
	delete(r.events, id)

	return nil
}
