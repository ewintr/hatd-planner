package command

import (
	"errors"
	"fmt"
	"strconv"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/format"
	"go-mod.ewintr.nl/planner/plan/storage"
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

func (s *Show) Do(deps Dependencies) (CommandResult, error) {
	id, err := deps.LocalIDRepo.FindOne(s.args.localID)
	switch {
	case errors.Is(err, storage.ErrNotFound):
		return nil, fmt.Errorf("could not find local id")
	case err != nil:
		return nil, err
	}

	tsk, err := deps.TaskRepo.FindOne(id)
	if err != nil {
		return nil, fmt.Errorf("could not find task")
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
		{"date", sr.Task.Date.String()},
		{"time", sr.Task.Time.String()},
		{"duration", sr.Task.Duration.String()},
		{"recur", recurStr},
		// {"id", s.Task.ID},
	}

	return fmt.Sprintf("\n%s\n", format.Table(data))
}
