package command

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/format"
	"go-mod.ewintr.nl/planner/plan/storage"
	"go-mod.ewintr.nl/planner/sync/client"
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

func (u Update) Do(repos Repositories, _ client.Client) (CommandResult, error) {
	tx, err := repos.Begin()
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %v", err)
	}
	defer tx.Rollback()

	id, err := repos.LocalID(tx).FindOne(u.args.LocalID)
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
	changes := make(map[string]string)
	oldTitle := tsk.Title

	if u.args.Title != "" {
		tsk.Title = u.args.Title
		changes["title"] = u.args.Title
	}
	if slices.Contains(u.args.NeedUpdate, "project") {
		tsk.Project = u.args.Project
		changes["project"] = tsk.Project
	}
	if slices.Contains(u.args.NeedUpdate, "date") {
		tsk.Date = u.args.Date
		changes["date"] = tsk.Date.String()
	}
	if slices.Contains(u.args.NeedUpdate, "time") {
		tsk.Time = u.args.Time
		changes["time"] = tsk.Time.String()
	}
	if slices.Contains(u.args.NeedUpdate, "duration") {
		tsk.Duration = u.args.Duration
		changes["duration"] = tsk.Duration.String()
	}
	if slices.Contains(u.args.NeedUpdate, "recurrer") {
		tsk.Recurrer = u.args.Recurrer
		tsk.RecurNext = tsk.Recurrer.First()
		changes["recurrer"] = tsk.Recurrer.String()
	}

	if !tsk.Valid() {
		return nil, fmt.Errorf("task is unvalid")
	}

	if err := repos.Task(tx).Store(tsk); err != nil {
		return nil, fmt.Errorf("could not store task: %v", err)
	}

	it, err := tsk.Item()
	if err != nil {
		return nil, fmt.Errorf("could not convert task to sync item: %v", err)
	}
	if err := repos.Sync(tx).Store(it); err != nil {
		return nil, fmt.Errorf("could not store sync item: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not update task: %v", err)
	}

	return UpdateResult{
		Title:   oldTitle,
		Changes: changes,
	}, nil
}

type UpdateResult struct {
	Title   string
	Changes map[string]string
}

func (ur UpdateResult) Render() string {
	chStr := make([]string, 0, len(ur.Changes))
	for k, v := range ur.Changes {
		chStr = append(chStr, fmt.Sprintf("%s to %s", format.Bold(k), format.Bold(v)))
	}
	return fmt.Sprintf("updated task %s, set %s", format.Bold(ur.Title), strings.Join(chStr, ", "))
}
