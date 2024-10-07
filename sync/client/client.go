package client

import (
	"time"

	"go-mod.ewintr.nl/planner/item"
)

type Client interface {
	Update(items []item.Item) error
	Updated(ks []item.Kind, ts time.Time) ([]item.Item, error)
}
