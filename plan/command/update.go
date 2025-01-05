package command

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage"
)

type UpdateArgs struct {
	fieldTPL   map[string][]string
	NeedUpdate []string
	LocalID    int
	Title      string
	Project    string
	Date       item.Date
	Time       item.Time
	Duration   time.Duration
	Recurrer   item.Recurrer
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
	if len(main) < 2 {
		return nil, ErrWrongCommand
	}
	aliases := []string{"u", "update", "m", "mod"}
	var localIDStr string
	switch {
	case slices.Contains(aliases, main[0]):
		localIDStr = main[1]
	case slices.Contains(aliases, main[1]):
		localIDStr = main[0]
	default:
		return nil, ErrWrongCommand
	}
	localID, err := strconv.Atoi(localIDStr)
	if err != nil {
		return nil, fmt.Errorf("not a local id: %v", main[1])
	}
	fields, err = ResolveFields(fields, ua.fieldTPL)
	if err != nil {
		return nil, err
	}
	args := UpdateArgs{
		NeedUpdate: make([]string, 0),
		LocalID:    localID,
		Title:      strings.Join(main[2:], " "),
	}

	if val, ok := fields["project"]; ok {
		args.NeedUpdate = append(args.NeedUpdate, "project")
		args.Project = val
	}
	if val, ok := fields["date"]; ok {
		args.NeedUpdate = append(args.NeedUpdate, "date")
		if val != "" {
			d := item.NewDateFromString(val)
			if d.IsZero() {
				return nil, fmt.Errorf("%w: could not parse date", ErrInvalidArg)
			}
			args.Date = d
		}
	}
	if val, ok := fields["time"]; ok {
		args.NeedUpdate = append(args.NeedUpdate, "time")
		if val != "" {
			t := item.NewTimeFromString(val)
			if t.IsZero() {
				return nil, fmt.Errorf("%w: could not parse time", ErrInvalidArg)
			}
			args.Time = t
		}
	}
	if val, ok := fields["duration"]; ok {
		args.NeedUpdate = append(args.NeedUpdate, "duration")
		if val != "" {
			d, err := time.ParseDuration(val)
			if err != nil {
				return nil, fmt.Errorf("%w: could not parse duration", ErrInvalidArg)
			}
			args.Duration = d
		}
	}
	if val, ok := fields["recurrer"]; ok {
		args.NeedUpdate = append(args.NeedUpdate, "recurrer")
		if val != "" {
			rec := item.NewRecurrer(val)
			if rec == nil {
				return nil, fmt.Errorf("%w: could not parse recurrer", ErrInvalidArg)
			}
			args.Recurrer = rec
		}
	}

	return &Update{args}, nil
}

type Update struct {
	args UpdateArgs
}

func (u Update) Do(deps Dependencies) (CommandResult, error) {
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
	if slices.Contains(u.args.NeedUpdate, "project") {
		tsk.Project = u.args.Project
	}
	if slices.Contains(u.args.NeedUpdate, "date") {
		tsk.Date = u.args.Date
	}
	if slices.Contains(u.args.NeedUpdate, "time") {
		tsk.Time = u.args.Time
	}
	if slices.Contains(u.args.NeedUpdate, "duration") {
		tsk.Duration = u.args.Duration
	}
	if slices.Contains(u.args.NeedUpdate, "recurrer") {
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
