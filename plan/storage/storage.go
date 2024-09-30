package storage

import "go-mod.ewintr.nl/planner/item"

type EventRepo interface {
	Store(event item.Event) error
	Find(id string) (item.Event, error)
	FindAll() ([]item.Event, error)
	Delete(id string) error
}
