package main

import (
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

func (m *Memory) Update(item item.Item) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

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
