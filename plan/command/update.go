package command

import (
	"fmt"
	"strconv"
	"strings"

	"go-mod.ewintr.nl/planner/plan/storage"
)

type Update struct {
	localIDRepo storage.LocalID
	taskRepo    storage.Task
	syncRepo    storage.Sync
	argSet      *ArgSet
	localID     int
}

func NewUpdate(localIDRepo storage.LocalID, taskRepo storage.Task, syncRepo storage.Sync) Command {
	return &Update{
		localIDRepo: localIDRepo,
		taskRepo:    taskRepo,
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
	for tid, lid := range idMap {
		if update.localID == lid {
			id = tid
		}
	}
	if id == "" {
		return fmt.Errorf("could not find local id")
	}

	tsk, err := update.taskRepo.Find(id)
	if err != nil {
		return fmt.Errorf("could not find task")
	}

	if as.Main != "" {
		tsk.Title = as.Main
	}
	if as.IsSet(FlagOn) {
		tsk.Date = as.GetDate(FlagOn)
	}
	if as.IsSet(FlagAt) {
		tsk.Time = as.GetTime(FlagAt)
	}
	if as.IsSet(FlagFor) {
		tsk.Duration = as.GetDuration(FlagFor)
	}
	if as.IsSet(FlagRec) {
		tsk.Recurrer = as.GetRecurrer(FlagRec)
	}

	if !tsk.Valid() {
		return fmt.Errorf("task is unvalid")
	}

	if err := update.taskRepo.Store(tsk); err != nil {
		return fmt.Errorf("could not store task: %v", err)
	}

	it, err := tsk.Item()
	if err != nil {
		return fmt.Errorf("could not convert task to sync item: %v", err)
	}
	if err := update.syncRepo.Store(it); err != nil {
		return fmt.Errorf("could not store sync item: %v", err)
	}

	return nil
}
