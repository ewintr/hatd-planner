package command

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"go-mod.ewintr.nl/planner/item"
)

const (
	DateFormat = "2006-01-02"
	TimeFormat = "15:04"
)

var (
	ErrWrongCommand = errors.New("wrong command")
	ErrInvalidArg   = errors.New("invalid argument")
)

type Flag interface {
	Set(val string) error
	IsSet() bool
	Get() any
}

type FlagString struct {
	Name  string
	Value string
}

func (fs *FlagString) Set(val string) error {
	fs.Value = val
	return nil
}

func (fs *FlagString) IsSet() bool {
	return fs.Value != ""
}

func (fs *FlagString) Get() any {
	return fs.Value
}

type FlagDate struct {
	Name  string
	Value item.Date
}

func (fd *FlagDate) Set(val string) error {
	d := item.NewDateFromString(val)
	if d.IsZero() {
		return fmt.Errorf("could not parse date: %v", d)
	}
	fd.Value = d

	return nil
}

func (fd *FlagDate) IsSet() bool {
	return !fd.Value.IsZero()
}

func (fd *FlagDate) Get() any {
	return fd.Value
}

type FlagTime struct {
	Name  string
	Value item.Time
}

func (ft *FlagTime) Set(val string) error {
	d := item.NewTimeFromString(val)
	if d.IsZero() {
		return fmt.Errorf("could not parse date: %v", d)
	}
	ft.Value = d

	return nil
}

func (fd *FlagTime) IsSet() bool {
	return !fd.Value.IsZero()
}

func (fs *FlagTime) Get() any {
	return fs.Value
}

type FlagDuration struct {
	Name  string
	Value time.Duration
}

func (fd *FlagDuration) Set(val string) error {
	dur, err := time.ParseDuration(val)
	if err != nil {
		return fmt.Errorf("could not parse duration: %v", err)
	}
	fd.Value = dur
	return nil
}

func (fd *FlagDuration) IsSet() bool {
	return fd.Value.String() != "0s"
}

func (fs *FlagDuration) Get() any {
	return fs.Value
}

type FlagRecurrer struct {
	Name  string
	Value item.Recurrer
}

func (fr *FlagRecurrer) Set(val string) error {
	fr.Value = item.NewRecurrer(val)
	if fr.Value == nil {
		return fmt.Errorf("not a valid recurrer: %v", val)
	}
	return nil
}

func (fr *FlagRecurrer) IsSet() bool {
	return fr.Value != nil
}

func (fr *FlagRecurrer) Get() any {
	return fr.Value
}

type FlagInt struct {
	Name  string
	Value int
}

func (fi *FlagInt) Set(val string) error {
	i, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("not a valid integer: %v", val)
	}

	fi.Value = i
	return nil
}

func (fi *FlagInt) IsSet() bool {
	return fi.Value != 0
}

func (fi *FlagInt) Get() any {
	return fi.Value
}
