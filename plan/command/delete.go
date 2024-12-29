package command

import (
	"fmt"
	"strconv"
)

type DeleteArgs struct {
	LocalID int
}

func NewDeleteArgs() DeleteArgs {
	return DeleteArgs{}
}

func (da DeleteArgs) Parse(main []string, flags map[string]string) (Command, error) {
	if len(main) < 2 || main[0] != "delete" {
		return nil, ErrWrongCommand
	}

	localID, err := strconv.Atoi(main[1])
	if err != nil {
		return nil, fmt.Errorf("not a local id: %v", main[1])
	}

	return &Delete{
		args: DeleteArgs{
			LocalID: localID,
		},
	}, nil
}

type Delete struct {
	args DeleteArgs
}

func (del *Delete) Do(deps Dependencies) ([][]string, error) {
	var id string
	idMap, err := deps.LocalIDRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("could not get local ids: %v", err)
	}
	for tskID, lid := range idMap {
		if del.args.LocalID == lid {
			id = tskID
		}
	}
	if id == "" {
		return nil, fmt.Errorf("could not find local id")
	}

	tsk, err := deps.TaskRepo.Find(id)
	if err != nil {
		return nil, fmt.Errorf("could not get task: %v", err)
	}

	it, err := tsk.Item()
	if err != nil {
		return nil, fmt.Errorf("could not convert task to sync item: %v", err)
	}
	it.Deleted = true
	if err := deps.SyncRepo.Store(it); err != nil {
		return nil, fmt.Errorf("could not store sync item: %v", err)
	}

	if err := deps.LocalIDRepo.Delete(id); err != nil {
		return nil, fmt.Errorf("could not delete local id: %v", err)
	}

	if err := deps.TaskRepo.Delete(id); err != nil {
		return nil, fmt.Errorf("could not delete task: %v", err)
	}

	return nil, nil
}
