package memory

import (
	"sort"
	"sync"
	"time"

	"go-mod.ewintr.nl/planner/item"
)

type Sync struct {
	items map[string]item.Item
	mutex sync.RWMutex
}

func NewSync() *Sync {
	return &Sync{
		items: make(map[string]item.Item),
	}
}

func (r *Sync) FindAll() ([]item.Item, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	items := make([]item.Item, 0, len(r.items))
	for _, item := range r.items {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	return items, nil
}

func (r *Sync) Store(e item.Item) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.items[e.ID] = e

	return nil
}

func (r *Sync) DeleteAll() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.items = make(map[string]item.Item)

	return nil
}

func (r *Sync) SetLastUpdate(ts time.Time) error {
	return nil
}

func (r *Sync) LastUpdate() (time.Time, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var last time.Time
	for _, i := range r.items {
		if i.Updated.After(last) {
			last = i.Updated
		}
	}

	return last, nil
}
