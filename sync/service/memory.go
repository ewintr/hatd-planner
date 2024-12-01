package main

import (
	"fmt"
	"slices"
	"sync"
	"time"

	"go-mod.ewintr.nl/planner/item"
)

type Memory struct {
	items map[string]item.Item
	mutex sync.RWMutex
}

func NewMemory() *Memory {
	return &Memory{
		items: make(map[string]item.Item),
	}
}

func (m *Memory) Update(item item.Item, ts time.Time) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	item.Updated = ts
	m.items[item.ID] = item

	return nil
}

func (m *Memory) Updated(kinds []item.Kind, timestamp time.Time) ([]item.Item, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make([]item.Item, 0)

	for _, i := range m.items {
		timeOK := timestamp.IsZero() || i.Updated.Equal(timestamp) || i.Updated.After(timestamp)
		kindOK := len(kinds) == 0 || slices.Contains(kinds, i.Kind)
		if timeOK && kindOK {
			result = append(result, i)
		}
	}

	return result, nil
}

func (m *Memory) RecursBefore(date time.Time) ([]item.Item, error) {
	res := make([]item.Item, 0)
	for _, i := range m.items {
		if i.Recurrer == nil {
			continue
		}
		if i.RecurNext.Before(date) {
			res = append(res, i)
		}
	}
	return res, nil
}

func (m *Memory) RecursNext(id string, date time.Time, ts time.Time) error {
	i, ok := m.items[id]
	if !ok {
		return ErrNotFound
	}
	if i.Recurrer == nil {
		return ErrNotARecurrer
	}
	if !i.Recurrer.On(date) {
		return fmt.Errorf("item does not recur on %v", date)
	}
	i.RecurNext = date
	i.Updated = ts
	m.items[id] = i

	return nil
}
