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
	r.logger.Info("start looking for recurring items")
	today := item.NewDateFromString(time.Now().Format(item.DateFormat))
	items, err := r.repoRecur.ShouldRecur(today)
	if err != nil {
		return err
	}
	r.logger.Info("found recurring items", "count", len(items))
	for _, i := range items {
		r.logger.Info("processing recurring item", "id", i.ID)
		// spawn instance
		newItem := i
		newItem.ID = uuid.New().String()
		newItem.Date = i.RecurNext
		newItem.Recurrer = nil
		newItem.RecurNext = item.Date{}
		if err := r.repoSync.Update(newItem, time.Now()); err != nil {
			return err
		}
		r.logger.Info("spawned instance", "newID", newItem.ID, "date", newItem.Date)

		// update recurrer
		i.RecurNext = item.FirstRecurAfter(i.Recurrer, i.RecurNext)
		if err := r.repoSync.Update(i, time.Now()); err != nil {
			return err
		}
		r.logger.Info("recurring item processed", "id", i.ID, "recurNext", i.RecurNext.String())
	}
	r.logger.Info("processed recurring items", "count", len(items))

	return nil
}
