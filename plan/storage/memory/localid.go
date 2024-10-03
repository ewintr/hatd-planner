package memory

import (
	"sync"

	"go-mod.ewintr.nl/planner/plan/storage"
)

type LocalID struct {
	ids   map[string]int
	mutex sync.RWMutex
}

func NewLocalID() *LocalID {
	return &LocalID{
		ids: make(map[string]int),
	}
}

func (ml *LocalID) FindAll() (map[string]int, error) {
	ml.mutex.RLock()
	defer ml.mutex.RUnlock()

	return ml.ids, nil
}

func (ml *LocalID) Next() (int, error) {
	ml.mutex.RLock()
	defer ml.mutex.RUnlock()

	cur := make([]int, 0, len(ml.ids))
	for _, i := range ml.ids {
		cur = append(cur, i)
	}

	localID := storage.NextLocalID(cur)

	return localID, nil
}

func (ml *LocalID) Store(id string, localID int) error {
	ml.mutex.Lock()
	defer ml.mutex.Unlock()

	ml.ids[id] = localID

	return nil
}

func (ml *LocalID) Delete(id string) error {
	ml.mutex.Lock()
	defer ml.mutex.Unlock()

	if _, ok := ml.ids[id]; !ok {
		return storage.ErrNotFound
	}

	delete(ml.ids, id)

	return nil
}
