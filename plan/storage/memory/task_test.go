package memory

import (
	"testing"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage"
)

func TestTask(t *testing.T) {
	t.Parallel()

	mem := NewTask()

	t.Log("empty")
	actTasks, actErr := mem.FindMany(storage.TaskListParams{})
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if len(actTasks) != 0 {
		t.Errorf("exp 0, got %d", len(actTasks))
	}

	t.Log("store")
	tsk1 := item.Task{
		ID:   "id-1",
		Date: item.NewDate(2024, 12, 29),
	}
	if err := mem.Store(tsk1); err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	tsk2 := item.Task{
		ID: "id-2",
	}
	if err := mem.Store(tsk2); err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	t.Log("find one")
	actTask, actErr := mem.FindOne(tsk1.ID)
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if actTask.ID != tsk1.ID {
		t.Errorf("exp %v, got %v", tsk1.ID, actTask.ID)
	}

	t.Log("find all")
	actTasks, actErr = mem.FindMany(storage.TaskListParams{})
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if diff := item.TaskDiffs([]item.Task{tsk1, tsk2}, actTasks); diff != "" {
		t.Errorf("(exp -, got +)\n%s", diff)
	}

	t.Log("find some")
	actTasks, actErr = mem.FindMany(storage.TaskListParams{
		From: item.NewDate(2024, 12, 29),
	})
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if diff := item.TaskDiffs([]item.Task{tsk1}, actTasks); diff != "" {
		t.Errorf("(exp -, got +)\n%s", diff)
	}

}
