package item

import (
	"slices"
	"time"
)

type RecurPeriod string

const (
	PeriodDay   RecurPeriod = "day"
	PeriodMonth RecurPeriod = "month"
)

var ValidPeriods = []RecurPeriod{PeriodDay, PeriodMonth}

type Recur struct {
	Start  time.Time   `json:"start"`
	Period RecurPeriod `json:"period"`
	Count  int         `json:"count"`
}

func (r *Recur) On(date time.Time) bool {
	switch r.Period {
	case PeriodDay:
		return r.onDays(date)
	case PeriodMonth:
		return r.onMonths(date)
	default:
		return false
	}
}

func (r *Recur) onDays(date time.Time) bool {
	if r.Start.After(date) {
		return false
	}

	testDate := r.Start
	for {
		if testDate.Equal(date) {
			return true
		}
		if testDate.After(date) {
			return false
		}

		dur := time.Duration(r.Count) * 24 * time.Hour
		testDate = testDate.Add(dur)
	}
}

func (r *Recur) onMonths(date time.Time) bool {
	if r.Start.After(date) {
		return false
	}

	tDate := r.Start
	for {
		if tDate.Equal(date) {
			return true
		}
		if tDate.After(date) {
			return false
		}

		y, m, d := tDate.Date()
		tDate = time.Date(y, m+time.Month(r.Count), d, 0, 0, 0, 0, time.UTC)
	}
}

func (r *Recur) NextAfter(old time.Time) time.Time {
	day, _ := time.ParseDuration("24h")
	test := old.Add(day)
	for {
		if r.On(test) || test.After(time.Date(2500, 1, 1, 0, 0, 0, 0, time.UTC)) {
			return test
		}
		test.Add(day)
	}
}

func (r *Recur) Valid() bool {
	return r.Start.IsZero() || !slices.Contains(ValidPeriods, r.Period)
}
