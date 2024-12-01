package main

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
	"go-mod.ewintr.nl/planner/item"
)

type Recur struct {
	repoSync  Syncer
	repoRecur Recurrer
	logger    *slog.Logger
}

func NewRecur(repoRecur Recurrer, repoSync Syncer, logger *slog.Logger) *Recur {
	r := &Recur{
		repoRecur: repoRecur,
		repoSync:  repoSync,
		logger:    logger,
	}

	return r
}

func (r *Recur) Run(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		if err := r.Recur(); err != nil {
			r.logger.Error("could not recur", "error", err)
		}
	}
}

func (r *Recur) Recur() error {
	items, err := r.repoRecur.RecursBefore(time.Now())
	if err != nil {
		return err
	}
	for _, i := range items {
		// spawn instance
		ne, err := item.NewEvent(i)
		if err != nil {
			return err
		}
		y, m, d := i.RecurNext.Date()
		ne.ID = uuid.New().String()
		ne.Recurrer = nil
		ne.RecurNext = time.Time{}
		ne.Start = time.Date(y, m, d, ne.Start.Hour(), ne.Start.Minute(), 0, 0, time.UTC)

		ni, err := ne.Item()
		if err != nil {
			return err
		}
		if err := r.repoSync.Update(ni, time.Now()); err != nil {
			return err
		}

		// set next
		if err := r.repoRecur.RecursNext(i.ID, i.Recurrer.NextAfter(i.RecurNext), time.Now()); err != nil {
			return err
		}
	}
	r.logger.Info("processed recurring items", "count", len(items))

	return nil
}
