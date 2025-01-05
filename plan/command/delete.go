package command

import (
	"fmt"
	"slices"
	"strconv"
)

type DeleteArgs struct {
	LocalID int
}

func NewDeleteArgs() DeleteArgs {
	return DeleteArgs{}
}

func (da DeleteArgs) Parse(main []string, flags map[string]string) (Command, error) {
	if len(main) != 2 {
		return nil, ErrWrongCommand
	}
	aliases := []string{"d", "delete", "done"}
	var localIDStr string
	switch {
	case slices.Contains(aliases, main[0]):
		localIDStr = main[1]
	case slices.Contains(aliases, main[1]):
		localIDStr = main[0]
	default:
		return nil, ErrWrongCommand
	}

	localID, err := strconv.Atoi(localIDStr)
	if err != nil {
		return nil, fmt.Errorf("not a local id: %v", main[1])
	}

	return &Delete{
		Args: DeleteArgs{
			LocalID: localID,
		},
	}, nil
}

type Delete struct {
	Args DeleteArgs
}

func (del Delete) Do(deps Dependencies) (CommandResult, error) {
	var id string
	idMap, err := deps.LocalIDRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("could not get local ids: %v", err)
	}
	for tskID, lid := range idMap {
		if del.Args.LocalID == lid {
			id = tskID
		}
	}
	if id == "" {
		return nil, fmt.Errorf("could not find local id")
	}

	tsk, err := deps.TaskRepo.FindOne(id)
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

	return DeleteResult{}, nil
}

type DeleteResult struct{}

func (dr DeleteResult) Render() string {
	return "task deleted"
}
