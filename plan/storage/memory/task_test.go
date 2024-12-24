package memory

import (
	"testing"

	"go-mod.ewintr.nl/planner/item"
)

func TestTask(t *testing.T) {
	t.Parallel()

	mem := NewTask()

	t.Log("empty")
	actTasks, actErr := mem.FindAll()
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if len(actTasks) != 0 {
		t.Errorf("exp 0, got %d", len(actTasks))
	}

	t.Log("store")
	tsk1 := item.Task{
		ID: "id-1",
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
	actTask, actErr := mem.Find(tsk1.ID)
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if actTask.ID != tsk1.ID {
		t.Errorf("exp %v, got %v", tsk1.ID, actTask.ID)
	}

	t.Log("find all")
	actTasks, actErr = mem.FindAll()
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if diff := item.TaskDiffs([]item.Task{tsk1, tsk2}, actTasks); diff != "" {
		t.Errorf("(exp -, got +)\n%s", diff)
	}
}
