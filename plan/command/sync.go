package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage"
)

type SyncArgs struct{}

func NewSyncArgs() SyncArgs {
	return SyncArgs{}
}

func (sa SyncArgs) Parse(main []string, flags map[string]string) (Command, error) {
	if len(main) == 0 || main[0] != "sync" {
		return nil, ErrWrongCommand
	}

	return &Sync{}, nil
}

type Sync struct{}

func (s Sync) Do(deps Dependencies) (CommandResult, error) {
	// local new and updated
	sendItems, err := deps.SyncRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("could not get updated items: %v", err)
	}
	if err := deps.SyncClient.Update(sendItems); err != nil {
		return nil, fmt.Errorf("could not send updated items: %v", err)
	}
	if err := deps.SyncRepo.DeleteAll(); err != nil {
		return nil, fmt.Errorf("could not clear updated items: %v", err)
	}

	// get new/updated items
	oldTS, err := deps.SyncRepo.LastUpdate()
	if err != nil {
		return nil, fmt.Errorf("could not find timestamp of last update: %v", err)
	}
	recItems, err := deps.SyncClient.Updated([]item.Kind{item.KindTask}, oldTS)
	if err != nil {
		return nil, fmt.Errorf("could not receive updates: %v", err)
	}

	updated := make([]item.Item, 0)
	var newTS time.Time
	for _, ri := range recItems {
		if ri.Updated.After(newTS) {
			newTS = ri.Updated
		}
		if ri.Deleted {
			if err := deps.LocalIDRepo.Delete(ri.ID); err != nil && !errors.Is(err, storage.ErrNotFound) {
				return nil, fmt.Errorf("could not delete local id: %v", err)
			}
			if err := deps.TaskRepo.Delete(ri.ID); err != nil && !errors.Is(err, storage.ErrNotFound) {
				return nil, fmt.Errorf("could not delete task: %v", err)
			}
			continue
		}
		updated = append(updated, ri)
	}

	lidMap, err := deps.LocalIDRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("could not get local ids: %v", err)
	}
	for _, u := range updated {
		var tskBody item.TaskBody
		if err := json.Unmarshal([]byte(u.Body), &tskBody); err != nil {
			return nil, fmt.Errorf("could not unmarshal task body: %v", err)
		}
		tsk := item.Task{
			ID:        u.ID,
			Date:      u.Date,
			Recurrer:  u.Recurrer,
			RecurNext: u.RecurNext,
			TaskBody:  tskBody,
		}
		if err := deps.TaskRepo.Store(tsk); err != nil {
			return nil, fmt.Errorf("could not store task: %v", err)
		}
		lid, ok := lidMap[u.ID]
		if !ok {
			lid, err = deps.LocalIDRepo.Next()
			if err != nil {
				return nil, fmt.Errorf("could not get next local id: %v", err)
			}

			if err := deps.LocalIDRepo.Store(u.ID, lid); err != nil {
				return nil, fmt.Errorf("could not store local id: %v", err)
			}
		}
	}

	if err := deps.SyncRepo.SetLastUpdate(newTS); err != nil {
		return nil, fmt.Errorf("could not store update timestamp: %v", err)
	}

	return SyncResult{}, nil
}

type SyncResult struct{}

func (sr SyncResult) Render() string { return "tasks synced" }
