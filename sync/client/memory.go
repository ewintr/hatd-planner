package client

import (
	"slices"
	"sync"
	"time"

	"go-mod.ewintr.nl/planner/item"
)

type Memory struct {
	items map[string]item.Item
	sync.RWMutex
}

func NewMemory() *Memory {
	return &Memory{
		items: make(map[string]item.Item, 0),
	}
}

func (m *Memory) Update(items []item.Item) error {
	m.Lock()
	defer m.Unlock()

	for _, i := range items {
		m.items[i.ID] = i
	}

	return nil
}

func (m *Memory) Updated(kw []item.Kind, ts time.Time) ([]item.Item, error) {
	m.RLock()
	defer m.RUnlock()

	res := make([]item.Item, 0)
	for _, i := range m.items {
		if slices.Contains(kw, i.Kind) && (i.Updated.After(ts) || i.Updated.Equal(ts)) {
			res = append(res, i)
		}
	}

	return res, nil
}
