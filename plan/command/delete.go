package command

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"go-mod.ewintr.nl/planner/plan/storage"
)

var DeleteCmd = &cli.Command{
	Name:  "delete",
	Usage: "Delete an event",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:     "localID",
			Aliases:  []string{"l"},
			Usage:    "The local id of the event",
			Required: true,
		},
	},
}

func NewDeleteCmd(localRepo storage.LocalID, eventRepo storage.Event) *cli.Command {
	DeleteCmd.Action = func(cCtx *cli.Context) error {
		return Delete(localRepo, eventRepo, cCtx.Int("localID"))
	}
	return DeleteCmd
}

func Delete(localRepo storage.LocalID, eventRepo storage.Event, localID int) error {
	var id string
	idMap, err := localRepo.FindAll()
	if err != nil {
		return fmt.Errorf("could not get local ids: %v", err)
	}
	for eid, lid := range idMap {
		if localID == lid {
			id = eid
		}
	}
	if id == "" {
		return fmt.Errorf("could not find local id")
	}

	if err := eventRepo.Delete(id); err != nil {
		return fmt.Errorf("could not delete event: %v", err)
	}

	return nil
}
