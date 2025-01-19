package command

import (
	"errors"

	"go-mod.ewintr.nl/planner/plan/storage"
	"go-mod.ewintr.nl/planner/sync/client"
)

const (
	DateFormat = "2006-01-02"
	TimeFormat = "15:04"
)

var (
	ErrWrongCommand = errors.New("wrong command")
	ErrInvalidArg   = errors.New("invalid argument")
)

type Repositories interface {
	Begin() (*storage.Tx, error)
	LocalID(tx *storage.Tx) storage.LocalID
	Sync(tx *storage.Tx) storage.Sync
	Task(tx *storage.Tx) storage.Task
	Schedule(tx *storage.Tx) storage.Schedule
}

type CommandArgs interface {
	Parse(main []string, fields map[string]string) (Command, error)
}

type Command interface {
	Do(repos Repositories, client client.Client) (CommandResult, error)
}

type CommandResult interface {
	Render() string
}
