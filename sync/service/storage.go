package main

import (
	"errors"
	"time"

	"go-mod.ewintr.nl/planner/sync/item"
)

var (
	ErrNotFound = errors.New("not found")
)

type Syncer interface {
	Update(item item.Item) error
	Updated(kind []item.Kind, t time.Time) ([]item.Item, error)
}
