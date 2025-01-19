package task

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/cli/arg"
	"go-mod.ewintr.nl/planner/plan/command"
	"go-mod.ewintr.nl/planner/plan/format"
	"go-mod.ewintr.nl/planner/sync/client"
)

type AddArgs struct {
	fieldTPL map[string][]string
	Task     item.Task
}

func NewAddArgs() AddArgs {
	return AddArgs{
		fieldTPL: map[string][]string{
			"date":     {"d", "date", "on"},
			"time":     {"t", "time", "at"},
			"project":  {"p", "project"},
			"duration": {"dur", "duration", "for"},
			"recurrer": {"rec", "recurrer"},
		},
	}
}

func (aa AddArgs) Parse(main []string, fields map[string]string) (command.Command, error) {
	if len(main) == 0 || !slices.Contains([]string{"add", "a", "new", "n"}, main[0]) {
		return nil, command.ErrWrongCommand
	}

	main = main[1:]
	if len(main) == 0 {
		return nil, fmt.Errorf("%w: title is required for add", command.ErrInvalidArg)
	}
	fields, err := arg.ResolveFields(fields, aa.fieldTPL)
	if err != nil {
		return nil, err
	}

	tsk := item.Task{
		ID: uuid.New().String(),
		TaskBody: item.TaskBody{
			Title: strings.Join(main, " "),
		},
	}

	if val, ok := fields["project"]; ok {
		tsk.Project = val
	}
	if val, ok := fields["date"]; ok {
		d := item.NewDateFromString(val)
		if d.IsZero() {
			return nil, fmt.Errorf("%w: could not parse date", command.ErrInvalidArg)
		}
		tsk.Date = d
	}
	if val, ok := fields["time"]; ok {
		t := item.NewTimeFromString(val)
		if t.IsZero() {
			return nil, fmt.Errorf("%w: could not parse time", command.ErrInvalidArg)
		}
		tsk.Time = t
	}
	if val, ok := fields["duration"]; ok {
		d, err := time.ParseDuration(val)
		if err != nil {
			return nil, fmt.Errorf("%w: could not parse duration", command.ErrInvalidArg)
		}
		tsk.Duration = d
	}
	if val, ok := fields["recurrer"]; ok {
		rec := item.NewRecurrer(val)
		if rec == nil {
			return nil, fmt.Errorf("%w: could not parse recurrer", command.ErrInvalidArg)
		}
		tsk.Recurrer = rec
		tsk.RecurNext = tsk.Recurrer.First()
	}

	return &Add{
		Args: AddArgs{
			Task: tsk,
		},
	}, nil
}

type Add struct {
	Args AddArgs
}

func (a Add) Do(repos command.Repositories, _ client.Client) (command.CommandResult, error) {
	tx, err := repos.Begin()
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %v", err)
	}
	defer tx.Rollback()

	if err := repos.Task(tx).Store(a.Args.Task); err != nil {
		return nil, fmt.Errorf("could not store task: %v", err)
	}

	localID, err := repos.LocalID(tx).Next()
	if err != nil {
		return nil, fmt.Errorf("could not create next local id: %v", err)
	}
	if err := repos.LocalID(tx).Store(a.Args.Task.ID, localID); err != nil {
		return nil, fmt.Errorf("could not store local id: %v", err)
	}

	it, err := a.Args.Task.Item()
	if err != nil {
		return nil, fmt.Errorf("could not convert task to sync item: %v", err)
	}
	if err := repos.Sync(tx).Store(it); err != nil {
		return nil, fmt.Errorf("could not store sync item: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not add task: %v", err)
	}

	return AddResult{
		LocalID: localID,
	}, nil
}

type AddResult struct {
	LocalID int
}

func (ar AddResult) Render() string {
	return fmt.Sprintf("stored task %s", format.Bold(fmt.Sprintf("%d", ar.LocalID)))
}
