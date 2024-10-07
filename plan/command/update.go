package command

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
	"go-mod.ewintr.nl/planner/plan/storage"
)

var UpdateCmd = &cli.Command{
	Name:  "update",
	Usage: "Update an event",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:     "localID",
			Aliases:  []string{"l"},
			Usage:    "The local id of the event",
			Required: true,
		},
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

func NewUpdateCmd(localRepo storage.LocalID, eventRepo storage.Event, syncRepo storage.Sync) *cli.Command {
	UpdateCmd.Action = func(cCtx *cli.Context) error {
		return Update(localRepo, eventRepo, syncRepo, cCtx.Int("localID"), cCtx.String("name"), cCtx.String("on"), cCtx.String("at"), cCtx.String("for"))
	}
	return UpdateCmd
}

func Update(localRepo storage.LocalID, eventRepo storage.Event, syncRepo storage.Sync, localID int, nameStr, onStr, atStr, frStr string) error {
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

	e, err := eventRepo.Find(id)
	if err != nil {
		return fmt.Errorf("could not find event")
	}

	if nameStr != "" {
		e.Title = nameStr
	}
	if onStr != "" || atStr != "" {
		oldStart := e.Start
		dateStr := oldStart.Format("2006-01-02")
		if onStr != "" {
			dateStr = onStr
		}
		timeStr := oldStart.Format("15:04")
		if atStr != "" {
			timeStr = atStr
		}
		newStart, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", dateStr, timeStr))
		if err != nil {
			return fmt.Errorf("could not parse new start: %v", err)
		}
		e.Start = newStart
	}

	if frStr != "" { // no check on at, can set a duration with at 00:00, making it not a whole day
		fr, err := time.ParseDuration(frStr)
		if err != nil {
			return fmt.Errorf("%w: could not parse duration: %s", ErrInvalidArg, err)
		}
		e.Duration = fr
	}
	if err := eventRepo.Store(e); err != nil {
		return fmt.Errorf("could not store event: %v", err)
	}

	it, err := e.Item()
	if err != nil {
		return fmt.Errorf("could not convert event to sync item: %v", err)
	}
	if err := syncRepo.Store(it); err != nil {
		return fmt.Errorf("could not store sync item: %v", err)
	}

	return nil
}
