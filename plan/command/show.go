package command

import (
	"errors"
	"fmt"
	"strconv"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/format"
	"go-mod.ewintr.nl/planner/plan/storage"
	"go-mod.ewintr.nl/planner/sync/client"
)

type ShowArgs struct {
	localID int
}

func NewShowArgs() ShowArgs {
	return ShowArgs{}
}

func (sa ShowArgs) Parse(main []string, fields map[string]string) (Command, error) {
	if len(main) != 1 {
		return nil, ErrWrongCommand
	}
	lid, err := strconv.Atoi(main[0])
	if err != nil {
		return nil, ErrWrongCommand
	}

	return &Show{
		args: ShowArgs{
			localID: lid,
		},
	}, nil
}

type Show struct {
	args ShowArgs
}

func (s Show) Do(repos Repositories, _ client.Client) (CommandResult, error) {
	tx, err := repos.Begin()
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %v", err)
	}
	defer tx.Rollback()

	id, err := repos.LocalID(tx).FindOne(s.args.localID)
	switch {
	case errors.Is(err, storage.ErrNotFound):
		return nil, fmt.Errorf("could not find local id")
	case err != nil:
		return nil, err
	}

	tsk, err := repos.Task(tx).FindOne(id)
	if err != nil {
		return nil, fmt.Errorf("could not find task")
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not show task: %v", err)
	}

	return ShowResult{
		LocalID: s.args.localID,
		Task:    tsk,
	}, nil
}

type ShowResult struct {
	LocalID int
	Task    item.Task
}

func (sr ShowResult) Render() string {

	var recurStr string
	if sr.Task.Recurrer != nil {
		recurStr = sr.Task.Recurrer.String()
	}
	data := [][]string{
		{"title", sr.Task.Title},
		{"local id", fmt.Sprintf("%d", sr.LocalID)},
		{"project", sr.Task.Project},
		{"date", sr.Task.Date.String()},
		{"time", sr.Task.Time.String()},
		{"duration", sr.Task.Duration.String()},
		{"recur", recurStr},
		// {"id", s.Task.ID},
	}

	return fmt.Sprintf("\n%s\n", format.Table(data))
}
