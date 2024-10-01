package command

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage"
)

var (
	ErrInvalidArg = errors.New("invalid argument")
)

var AddCmd = &cli.Command{
	Name:  "add",
	Usage: "Add a new event",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "The event that will happen",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "on",
			Aliases:  []string{"o"},
			Usage:    "The date, in YYYY-MM-DD format",
			Required: true,
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
	AddCmd.Action = func(cCtx *cli.Context) error {
		return Add(cCtx.String("name"), cCtx.String("on"), cCtx.String("at"), cCtx.String("for"), repo)
	}
	return AddCmd
}

func Add(nameStr, onStr, atStr, frStr string, repo storage.EventRepo) error {
	if nameStr == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidArg)
	}
	if onStr == "" {
		return fmt.Errorf("%w: date is required", ErrInvalidArg)
	}
	if atStr == "" && frStr != "" {
		return fmt.Errorf("%w: can not have duration without start time", ErrInvalidArg)
	}
	if atStr == "" && frStr == "" {
		frStr = "24h"
	}

	startFormat := "2006-01-02"
	startStr := onStr
	if atStr != "" {
		startFormat = fmt.Sprintf("%s 15:04", startFormat)
		startStr = fmt.Sprintf("%s %s", startStr, atStr)
	}
	start, err := time.Parse(startFormat, startStr)
	if err != nil {
		return fmt.Errorf("%w: could not parse start time and date: %v", ErrInvalidArg, err)
	}

	e := item.Event{
		ID: uuid.New().String(),
		EventBody: item.EventBody{
			Title: nameStr,
			Start: start,
		},
	}

	if frStr != "" {
		fr, err := time.ParseDuration(frStr)
		if err != nil {
			return fmt.Errorf("%w: could not parse time: %s", ErrInvalidArg, err)
		}
		e.Duration = fr
	}
	if err := repo.Store(e); err != nil {
		return fmt.Errorf("could not store event: %v", err)
	}

	return nil
}
