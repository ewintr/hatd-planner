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
	days      int
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

func (r *Recur) Run(days int, interval time.Duration) {
	r.days = days
	ticker := time.NewTicker(interval)

	for range ticker.C {
		today := item.NewDateFromString(time.Now().Format(item.DateFormat))
		last := today.Add(r.days)
		if err := r.Recur(last); err != nil {
			r.logger.Error("could not recur", "error", err)
		}
	}
}

func (r *Recur) Recur(until item.Date) error {
	r.logger.Info("start looking for recurring items", "until", until.String())

	items, err := r.repoRecur.ShouldRecur(until)
	if err != nil {
		return err
	}

	r.logger.Info("found recurring items", "count", len(items))
	for _, i := range items {
		r.logger.Info("processing recurring item", "id", i.ID)
		newRecurNext := item.FirstRecurAfter(i.Recurrer, i.RecurNext)

		for {
			// spawn instance
			newItem := i
			newItem.ID = uuid.New().String()
			newItem.Date = newRecurNext
			newItem.Recurrer = nil
			newItem.RecurNext = item.Date{}
			if err := r.repoSync.Update(newItem, time.Now()); err != nil {
				return err
			}
			r.logger.Info("spawned instance", "newID", newItem.ID, "date", newItem.Date)

			newRecurNext = item.FirstRecurAfter(i.Recurrer, newRecurNext)

			if newRecurNext.After(until) {
				break
			}
		}

		// update recurrer
		i.RecurNext = newRecurNext
		if err := r.repoSync.Update(i, time.Now()); err != nil {
			return err
		}
		r.logger.Info("recurring item processed", "id", i.ID, "recurNext", i.RecurNext.String())
	}
	r.logger.Info("processed recurring items", "count", len(items))

	return nil
}
