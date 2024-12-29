package command

import (
	"fmt"
	"slices"
	"time"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/format"
	"go-mod.ewintr.nl/planner/plan/storage"
)

type ListArgs struct {
	params storage.TaskListParams
}

func NewListArgs() ListArgs {
	return ListArgs{}
}

func (la ListArgs) Parse(main []string, flags map[string]string) (Command, error) {
	if len(main) > 2 {
		return nil, ErrWrongCommand
	}

	now := time.Now()
	today := item.NewDate(now.Year(), int(now.Month()), now.Day())
	tomorrow := item.NewDate(now.Year(), int(now.Month()), now.Day()+1)
	var date item.Date
	var includeBefore, recurrer bool

	switch len(main) {
	case 0:
		date = today
		includeBefore = true
	case 1:
		switch {
		case slices.Contains([]string{"today", "tod"}, main[0]):
			date = today
			includeBefore = true
		case slices.Contains([]string{"tomorrow", "tom"}, main[0]):
			date = tomorrow
		case main[0] == "list":
		default:
			return nil, ErrWrongCommand
		}
	case 2:
		if main[0] == "list" && main[1] == "recur" {
			recurrer = true
		} else {
			return nil, ErrWrongCommand
		}
	default:
		return nil, ErrWrongCommand
	}

	return &List{
		args: ListArgs{
			params: storage.TaskListParams{
				Date:          date,
				IncludeBefore: includeBefore,
				Recurrer:      recurrer,
			},
		},
	}, nil
}

type List struct {
	args ListArgs
}

func (list *List) Do(deps Dependencies) (CommandResult, error) {
	localIDs, err := deps.LocalIDRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("could not get local ids: %v", err)
	}
	all, err := deps.TaskRepo.FindMany(list.args.params)
	if err != nil {
		return nil, err
	}

	res := make([]TaskWithLID, 0, len(all))
	for _, tsk := range all {
		lid, ok := localIDs[tsk.ID]
		if !ok {
			return nil, fmt.Errorf("could not find local id for %s", tsk.ID)
		}
		res = append(res, TaskWithLID{
			LocalID: lid,
			Task:    tsk,
		})
	}
	return ListResult{
		Tasks: res,
	}, nil
}

type TaskWithLID struct {
	LocalID int
	Task    item.Task
}

type ListResult struct {
	Tasks []TaskWithLID
}

func (lr ListResult) Render() string {
	data := [][]string{{"id", "date", "dur", "title"}}
	for _, tl := range lr.Tasks {
		data = append(data, []string{fmt.Sprintf("%d", tl.LocalID), tl.Task.Date.String(), tl.Task.Duration.String(), tl.Task.Title})
	}

	return fmt.Sprintf("\n%s\n", format.Table(data))
}
