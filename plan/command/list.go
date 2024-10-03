package command

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
	"go-mod.ewintr.nl/planner/plan/storage"
)

var ListCmd = &cli.Command{
	Name:  "list",
	Usage: "List everything",
}

func NewListCmd(localRepo storage.LocalID, eventRepo storage.Event) *cli.Command {
	ListCmd.Action = func(cCtx *cli.Context) error {
		return List(localRepo, eventRepo)
	}
	return ListCmd
}

func List(localRepo storage.LocalID, eventRepo storage.Event) error {
	localIDs, err := localRepo.FindAll()
	if err != nil {
		return fmt.Errorf("could not get local ids: %v", err)
	}
	all, err := eventRepo.FindAll()
	if err != nil {
		return err
	}
	for _, e := range all {
		lid, ok := localIDs[e.ID]
		if !ok {
			return fmt.Errorf("could not find local id for %s", e.ID)
		}
		fmt.Printf("%s\t%d\t%s\t%s\t%s\n", e.ID, lid, e.Title, e.Start.Format(time.DateTime), e.Duration.String())
	}

	return nil
}
