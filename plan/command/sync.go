package command

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage"
	"go-mod.ewintr.nl/planner/sync/client"
)

var SyncCmd = &cli.Command{
	Name:  "sync",
	Usage: "Synchronize with server",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "full",
			Aliases: []string{"f"},
			Usage:   "Force full sync",
		},
	},
}

func NewSyncCmd(client client.Client, syncRepo storage.Sync, localIDRepo storage.LocalID, eventRepo storage.Event) *cli.Command {
	SyncCmd.Action = func(cCtx *cli.Context) error {
		return Sync(client, syncRepo, localIDRepo, eventRepo, cCtx.Bool("full"))
	}
	return SyncCmd
}

func Sync(client client.Client, syncRepo storage.Sync, localIDRepo storage.LocalID, eventRepo storage.Event, full bool) error {
	// local new and updated
	sendItems, err := syncRepo.FindAll()
	if err != nil {
		return fmt.Errorf("could not get updated items: %v", err)
	}
	if err := client.Update(sendItems); err != nil {
		return fmt.Errorf("could not send updated items: %v", err)
	}
	if err := syncRepo.DeleteAll(); err != nil {
		return fmt.Errorf("could not clear updated items: %v", err)
	}

	// get new/updated items
	ts, err := syncRepo.LastUpdate()
	if err != nil {
		return fmt.Errorf("could not find timestamp of last update: %v", err)
	}
	recItems, err := client.Updated([]item.Kind{item.KindEvent}, ts)
	if err != nil {
		return fmt.Errorf("could not receive updates: %v", err)
	}

	updated := make([]item.Item, 0)
	for _, ri := range recItems {
		if ri.Deleted {
			if err := localIDRepo.Delete(ri.ID); err != nil {
				return fmt.Errorf("could not delete local id: %v", err)
			}
			if err := eventRepo.Delete(ri.ID); err != nil && !errors.Is(err, storage.ErrNotFound) {
				return fmt.Errorf("could not delete event: %v", err)
			}
			continue
		}
		updated = append(updated, ri)
	}

	lidMap, err := localIDRepo.FindAll()
	if err != nil {
		return fmt.Errorf("could not get local ids: %v", err)
	}
	for _, u := range updated {
		var eBody item.EventBody
		if err := json.Unmarshal([]byte(u.Body), &eBody); err != nil {
			return fmt.Errorf("could not unmarshal event body: %v", err)
		}
		e := item.Event{
			ID:        u.ID,
			EventBody: eBody,
		}
		if err := eventRepo.Store(e); err != nil {
			return fmt.Errorf("could not store event: %v", err)
		}
		lid, ok := lidMap[u.ID]
		if !ok {
			lid, err = localIDRepo.Next()
			if err != nil {
				return fmt.Errorf("could not get next local id: %v", err)
			}

			if err := localIDRepo.Store(u.ID, lid); err != nil {
				return fmt.Errorf("could not store local id: %v", err)
			}
		}
	}

	return nil
}
