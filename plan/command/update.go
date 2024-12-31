package command

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage"
)

type UpdateArgs struct {
	fieldTPL map[string][]string
	LocalID  int
	Title    string
	Project  string
	Date     item.Date
	Time     item.Time
	Duration time.Duration
	Recurrer item.Recurrer
}

func NewUpdateArgs() UpdateArgs {
	return UpdateArgs{
		fieldTPL: map[string][]string{
			"project":  {"p", "project"},
			"date":     {"d", "date", "on"},
			"time":     {"t", "time", "at"},
			"duration": {"dur", "duration", "for"},
			"recurrer": {"rec", "recurrer"},
		},
	}
}

func (ua UpdateArgs) Parse(main []string, fields map[string]string) (Command, error) {
	if len(main) < 2 || main[0] != "update" {
		return nil, ErrWrongCommand
	}
	localID, err := strconv.Atoi(main[1])
	if err != nil {
		return nil, fmt.Errorf("not a local id: %v", main[1])
	}
	fields, err = ResolveFields(fields, ua.fieldTPL)
	if err != nil {
		return nil, err
	}
	args := UpdateArgs{
		LocalID: localID,
		Title:   strings.Join(main[2:], " "),
	}

	if val, ok := fields["project"]; ok {
		args.Project = val
	}
	if val, ok := fields["date"]; ok {
		d := item.NewDateFromString(val)
		if d.IsZero() {
			return nil, fmt.Errorf("%w: could not parse date", ErrInvalidArg)
		}
		args.Date = d
	}
	if val, ok := fields["time"]; ok {
		t := item.NewTimeFromString(val)
		if t.IsZero() {
			return nil, fmt.Errorf("%w: could not parse time", ErrInvalidArg)
		}
		args.Time = t
	}
	if val, ok := fields["duration"]; ok {
		d, err := time.ParseDuration(val)
		if err != nil {
			return nil, fmt.Errorf("%w: could not parse duration", ErrInvalidArg)
		}
		args.Duration = d
	}
	if val, ok := fields["recurrer"]; ok {
		rec := item.NewRecurrer(val)
		if rec == nil {
			return nil, fmt.Errorf("%w: could not parse recurrer", ErrInvalidArg)
		}
		args.Recurrer = rec
	}

	return &Update{args}, nil
}

type Update struct {
	args UpdateArgs
}

func (u *Update) Do(deps Dependencies) (CommandResult, error) {
	id, err := deps.LocalIDRepo.FindOne(u.args.LocalID)
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

	if u.args.Title != "" {
		tsk.Title = u.args.Title
	}
	if u.args.Project != "" {
		tsk.Project = u.args.Project
	}
	if !u.args.Date.IsZero() {
		tsk.Date = u.args.Date
	}
	if !u.args.Time.IsZero() {
		tsk.Time = u.args.Time
	}
	if u.args.Duration != 0 {
		tsk.Duration = u.args.Duration
	}
	if u.args.Recurrer != nil {
		tsk.Recurrer = u.args.Recurrer
		tsk.RecurNext = tsk.Recurrer.First()
	}

	if !tsk.Valid() {
		return nil, fmt.Errorf("task is unvalid")
	}

	if err := deps.TaskRepo.Store(tsk); err != nil {
		return nil, fmt.Errorf("could not store task: %v", err)
	}

	it, err := tsk.Item()
	if err != nil {
		return nil, fmt.Errorf("could not convert task to sync item: %v", err)
	}
	if err := deps.SyncRepo.Store(it); err != nil {
		return nil, fmt.Errorf("could not store sync item: %v", err)
	}

	return UpdateResult{}, nil
}

type UpdateResult struct{}

func (ur UpdateResult) Render() string {
	return "task updated"
}
