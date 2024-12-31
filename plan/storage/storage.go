package storage

import (
	"errors"
	"sort"
	"time"

	"go-mod.ewintr.nl/planner/item"
)

var (
	ErrNotFound = errors.New("not found")
)

type LocalID interface {
	FindOne(lid int) (string, error)
	FindAll() (map[string]int, error)
	FindOrNext(id string) (int, error)
	Next() (int, error)
	Store(id string, localID int) error
	Delete(id string) error
}

type Sync interface {
	FindAll() ([]item.Item, error)
	Store(i item.Item) error
	DeleteAll() error
	LastUpdate() (time.Time, error)
}

type TaskListParams struct {
	Recurrer      bool
	Date          item.Date
	IncludeBefore bool
	Project       string
}

type Task interface {
	Store(task item.Task) error
	FindOne(id string) (item.Task, error)
	FindMany(params TaskListParams) ([]item.Task, error)
	Delete(id string) error
}

func Match(tsk item.Task, params TaskListParams) bool {
	if params.Recurrer && tsk.Recurrer == nil {
		return false
	}
	if !params.Date.IsZero() {
		if !params.IncludeBefore && !params.Date.Equal(tsk.Date) {
			return false
		}
		if params.IncludeBefore && tsk.Date.After(params.Date) {
			return false
		}
	}
	if params.Project != "" && params.Project != tsk.Project {
		return false
	}

	return true
}

func NextLocalID(used []int) int {
	if len(used) == 0 {
		return 1
	}

	sort.Ints(used)
	usedMax := 1
	for _, u := range used {
		if u > usedMax {
			usedMax = u
		}
	}

	var limit int
	for limit = 1; limit <= len(used) || limit < usedMax; limit *= 10 {
	}

	newId := used[len(used)-1] + 1
	if newId < limit {
		return newId
	}

	usedMap := map[int]bool{}
	for _, u := range used {
		usedMap[u] = true
	}

	for i := 1; i < limit; i++ {
		if _, ok := usedMap[i]; !ok {
			return i
		}
	}

	return limit
}
