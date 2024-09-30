package command

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage"
)

var AddCmd = &cli.Command{
	Name:  "add",
	Usage: "Add a new event",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Usage:   "The event that will happen",
		},
		&cli.StringFlag{
			Name:    "on",
			Aliases: []string{"o"},
			Usage:   "The date, in YYYY-MM-DD format",
		},
		&cli.StringFlag{
			Name:    "at",
			Aliases: []string{"a"},
			Usage:   "The time, in HH:MM format. If omitted, the event will last the whole day",
		},
		&cli.StringFlag{
			Name:    "for",
			Aliases: []string{"f"},
			Usage:   "The duration, in show format (e.g. 1h30m)",
		},
	},
}

func NewAddCmd(repo storage.EventRepo) *cli.Command {
	AddCmd.Action = NewAddAction(repo)
	return AddCmd
}

func NewAddAction(repo storage.EventRepo) func(*cli.Context) error {
	return func(cCtx *cli.Context) error {
		desc := cCtx.String("name")
		date, err := time.Parse("2006-01-02", cCtx.String("date"))
		if err != nil {
			return fmt.Errorf("could not parse date: %v", err)
		}

		one := item.Event{
			ID: "a",
			EventBody: item.EventBody{
				Title: desc,
				Start: date,
			},
		}
		if err := repo.Store(one); err != nil {
			return fmt.Errorf("could not store event: %v", err)
		}

		return nil
	}
}
