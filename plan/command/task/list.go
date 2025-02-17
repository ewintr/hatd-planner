package task

import (
	"fmt"
	"slices"
	"sort"
	"time"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/cli/arg"
	"go-mod.ewintr.nl/planner/plan/command"
	"go-mod.ewintr.nl/planner/plan/format"
	"go-mod.ewintr.nl/planner/plan/storage"
	"go-mod.ewintr.nl/planner/sync/client"
)

type ListArgs struct {
	fieldTPL    map[string][]string
	HasRecurrer bool
	From        item.Date
	To          item.Date
	Project     string
}

func NewListArgs() ListArgs {
	return ListArgs{
		fieldTPL: map[string][]string{
			"project":   {"p", "project"},
			"from":      {"f", "from"},
			"to":        {"t", "to"},
			"recurring": {"rec", "recurring"},
		},
	}
}

func (la ListArgs) Parse(main []string, fields map[string]string) (command.Command, error) {
	if len(main) > 1 {
		return nil, command.ErrWrongCommand
	}

	fields, err := arg.ResolveFields(fields, la.fieldTPL)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	today := item.NewDate(now.Year(), int(now.Month()), now.Day())

	switch len(main) {
	case 0:
		// fields["to"] = today.String()
	case 1:
		switch {
		case slices.Contains([]string{"tod", "today"}, main[0]):
			fields["to"] = today.String()
		case slices.Contains([]string{"tom", "tomorrow"}, main[0]):
			fields["from"] = today.Add(1).String()
			fields["to"] = today.Add(1).String()
		case main[0] == "week":
			fields["from"] = today.String()
			fields["to"] = today.Add(7).String()
		case main[0] == "recur":
			fields["recurrer"] = "true"
		// case main[0] == "list":
		// 	fields["from"] = today.String()
		// 	fields["to"] = today.String()
		default:
			return nil, command.ErrWrongCommand
		}
	}

	var fromDate, toDate item.Date
	var hasRecurrer bool
	var project string
	if val, ok := fields["from"]; ok {
		fromDate = item.NewDateFromString(val)
	}
	if val, ok := fields["to"]; ok {
		toDate = item.NewDateFromString(val)
	}
	if val, ok := fields["recurrer"]; ok && val == "true" {
		hasRecurrer = true
	}
	if val, ok := fields["project"]; ok {
		project = val
	}

	return List{
		Args: ListArgs{
			HasRecurrer: hasRecurrer,
			From:        fromDate,
			To:          toDate,
			Project:     project,
		},
	}, nil
}

type List struct {
	Args ListArgs
}

func (list List) Do(repos command.Repositories, _ client.Client) (command.CommandResult, error) {
	tx, err := repos.Begin()
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %v", err)
	}
	defer tx.Rollback()

	localIDs, err := repos.LocalID(tx).FindAll()
	if err != nil {
		return nil, fmt.Errorf("could not get local ids: %v", err)
	}
	all, err := repos.Task(tx).FindMany(storage.TaskListParams{
		HasRecurrer: list.Args.HasRecurrer,
		From:        list.Args.From,
		To:          list.Args.To,
		Project:     list.Args.Project,
	})
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
	if len(lr.Tasks) == 0 {
		return "\nno tasks to display\n"
	}

	sort.Slice(lr.Tasks, func(i, j int) bool {
		if lr.Tasks[i].Task.Date.After(lr.Tasks[j].Task.Date) {
			return false
		}
		if lr.Tasks[j].Task.Date.After(lr.Tasks[i].Task.Date) {
			return true
		}
		if lr.Tasks[i].Task.Project < lr.Tasks[j].Task.Project {
			return true
		}
		if lr.Tasks[i].Task.Project > lr.Tasks[j].Task.Project {
			return false
		}
		if lr.Tasks[i].Task.Recurrer == nil && lr.Tasks[j].Task.Recurrer != nil {
			return true
		}
		if lr.Tasks[i].Task.Recurrer != nil && lr.Tasks[j].Task.Recurrer == nil {
			return false
		}
		return lr.Tasks[i].LocalID < lr.Tasks[j].LocalID
	})

	var showRec, showTime, showDur bool
	for _, tl := range lr.Tasks {
		if tl.Task.Recurrer != nil {
			showRec = true
		}
		if tl.Task.Duration > time.Duration(0) {
			showDur = true
		}
		if !tl.Task.Time.IsZero() {
			showTime = true
		}
	}

	title := []string{"id"}
	if showRec {
		title = append(title, "rec")
	}
	title = append(title, "project", "date")
	if showTime {
		title = append(title, "time")
	}
	if showDur {
		title = append(title, "dur")
	}
	title = append(title, "title")

	data := [][]string{title}
	for _, tl := range lr.Tasks {
		row := []string{fmt.Sprintf("%d", tl.LocalID)}
		if showRec {
			recStr := ""
			if tl.Task.Recurrer != nil {
				recStr = "*"
			}
			row = append(row, recStr)
		}
		row = append(row, tl.Task.Project, tl.Task.Date.String())
		if showTime {
			row = append(row, tl.Task.Time.String())
		}
		if showDur {
			durStr := ""
			if tl.Task.Duration > time.Duration(0) {
				durStr = tl.Task.Duration.String()
			}
			row = append(row, durStr)
		}
		row = append(row, tl.Task.Title)
		data = append(data, row)
	}

	return fmt.Sprintf("\n%s\n", format.Table(data))
}
