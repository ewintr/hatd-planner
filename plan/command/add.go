package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage"
)

type Add struct {
	localIDRepo storage.LocalID
	eventRepo   storage.Event
	syncRepo    storage.Sync
	argSet      *ArgSet
}

func NewAdd(localRepo storage.LocalID, eventRepo storage.Event, syncRepo storage.Sync) Command {
	return &Add{
		localIDRepo: localRepo,
		eventRepo:   eventRepo,
		syncRepo:    syncRepo,
		argSet: &ArgSet{
			Flags: map[string]Flag{
				FlagOn:  &FlagDate{},
				FlagAt:  &FlagTime{},
				FlagFor: &FlagDuration{},
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
	start := as.GetTime(FlagOn)
	if as.IsSet(FlagAt) {
		at := as.GetTime(FlagAt)
		h := time.Duration(at.Hour()) * time.Hour
		m := time.Duration(at.Minute()) * time.Minute
		start = start.Add(h).Add(m)
	}

	e := item.Event{
		ID: uuid.New().String(),
		EventBody: item.EventBody{
			Title: as.Main,
			Start: start,
		},
	}

	if as.IsSet(FlagFor) {
		e.Duration = as.GetDuration(FlagFor)
	}
	if err := add.eventRepo.Store(e); err != nil {
		return fmt.Errorf("could not store event: %v", err)
	}

	localID, err := add.localIDRepo.Next()
	if err != nil {
		return fmt.Errorf("could not create next local id: %v", err)
	}
	if err := add.localIDRepo.Store(e.ID, localID); err != nil {
		return fmt.Errorf("could not store local id: %v", err)
	}

	it, err := e.Item()
	if err != nil {
		return fmt.Errorf("could not convert event to sync item: %v", err)
	}
	if err := add.syncRepo.Store(it); err != nil {
		return fmt.Errorf("could not store sync item: %v", err)
	}

	return nil
}
