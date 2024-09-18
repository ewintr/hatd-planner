package main

import (
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
)

type Syncer interface {
	Update(item Item) error
	Updated(kind []Kind, t time.Time) ([]Item, error)
}
