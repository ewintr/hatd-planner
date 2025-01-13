package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage"
	"go-mod.ewintr.nl/planner/sync/client"
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

func (s Sync) Do(repos Repositories, client client.Client) (CommandResult, error) {
	tx, err := repos.Begin()
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %v", err)
	}
	defer tx.Rollback()

	// local new and updated
	sendItems, err := repos.Sync(tx).FindAll()
	if err != nil {
		return nil, fmt.Errorf("could not get updated items: %v", err)
	}
	if err := client.Update(sendItems); err != nil {
		return nil, fmt.Errorf("could not send updated items: %v", err)
	}
	if err := repos.Sync(tx).DeleteAll(); err != nil {
		return nil, fmt.Errorf("could not clear updated items: %v", err)
	}

	// get new/updated items
	oldTS, err := repos.Sync(tx).LastUpdate()
	if err != nil {
		return nil, fmt.Errorf("could not find timestamp of last update: %v", err)
	}
	recItems, err := client.Updated([]item.Kind{item.KindTask}, oldTS)
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
			if err := repos.LocalID(tx).Delete(ri.ID); err != nil && !errors.Is(err, storage.ErrNotFound) {
				return nil, fmt.Errorf("could not delete local id: %v", err)
			}
			if err := repos.Task(tx).Delete(ri.ID); err != nil && !errors.Is(err, storage.ErrNotFound) {
				return nil, fmt.Errorf("could not delete task: %v", err)
			}
			continue
		}
		updated = append(updated, ri)
	}

	lidMap, err := repos.LocalID(tx).FindAll()
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
		if err := repos.Task(tx).Store(tsk); err != nil {
			return nil, fmt.Errorf("could not store task: %v", err)
		}
		lid, ok := lidMap[u.ID]
		if !ok {
			lid, err = repos.LocalID(tx).Next()
			if err != nil {
				return nil, fmt.Errorf("could not get next local id: %v", err)
			}

			if err := repos.LocalID(tx).Store(u.ID, lid); err != nil {
				return nil, fmt.Errorf("could not store local id: %v", err)
			}
		}
	}

	if err := repos.Sync(tx).SetLastUpdate(newTS); err != nil {
		return nil, fmt.Errorf("could not store update timestamp: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not sync items: %v", err)
	}

	return SyncResult{}, nil
}

type SyncResult struct{}

func (sr SyncResult) Render() string { return "tasks synced" }
