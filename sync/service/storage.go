package main

import (
	"errors"
	"time"

	"go-mod.ewintr.nl/planner/item"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrNotARecurrer = errors.New("not a recurrer")
)

type Syncer interface {
	Update(item item.Item, t time.Time) error
	Updated(kind []item.Kind, t time.Time) ([]item.Item, error)
}

type Recurrer interface {
	RecursBefore(date time.Time) ([]item.Item, error)
	RecursNext(id string, date time.Time, t time.Time) error
}
