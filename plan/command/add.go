package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go-mod.ewintr.nl/planner/item"
)

type AddArgs struct {
	fieldTPL map[string][]string
	task     item.Task
}

func NewAddArgs() AddArgs {
	return AddArgs{
		fieldTPL: map[string][]string{
			"date":     []string{"d", "date", "on"},
			"time":     []string{"t", "time", "at"},
			"duration": []string{"dur", "duration", "for"},
			"recurrer": []string{"rec", "recurrer"},
		},
	}
}

func (aa AddArgs) Parse(main []string, fields map[string]string) (Command, error) {
	if len(main) == 0 || main[0] != "add" {
		return nil, ErrWrongCommand
	}
	main = main[1:]
	if len(main) == 0 {
		return nil, fmt.Errorf("%w: title is required for add", ErrInvalidArg)
	}
	fields, err := ResolveFields(fields, aa.fieldTPL)
	if err != nil {
		return nil, err
	}

	tsk := item.Task{
		ID: uuid.New().String(),
		TaskBody: item.TaskBody{
			Title: strings.Join(main, ","),
		},
	}

	if val, ok := fields["date"]; ok {
		d := item.NewDateFromString(val)
		if d.IsZero() {
			return nil, fmt.Errorf("%w: could not parse date", ErrInvalidArg)
		}
		tsk.Date = d
	}
	if val, ok := fields["time"]; ok {
		t := item.NewTimeFromString(val)
		if t.IsZero() {
			return nil, fmt.Errorf("%w: could not parse time", ErrInvalidArg)
		}
		tsk.Time = t
	}
	if val, ok := fields["duration"]; ok {
		d, err := time.ParseDuration(val)
		if err != nil {
			return nil, fmt.Errorf("%w: could not parse duration", ErrInvalidArg)
		}
		tsk.Duration = d
	}
	if val, ok := fields["recurrer"]; ok {
		rec := item.NewRecurrer(val)
		if rec == nil {
			return nil, fmt.Errorf("%w: could not parse recurrer", ErrInvalidArg)
		}
		tsk.Recurrer = rec
		tsk.RecurNext = tsk.Recurrer.First()
	}

	return &Add{
		args: AddArgs{
			task: tsk,
		},
	}, nil
}

type Add struct {
	args AddArgs
}

func (a *Add) Do(deps Dependencies) ([][]string, error) {
	if err := deps.TaskRepo.Store(a.args.task); err != nil {
		return nil, fmt.Errorf("could not store event: %v", err)
	}

	localID, err := deps.LocalIDRepo.Next()
	if err != nil {
		return nil, fmt.Errorf("could not create next local id: %v", err)
	}
	if err := deps.LocalIDRepo.Store(a.args.task.ID, localID); err != nil {
		return nil, fmt.Errorf("could not store local id: %v", err)
	}

	it, err := a.args.task.Item()
	if err != nil {
		return nil, fmt.Errorf("could not convert event to sync item: %v", err)
	}
	if err := deps.SyncRepo.Store(it); err != nil {
		return nil, fmt.Errorf("could not store sync item: %v", err)
	}

	return nil, nil
}
