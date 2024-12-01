package command

import (
	"errors"
	"fmt"
	"slices"
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
	Value time.Time
}

func (ft *FlagDate) Set(val string) error {
	d, err := time.Parse(DateFormat, val)
	if err != nil {
		return fmt.Errorf("could not parse date: %v", d)
	}
	ft.Value = d

	return nil
}

func (ft *FlagDate) IsSet() bool {
	return !ft.Value.IsZero()
}

func (fs *FlagDate) Get() any {
	return fs.Value
}

type FlagTime struct {
	Name  string
	Value time.Time
}

func (ft *FlagTime) Set(val string) error {
	d, err := time.Parse(TimeFormat, val)
	if err != nil {
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

type FlagPeriod struct {
	Name  string
	Value item.RecurPeriod
}

func (fp *FlagPeriod) Set(val string) error {
	if !slices.Contains(item.ValidPeriods, item.RecurPeriod(val)) {
		return fmt.Errorf("not a valid period: %v", val)
	}
	fp.Value = item.RecurPeriod(val)
	return nil
}

func (fp *FlagPeriod) IsSet() bool {
	return fp.Value != ""
}

func (fp *FlagPeriod) Get() any {
	return fp.Value
}
