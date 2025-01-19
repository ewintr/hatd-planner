package schedule

import (
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"
	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/cli/arg"
	"go-mod.ewintr.nl/planner/plan/command"
	"go-mod.ewintr.nl/planner/plan/format"
	"go-mod.ewintr.nl/planner/sync/client"
)

type AddArgs struct {
	fieldTPL map[string][]string
	Schedule item.Schedule
}

func NewAddArgs() AddArgs {
	return AddArgs{
		fieldTPL: map[string][]string{
			"date":     {"d", "date", "on"},
			"recurrer": {"rec", "recurrer"},
		},
	}
}

func (aa AddArgs) Parse(main []string, fields map[string]string) (command.Command, error) {
	if len(main) == 0 || !slices.Contains([]string{"s", "sched", "schedule"}, main[0]) {
		return nil, command.ErrWrongCommand
	}
	main = main[1:]
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

	sched := item.Schedule{
		ID: uuid.New().String(),
		ScheduleBody: item.ScheduleBody{
			Title: strings.Join(main, " "),
		},
	}

	if val, ok := fields["date"]; ok {
		d := item.NewDateFromString(val)
		if d.IsZero() {
			return nil, fmt.Errorf("%w: could not parse date", command.ErrInvalidArg)
		}
		sched.Date = d
	}
	if val, ok := fields["recurrer"]; ok {
		rec := item.NewRecurrer(val)
		if rec == nil {
			return nil, fmt.Errorf("%w: could not parse recurrer", command.ErrInvalidArg)
		}
		sched.Recurrer = rec
		sched.RecurNext = sched.Recurrer.First()
	}

	return &Add{
		Args: AddArgs{
			Schedule: sched,
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

	if err := repos.Schedule(tx).Store(a.Args.Schedule); err != nil {
		return nil, fmt.Errorf("could not store schedule: %v", err)
	}

	localID, err := repos.LocalID(tx).Next()
	if err != nil {
		return nil, fmt.Errorf("could not create next local id: %v", err)
	}
	if err := repos.LocalID(tx).Store(a.Args.Schedule.ID, localID); err != nil {
		return nil, fmt.Errorf("could not store local id: %v", err)
	}

	it, err := a.Args.Schedule.Item()
	if err != nil {
		return nil, fmt.Errorf("could not convert schedule to sync item: %v", err)
	}
	if err := repos.Sync(tx).Store(it); err != nil {
		return nil, fmt.Errorf("could not store sync item: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not add schedule: %v", err)
	}

	return AddResult{
		LocalID: localID,
	}, nil
}

type AddResult struct {
	LocalID int
}

func (ar AddResult) Render() string {
	return fmt.Sprintf("stored schedule %s", format.Bold(fmt.Sprintf("%d", ar.LocalID)))
}
