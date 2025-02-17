package memory

import (
	"sort"
	"sync"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage"
)

type Task struct {
	tasks map[string]item.Task
	mutex sync.RWMutex
}

func NewTask() *Task {
	return &Task{
		tasks: make(map[string]item.Task),
	}
}

func (t *Task) FindOne(id string) (item.Task, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	task, exists := t.tasks[id]
	if !exists {
		return item.Task{}, storage.ErrNotFound
	}
	return task, nil
}

func (t *Task) FindMany(params storage.TaskListParams) ([]item.Task, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	tasks := make([]item.Task, 0, len(t.tasks))
	for _, tsk := range t.tasks {
		if storage.MatchTask(tsk, params) {
			tasks = append(tasks, tsk)
		}
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})

	return tasks, nil
}

func (t *Task) Store(tsk item.Task) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.tasks[tsk.ID] = tsk

	return nil
}

func (t *Task) Delete(id string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if _, exists := t.tasks[id]; !exists {
		return storage.ErrNotFound
	}
	delete(t.tasks, id)

	return nil
}

func (t *Task) Projects() (map[string]int, error) {
	projects := make(map[string]int)
	for _, tsk := range t.tasks {
		if _, ok := projects[tsk.Project]; !ok {
			projects[tsk.Project] = 0
		}
		projects[tsk.Project]++
	}

	return projects, nil
}
