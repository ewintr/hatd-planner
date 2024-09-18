package main

import (
	"time"

	"github.com/google/uuid"
)

type Kind string

const (
	KindTask  Kind = "task"
	KindEvent Kind = "event"
)

var (
	KnownKinds = []Kind{KindTask, KindEvent}
)

type Item struct {
	ID      string    `json:"id"`
	Kind    Kind      `json:"kind"`
	Updated time.Time `json:"updated"`
	Deleted bool      `json:"deleted"`
	Body    string    `json:"body"`
}

func NewItem(k Kind, body string) Item {
	return Item{
		ID:      uuid.New().String(),
		Kind:    k,
		Updated: time.Now(),
		Body:    body,
	}
}
