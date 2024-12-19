package item_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go-mod.ewintr.nl/planner/item"
)

func TestDaily(t *testing.T) {
	t.Parallel()

	daily := item.Daily{
		Start: item.NewDate(2021, 1, 31), // a sunday
	}
	dailyStr := "2021-01-31, daily"

	t.Run("parse", func(t *testing.T) {
		if diff := cmp.Diff(daily, item.NewRecurrer(dailyStr)); diff != "" {
			t.Errorf("(-exp +got):\n%s", diff)
		}
	})

	t.Run("string", func(t *testing.T) {
		if dailyStr != daily.String() {
			t.Errorf("exp %v, got %v", dailyStr, daily.String())
		}
	})

	t.Run("recurs_on", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			date item.Date
			exp  bool
		}{
			{
				name: "before",
				date: item.NewDate(2021, 1, 30),
			},
			{
				name: "on",
				date: daily.Start,
				exp:  true,
			},
			{
				name: "after",
				date: item.NewDate(2021, 2, 1),
				exp:  true,
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				if tc.exp != daily.RecursOn(tc.date) {
					t.Errorf("exp %v, got %v", tc.exp, daily.RecursOn(tc.date))
				}
			})
		}
	})
}

func TestEveryNDays(t *testing.T) {
	t.Parallel()

	every := item.EveryNDays{
		Start: item.NewDate(2022, 6, 8),
		N:     5,
	}
	everyStr := "2022-06-08, every 5 days"

	t.Run("parse", func(t *testing.T) {
		if diff := cmp.Diff(every, item.NewRecurrer(everyStr)); diff != "" {
			t.Errorf("(-exp +got):\n%s", diff)
		}
	})

	t.Run("string", func(t *testing.T) {
		if everyStr != every.String() {
			t.Errorf("exp %v, got %v", everyStr, every.String())
		}
	})

	t.Run("recurs on", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			date item.Date
			exp  bool
		}{
			{
				name: "before",
				date: item.NewDate(2022, 1, 1),
			},
			{
				name: "start",
				date: every.Start,
				exp:  true,
			},
			{
				name: "after true",
				date: every.Start.Add(15),
				exp:  true,
			},
			{
				name: "after false",
				date: every.Start.Add(16),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				if tc.exp != every.RecursOn(tc.date) {
					t.Errorf("exp %v, got %v", tc.exp, tc.date)
				}
			})
		}
	})
}

func TestParseWeekly(t *testing.T) {
	t.Parallel()

	start := item.NewDate(2021, 2, 7)
	for _, tc := range []struct {
		name      string
		input     []string
		expOK     bool
		expWeekly item.Weekly
	}{
		{
			name: "empty",
		},
		{
			name:  "wrong type",
			input: []string{"daily"},
		},
		{
			name:  "wrong count",
			input: []string{"weeekly"},
		},
		{
			name:  "unknown day",
			input: []string{"weekly", "festivus"},
		},
		{
			name:  "one day",
			input: []string{"weekly", "monday"},
			expOK: true,
			expWeekly: item.Weekly{
				Start: start,
				Weekdays: item.Weekdays{
					time.Monday,
				},
			},
		},
		{
			name:  "multiple days",
			input: []string{"weekly", "monday & thursday & saturday"},
			expOK: true,
			expWeekly: item.Weekly{
				Start: start,
				Weekdays: item.Weekdays{
					time.Monday,
					time.Thursday,
					time.Saturday,
				},
			},
		},
		{
			name:  "wrong order",
			input: []string{"weekly", "sunday & thursday & wednesday"},
			expOK: true,
			expWeekly: item.Weekly{
				Start: start,
				Weekdays: item.Weekdays{
					time.Wednesday,
					time.Thursday,
					time.Sunday,
				},
			},
		},
		{
			name:  "doubles",
			input: []string{"weekly", "sunday & sunday & monday"},
			expOK: true,
			expWeekly: item.Weekly{
				Start: start,
				Weekdays: item.Weekdays{
					time.Monday,
					time.Sunday,
				},
			},
		},
		{
			name:  "one unknown",
			input: []string{"weekly", "sunday & someday"},
			expOK: true,
			expWeekly: item.Weekly{
				Start: start,
				Weekdays: item.Weekdays{
					time.Sunday,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actWeekly, actOK := item.ParseWeekly(start, tc.input)
			if tc.expOK != actOK {
				t.Errorf("exp %v, got %v", tc.expOK, actOK)
			}
			if !tc.expOK {
				return
			}
			if diff := cmp.Diff(tc.expWeekly, actWeekly); diff != "" {
				t.Errorf("(-exp, +got)%s\n", diff)
			}
		})
	}
}

func TestWeekly(t *testing.T) {
	t.Parallel()

	weekly := item.Weekly{
		Start: item.NewDate(2021, 1, 31), // a sunday
		Weekdays: item.Weekdays{
			time.Monday,
			time.Wednesday,
			time.Thursday,
		},
	}
	weeklyStr := "2021-01-31, weekly, monday & wednesday & thursday"

	t.Run("parse", func(t *testing.T) {
		if diff := cmp.Diff(weekly, item.NewRecurrer(weeklyStr)); diff != "" {
			t.Errorf("(-exp, +got)%s\n", diff)
		}
	})

	t.Run("string", func(t *testing.T) {
		if weeklyStr != weekly.String() {
			t.Errorf("exp %v, got %v", weeklyStr, weekly.String())
		}
	})

	t.Run("recurs_on", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			date item.Date
			exp  bool
		}{
			{
				name: "before start",
				date: item.NewDate(2021, 1, 27), // a wednesday
			},
			{
				name: "right weekday",
				date: item.NewDate(2021, 2, 1), // a monday
				exp:  true,
			},
			{
				name: "another right day",
				date: item.NewDate(2021, 2, 3), // a wednesday
				exp:  true,
			},
			{
				name: "wrong weekday",
				date: item.NewDate(2021, 2, 5), // a friday
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				if tc.exp != weekly.RecursOn(tc.date) {
					t.Errorf("exp %v, got %v", tc.exp, weekly.RecursOn(tc.date))
				}
			})
		}
	})
}

