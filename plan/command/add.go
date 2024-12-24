package command

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage"
)

type Add struct {
	localIDRepo storage.LocalID
	taskRepo    storage.Task
	syncRepo    storage.Sync
	argSet      *ArgSet
}

func NewAdd(localRepo storage.LocalID, taskRepo storage.Task, syncRepo storage.Sync) Command {
	return &Add{
		localIDRepo: localRepo,
		taskRepo:    taskRepo,
		syncRepo:    syncRepo,
		argSet: &ArgSet{
			Flags: map[string]Flag{
				FlagOn:  &FlagDate{},
				FlagAt:  &FlagTime{},
				FlagFor: &FlagDuration{},
				FlagRec: &FlagRecurrer{},
			},
		},
	}
}

func (add *Add) Execute(main []string, flags map[string]string) error {
	if len(main) == 0 || main[0] != "add" {
		return ErrWrongCommand
	}
	as := add.argSet
	if len(main) > 1 {
		as.Main = strings.Join(main[1:], " ")
	}
	for k := range as.Flags {
		v, ok := flags[k]
		if !ok {
			continue
		}
		if err := as.Set(k, v); err != nil {
			return fmt.Errorf("could not set %s: %v", k, err)
		}
	}
	if as.Main == "" {
		return fmt.Errorf("%w: title is required", ErrInvalidArg)
	}
	if !as.IsSet(FlagOn) {
		return fmt.Errorf("%w: date is required", ErrInvalidArg)
	}
	if !as.IsSet(FlagAt) && as.IsSet(FlagFor) {
		return fmt.Errorf("%w: can not have duration without start time", ErrInvalidArg)
	}
	if as.IsSet(FlagAt) && !as.IsSet(FlagFor) {
		if err := as.Flags[FlagFor].Set("1h"); err != nil {
			return fmt.Errorf("could not set duration to one hour")
		}
	}
	if !as.IsSet(FlagAt) && !as.IsSet(FlagFor) {
		if err := as.Flags[FlagFor].Set("24h"); err != nil {
			return fmt.Errorf("could not set duration to 24 hours")
		}
	}

	return add.do()
}

func (add *Add) do() error {
	as := add.argSet
	rec := as.GetRecurrer(FlagRec)
	tsk := item.Task{
		ID:       uuid.New().String(),
		Date:     as.GetDate(FlagOn),
		Recurrer: rec,
		TaskBody: item.TaskBody{
			Title:    as.Main,
			Time:     as.GetTime(FlagAt),
			Duration: as.GetDuration(FlagFor),
		},
	}
	if rec != nil {
		tsk.RecurNext = rec.First()
	}

	if err := add.taskRepo.Store(tsk); err != nil {
		return fmt.Errorf("could not store event: %v", err)
	}

	localID, err := add.localIDRepo.Next()
	if err != nil {
		return fmt.Errorf("could not create next local id: %v", err)
	}
	if err := add.localIDRepo.Store(tsk.ID, localID); err != nil {
		return fmt.Errorf("could not store local id: %v", err)
	}

	it, err := tsk.Item()
	if err != nil {
		return fmt.Errorf("could not convert event to sync item: %v", err)
	}
	if err := add.syncRepo.Store(it); err != nil {
		return fmt.Errorf("could not store sync item: %v", err)
	}

	return nil
}
