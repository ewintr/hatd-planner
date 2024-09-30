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

func NewListCmd(repo storage.EventRepo) *cli.Command {
	ListCmd.Action = NewListAction(repo)
	return ListCmd
}

func NewListAction(repo storage.EventRepo) func(*cli.Context) error {
	return func(cCtx *cli.Context) error {
		all, err := repo.FindAll()
		if err != nil {
			return err
		}
		for _, e := range all {
			fmt.Printf("%s\t%s\t%s\t%s\n", e.ID, e.Title, e.Start.Format(time.DateTime), e.Duration.String())
		}

		return nil
	}

}
