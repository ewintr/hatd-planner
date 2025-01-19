package memory

import (
	"sync"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage"
)

type Schedule struct {
	scheds map[string]item.Schedule
	mutex  sync.RWMutex
}

func NewSchedule() *Schedule {
	return &Schedule{
		scheds: make(map[string]item.Schedule),
	}
}

func (s *Schedule) Store(sched item.Schedule) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.scheds[sched.ID] = sched
	return nil
}

func (s *Schedule) Find(start, end item.Date) ([]item.Schedule, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	res := make([]item.Schedule, 0)
	for _, sched := range s.scheds {
		if start.After(sched.Date) || sched.Date.After(end) {
			continue
		}
		res = append(res, sched)
	}

	return res, nil
}

func (s *Schedule) Delete(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.scheds[id]; !exists {
		return storage.ErrNotFound
	}
	delete(s.scheds, id)

	return nil
}
