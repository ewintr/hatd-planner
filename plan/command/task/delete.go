package task

import (
	"fmt"
	"slices"
	"strconv"

	"go-mod.ewintr.nl/planner/plan/command"
	"go-mod.ewintr.nl/planner/plan/format"
	"go-mod.ewintr.nl/planner/sync/client"
)

type DeleteArgs struct {
	LocalID int
}

func NewDeleteArgs() DeleteArgs {
	return DeleteArgs{}
}

func (da DeleteArgs) Parse(main []string, flags map[string]string) (command.Command, error) {
	if len(main) != 2 {
		return nil, command.ErrWrongCommand
	}
	aliases := []string{"d", "delete", "done"}
	var localIDStr string
	switch {
	case slices.Contains(aliases, main[0]):
		localIDStr = main[1]
	case slices.Contains(aliases, main[1]):
		localIDStr = main[0]
	default:
		return nil, command.ErrWrongCommand
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

func (del Delete) Do(repos command.Repositories, _ client.Client) (command.CommandResult, error) {
	tx, err := repos.Begin()
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %v", err)
	}
	defer tx.Rollback()

	var id string
	idMap, err := repos.LocalID(tx).FindAll()
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

	tsk, err := repos.Task(tx).FindOne(id)
	if err != nil {
		return nil, fmt.Errorf("could not get task: %v", err)
	}

	it, err := tsk.Item()
	if err != nil {
		return nil, fmt.Errorf("could not convert task to sync item: %v", err)
	}
	it.Deleted = true
	if err := repos.Sync(tx).Store(it); err != nil {
		return nil, fmt.Errorf("could not store sync item: %v", err)
	}

	if err := repos.LocalID(tx).Delete(id); err != nil {
		return nil, fmt.Errorf("could not delete local id: %v", err)
	}

	if err := repos.Task(tx).Delete(id); err != nil {
		return nil, fmt.Errorf("could not delete task: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not delete task: %v", err)
	}

	return DeleteResult{
		Title: tsk.Title,
	}, nil
}

type DeleteResult struct {
	Title string
}

func (dr DeleteResult) Render() string {
	return fmt.Sprintf("removed task %s", format.Bold(dr.Title))
}
