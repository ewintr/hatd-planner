package main

import (
	"slices"
	"time"
)

type Memory struct {
	items map[string]Item
}

func NewMemory() *Memory {
	return &Memory{
		items: make(map[string]Item),
	}
}

func (m *Memory) Update(item Item) error {
	m.items[item.ID] = item

	return nil
}

func (m *Memory) Updated(kinds []Kind, timestamp time.Time) ([]Item, error) {
	result := make([]Item, 0)

	for _, i := range m.items {
		timeOK := timestamp.IsZero() || i.Updated.Equal(timestamp) || i.Updated.After(timestamp)
		kindOK := len(kinds) == 0 || slices.Contains(kinds, i.Kind)
		if timeOK && kindOK {
			result = append(result, i)
		}
	}

	return result, nil
}
