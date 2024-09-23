package main

import (
	"errors"
	"sort"
	"sync"

	"go-mod.ewintr.nl/planner/item"
)

type Memory struct {
	events map[string]item.Event
	mutex  sync.RWMutex
}

func NewMemory() *Memory {
	return &Memory{
		events: make(map[string]item.Event),
	}
}

func (r *Memory) Find(id string) (item.Event, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	event, exists := r.events[id]
	if !exists {
		return item.Event{}, errors.New("event not found")
	}
	return event, nil
}

func (r *Memory) FindAll() ([]item.Event, error) {
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

func (r *Memory) Store(e item.Event) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.events[e.ID] = e
	return nil
}

func (r *Memory) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.events[id]; !exists {
		return errors.New("event not found")
	}
	delete(r.events, id)
	return nil
}