func TestEveryNWeeks(t *testing.T) {
	t.Parallel()

	everyNWeeks := item.EveryNWeeks{
		Start: item.NewDate(2021, 2, 3),
		N:     3,
	}
	everyNWeeksStr := "2021-02-03, every 3 weeks"

	t.Run("parse", func(t *testing.T) {
		if everyNWeeks != item.NewRecurrer(everyNWeeksStr) {
			t.Errorf("exp %v, got %v", everyNWeeks, item.NewRecurrer(everyNWeeksStr))
		}
	})

	t.Run("string", func(t *testing.T) {
		if everyNWeeksStr != everyNWeeks.String() {
			t.Errorf("exp %v, got %v", everyNWeeksStr, everyNWeeks.String())
		}
	})

	t.Run("recurs on", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			date item.Date
			exp  bool
		}{
			{
				name: "before start",
				date: item.NewDate(2021, 1, 27),
			},
			{
				name: "on start",
				date: item.NewDate(2021, 2, 3),
				exp:  true,
			},
			{
				name: "wrong day",
				date: item.NewDate(2021, 2, 4),
			},
			{
				name: "one week after",
				date: item.NewDate(2021, 2, 10),
			},
			{
				name: "first interval",
				date: item.NewDate(2021, 2, 24),
				exp:  true,
			},
			{
				name: "second interval",
				date: item.NewDate(2021, 3, 17),
				exp:  true,
			},
			{
				name: "second interval plus one week",
				date: item.NewDate(2021, 3, 24),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				if tc.exp != everyNWeeks.RecursOn(tc.date) {
					t.Errorf("exp %v, got %v", tc.exp, everyNWeeks.RecursOn(tc.date))
				}
			})
		}
	})
}

func TestEveryNMonths(t *testing.T) {
	everyNMonths := item.EveryNMonths{
		Start: item.NewDate(2021, 2, 3),
		N:     3,
	}
	everyNMonthsStr := "2021-02-03, every 3 months"

	t.Run("parse", func(t *testing.T) {
		if diff := cmp.Diff(everyNMonths, item.NewRecurrer(everyNMonthsStr)); diff != "" {
			t.Errorf("(-exp, +got)%s\n", diff)
		}
	})

	t.Run("string", func(t *testing.T) {
		if everyNMonthsStr != everyNMonths.String() {
			t.Errorf("exp %v, got %v", everyNMonthsStr, everyNMonths.String())
		}
	})

	t.Run("recurs on", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			date item.Date
			exp  bool
		}{
			{
				name: "before start",
				date: item.NewDate(2021, 1, 27),
			},
			{
				name: "on start",
				date: item.NewDate(2021, 2, 3),
				exp:  true,
			},
			{
				name: "8 weeks after",
				date: item.NewDate(2021, 3, 31),
			},
			{
				name: "one month",
				date: item.NewDate(2021, 3, 3),
			},
			{
				name: "3 months",
				date: item.NewDate(2021, 5, 3),
				exp:  true,
			},
			{
				name: "4 months",
				date: item.NewDate(2021, 6, 3),
			},
			{
				name: "6 months",
				date: item.NewDate(2021, 8, 3),
				exp:  true,
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				if tc.exp != everyNMonths.RecursOn(tc.date) {
					t.Errorf("exp %v, got %v", tc.exp, everyNMonths.RecursOn(tc.date))
				}
			})
		}
	})

	t.Run("recurs every year", func(t *testing.T) {
		recur := item.EveryNMonths{
			Start: item.NewDate(2021, 3, 1),
			N:     12,
		}
		if recur.RecursOn(item.NewDate(2021, 3, 9)) {
			t.Errorf("exp false, got true")
		}
	})

	t.Run("bug", func(t *testing.T) {
		recur := item.EveryNMonths{
			Start: item.NewDate(2021, 3, 1),
			N:     1,
		}
		if recur.RecursOn(item.NewDate(2021, 11, 3)) {
			t.Errorf("exp false, got true")
		}
	})
}
