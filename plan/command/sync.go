package command

import (
	"encoding/json"
	"errors"
	"fmt"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage"
	"go-mod.ewintr.nl/planner/sync/client"
)

type Sync struct {
	client      client.Client
	syncRepo    storage.Sync
	localIDRepo storage.LocalID
	taskRepo    storage.Task
}

func NewSync(client client.Client, syncRepo storage.Sync, localIDRepo storage.LocalID, taskRepo storage.Task) Command {
	return &Sync{
		client:      client,
		syncRepo:    syncRepo,
		localIDRepo: localIDRepo,
		taskRepo:    taskRepo,
	}
}

func (sync *Sync) Execute(main []string, flags map[string]string) error {
	if len(main) == 0 || main[0] != "sync" {
		return ErrWrongCommand
	}

	return sync.do()
}

func (sync *Sync) do() error {
	// local new and updated
	sendItems, err := sync.syncRepo.FindAll()
	if err != nil {
		return fmt.Errorf("could not get updated items: %v", err)
	}
	if err := sync.client.Update(sendItems); err != nil {
		return fmt.Errorf("could not send updated items: %v", err)
	}
	if err := sync.syncRepo.DeleteAll(); err != nil {
		return fmt.Errorf("could not clear updated items: %v", err)
	}

	// get new/updated items
	ts, err := sync.syncRepo.LastUpdate()
	if err != nil {
		return fmt.Errorf("could not find timestamp of last update: %v", err)
	}
	recItems, err := sync.client.Updated([]item.Kind{item.KindTask}, ts)
	if err != nil {
		return fmt.Errorf("could not receive updates: %v", err)
	}

	updated := make([]item.Item, 0)
	for _, ri := range recItems {
		if ri.Deleted {
			if err := sync.localIDRepo.Delete(ri.ID); err != nil && !errors.Is(err, storage.ErrNotFound) {
				return fmt.Errorf("could not delete local id: %v", err)
			}
			if err := sync.taskRepo.Delete(ri.ID); err != nil && !errors.Is(err, storage.ErrNotFound) {
				return fmt.Errorf("could not delete task: %v", err)
			}
			continue
		}
		updated = append(updated, ri)
	}

	lidMap, err := sync.localIDRepo.FindAll()
	if err != nil {
		return fmt.Errorf("could not get local ids: %v", err)
	}
	for _, u := range updated {
		var tskBody item.TaskBody
		if err := json.Unmarshal([]byte(u.Body), &tskBody); err != nil {
			return fmt.Errorf("could not unmarshal task body: %v", err)
		}
		tsk := item.Task{
			ID:        u.ID,
			Date:      u.Date,
			Recurrer:  u.Recurrer,
			RecurNext: u.RecurNext,
			TaskBody:  tskBody,
		}
		if err := sync.taskRepo.Store(tsk); err != nil {
			return fmt.Errorf("could not store task: %v", err)
		}
		lid, ok := lidMap[u.ID]
		if !ok {
			lid, err = sync.localIDRepo.Next()
			if err != nil {
				return fmt.Errorf("could not get next local id: %v", err)
			}

			if err := sync.localIDRepo.Store(u.ID, lid); err != nil {
				return fmt.Errorf("could not store local id: %v", err)
			}
		}
	}

	return nil
}
