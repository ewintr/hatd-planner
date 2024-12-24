package command

import (
	"fmt"
	"strconv"

	"go-mod.ewintr.nl/planner/plan/storage"
)

type Delete struct {
	localIDRepo storage.LocalID
	taskRepo    storage.Task
	syncRepo    storage.Sync
	localID     int
}

func NewDelete(localIDRepo storage.LocalID, taskRepo storage.Task, syncRepo storage.Sync) Command {
	return &Delete{
		localIDRepo: localIDRepo,
		taskRepo:    taskRepo,
		syncRepo:    syncRepo,
	}
}

func (del *Delete) Execute(main []string, flags map[string]string) error {
	if len(main) < 2 || main[0] != "delete" {
		return ErrWrongCommand
	}
	localID, err := strconv.Atoi(main[1])
	if err != nil {
		return fmt.Errorf("not a local id: %v", main[1])
	}
	del.localID = localID

	return del.do()
}

func (del *Delete) do() error {
	var id string
	idMap, err := del.localIDRepo.FindAll()
	if err != nil {
		return fmt.Errorf("could not get local ids: %v", err)
	}
	for tskID, lid := range idMap {
		if del.localID == lid {
			id = tskID
		}
	}
	if id == "" {
		return fmt.Errorf("could not find local id")
	}

	tsk, err := del.taskRepo.Find(id)
	if err != nil {
		return fmt.Errorf("could not get task: %v", err)
	}

	it, err := tsk.Item()
	if err != nil {
		return fmt.Errorf("could not convert task to sync item: %v", err)
	}
	it.Deleted = true
	if err := del.syncRepo.Store(it); err != nil {
		return fmt.Errorf("could not store sync item: %v", err)
	}

	if err := del.localIDRepo.Delete(id); err != nil {
		return fmt.Errorf("could not delete local id: %v", err)
	}

	if err := del.taskRepo.Delete(id); err != nil {
		return fmt.Errorf("could not delete task: %v", err)
	}

	return nil
}
