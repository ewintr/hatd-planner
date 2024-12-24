package command

import (
	"fmt"

	"go-mod.ewintr.nl/planner/plan/storage"
)

type List struct {
	localIDRepo storage.LocalID
	taskRepo    storage.Task
}

func NewList(localIDRepo storage.LocalID, taskRepo storage.Task) Command {
	return &List{
		localIDRepo: localIDRepo,
		taskRepo:    taskRepo,
	}
}

func (list *List) Execute(main []string, flags map[string]string) error {
	if len(main) > 0 && main[0] != "list" {
		return ErrWrongCommand
	}

	return list.do()
}

func (list *List) do() error {
	localIDs, err := list.localIDRepo.FindAll()
	if err != nil {
		return fmt.Errorf("could not get local ids: %v", err)
	}
	all, err := list.taskRepo.FindAll()
	if err != nil {
		return err
	}
	for _, e := range all {
		lid, ok := localIDs[e.ID]
		if !ok {
			return fmt.Errorf("could not find local id for %s", e.ID)
		}
		fmt.Printf("%s\t%d\t%s\t%s\t%s\n", e.ID, lid, e.Title, e.Date.String(), e.Duration.String())
	}

	return nil
}
