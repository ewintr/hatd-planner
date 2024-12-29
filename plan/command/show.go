package command

import (
	"errors"
	"fmt"
	"strconv"

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

func (s *Show) Do(deps Dependencies) ([][]string, error) {
	id, err := deps.LocalIDRepo.FindOne(s.args.localID)
	switch {
	case errors.Is(err, storage.ErrNotFound):
		return nil, fmt.Errorf("could not find local id")
	case err != nil:
		return nil, err
	}

	tsk, err := deps.TaskRepo.Find(id)
	if err != nil {
		return nil, fmt.Errorf("could not find task")
	}

	var recurStr string
	if tsk.Recurrer != nil {
		recurStr = tsk.Recurrer.String()
	}
	data := [][]string{
		{"title", tsk.Title},
		{"local id", fmt.Sprintf("%d", s.args.localID)},
		{"date", tsk.Date.String()},
		{"time", tsk.Time.String()},
		{"duration", tsk.Duration.String()},
		{"recur", recurStr},
		// {"id", tsk.ID},
	}

	return data, nil
}
