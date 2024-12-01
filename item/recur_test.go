package item_test

import (
	"testing"
	"time"

	"go-mod.ewintr.nl/planner/item"
)

func TestRecur(t *testing.T) {
	t.Parallel()

	t.Run("days", func(t *testing.T) {
		r := item.Recur{
			Start:  time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			Period: item.PeriodDay,
			Count:  5,
		}
		day := 24 * time.Hour

		for _, tc := range []struct {
			name string
			date time.Time
			exp  bool
		}{
			{
				name: "before",
				date: time.Date(202, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				name: "start",
				date: r.Start,
				exp:  true,
			},
			{
				name: "after true",
				date: r.Start.Add(15 * day),
				exp:  true,
			},
			{
				name: "after false",
				date: r.Start.Add(16 * day),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				if act := r.On(tc.date); tc.exp != act {
					t.Errorf("exp %v, got %v", tc.exp, act)
				}
			})
		}
	})

	t.Run("months", func(t *testing.T) {
		r := item.Recur{
			Start:  time.Date(2021, 2, 3, 0, 0, 0, 0, time.UTC),
			Period: item.PeriodMonth,
			Count:  3,
		}

		for _, tc := range []struct {
			name string
			date time.Time
			exp  bool
		}{
			{
				name: "before start",
				date: time.Date(2021, 1, 27, 0, 0, 0, 0, time.UTC),
			},
			{
				name: "on start",
				date: time.Date(2021, 2, 3, 0, 0, 0, 0, time.UTC),
				exp:  true,
			},
			{
				name: "8 weeks after",
				date: time.Date(2021, 3, 31, 0, 0, 0, 0, time.UTC),
			},
			{
				name: "one month",
				date: time.Date(2021, 3, 3, 0, 0, 0, 0, time.UTC),
			},
			{
				name: "3 months",
				date: time.Date(2021, 5, 3, 0, 0, 0, 0, time.UTC),
				exp:  true,
			},
			{
				name: "4 months",
				date: time.Date(2021, 6, 3, 0, 0, 0, 0, time.UTC),
			},
			{
				name: "6 months",
				date: time.Date(2021, 8, 3, 0, 0, 0, 0, time.UTC),
				exp:  true,
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				if act := r.On(tc.date); tc.exp != act {
					t.Errorf("exp %v, got %v", tc.exp, act)
				}
			})
		}
	})

}
