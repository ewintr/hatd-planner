package command

import (
	"fmt"
	"strconv"
	"strings"

	"go-mod.ewintr.nl/planner/plan/storage"
)

type Update struct {
	localIDRepo storage.LocalID
	eventRepo   storage.Event
	syncRepo    storage.Sync
	argSet      *ArgSet
	localID     int
}

func NewUpdate(localIDRepo storage.LocalID, eventRepo storage.Event, syncRepo storage.Sync) Command {
	return &Update{
		localIDRepo: localIDRepo,
		eventRepo:   eventRepo,
		syncRepo:    syncRepo,
		argSet: &ArgSet{
			Flags: map[string]Flag{
				FlagTitle: &FlagString{},
				FlagOn:    &FlagDate{},
				FlagAt:    &FlagTime{},
				FlagFor:   &FlagDuration{},
				FlagRec:   &FlagRecurrer{},
			},
		},
	}
}

func (update *Update) Execute(main []string, flags map[string]string) error {
	if len(main) < 2 || main[0] != "update" {
		return ErrWrongCommand
	}
	localID, err := strconv.Atoi(main[1])
	if err != nil {
		return fmt.Errorf("not a local id: %v", main[1])
	}
	update.localID = localID
	main = main[2:]

	as := update.argSet
	as.Main = strings.Join(main, " ")
	for k := range as.Flags {
		v, ok := flags[k]
		if !ok {
			continue
		}
		if err := as.Set(k, v); err != nil {
			return fmt.Errorf("could not set %s: %v", k, err)
		}
	}
	update.argSet = as

	return update.do()
}

func (update *Update) do() error {
	as := update.argSet
	var id string
	idMap, err := update.localIDRepo.FindAll()
	if err != nil {
		return fmt.Errorf("could not get local ids: %v", err)
	}
	for eid, lid := range idMap {
		if update.localID == lid {
			id = eid
		}
	}
	if id == "" {
		return fmt.Errorf("could not find local id")
	}

	e, err := update.eventRepo.Find(id)
	if err != nil {
		return fmt.Errorf("could not find event")
	}

	if as.Main != "" {
		e.Title = as.Main
	}
	if as.IsSet(FlagOn) {
		e.Date = as.GetDate(FlagOn)
	}
	if as.IsSet(FlagAt) {
		e.Time = as.GetTime(FlagAt)
	}
	if as.IsSet(FlagFor) {
		e.Duration = as.GetDuration(FlagFor)
	}
	if as.IsSet(FlagRec) {
		e.Recurrer = as.GetRecurrer(FlagRec)
	}

	if !e.Valid() {
		return fmt.Errorf("event is unvalid")
	}

	if err := update.eventRepo.Store(e); err != nil {
		return fmt.Errorf("could not store event: %v", err)
	}

	it, err := e.Item()
	if err != nil {
		return fmt.Errorf("could not convert event to sync item: %v", err)
	}
	if err := update.syncRepo.Store(it); err != nil {
		return fmt.Errorf("could not store sync item: %v", err)
	}

	return nil
}
