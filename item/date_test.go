package item_test

import (
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go-mod.ewintr.nl/planner/item"
)

func TestWeekdaysSort(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input item.Weekdays
		exp   item.Weekdays
	}{
		{
			name: "empty",
		},
		{
			name:  "one",
			input: item.Weekdays{time.Tuesday},
			exp:   item.Weekdays{time.Tuesday},
		},
		{
			name:  "multiple",
			input: item.Weekdays{time.Wednesday, time.Tuesday, time.Monday},
			exp:   item.Weekdays{time.Monday, time.Tuesday, time.Wednesday},
		},
		{
			name:  "sunday is last",
			input: item.Weekdays{time.Saturday, time.Sunday, time.Monday},
			exp:   item.Weekdays{time.Monday, time.Saturday, time.Sunday},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sort.Sort(tc.input)
			if diff := cmp.Diff(tc.exp, tc.input); diff != "" {
				t.Errorf("(-exp, +got)%s\n", diff)
			}
		})
	}
}

func TestWeekdaysUnique(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input item.Weekdays
		exp   item.Weekdays
	}{
		{
			name:  "empty",
			input: item.Weekdays{},
			exp:   item.Weekdays{},
		},
		{
			name:  "single",
			input: item.Weekdays{time.Monday},
			exp:   item.Weekdays{time.Monday},
		},
		{
			name:  "no doubles",
			input: item.Weekdays{time.Monday, time.Tuesday, time.Wednesday},
			exp:   item.Weekdays{time.Monday, time.Tuesday, time.Wednesday},
		},
		{
			name:  "doubles",
			input: item.Weekdays{time.Monday, time.Monday, time.Wednesday, time.Monday},
			exp:   item.Weekdays{time.Monday, time.Wednesday},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if diff := cmp.Diff(tc.exp, tc.input.Unique()); diff != "" {
				t.Errorf("(-exp, +got)%s\n", diff)
			}
		})
	}
}

func TestNewDateFromString(t *testing.T) {
	t.Parallel()

	t.Run("simple", func(t *testing.T) {
		for _, tc := range []struct {
			name  string
			input string
			exp   item.Date
		}{
			{
				name: "empty",
				exp:  item.Date{},
			},
			{
				name:  "no date",
				input: "no date",
				exp:   item.Date{},
			},
			{
				name:  "short",
				input: "2021-01-30",
				exp:   item.NewDate(2021, 1, 30),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				if diff := cmp.Diff(tc.exp, item.NewDateFromString(tc.input)); diff != "" {
					t.Errorf("(-exp, +got)%s\n", diff)
				}
			})
		}
	})

	t.Run("day name", func(t *testing.T) {
		monday := item.Today().Add(1)
		for {
			if monday.Weekday() == time.Monday {
				break
			}
			monday = monday.Add(1)
		}
		for _, tc := range []struct {
			name  string
			input string
		}{
			{
				name:  "dayname lowercase",
				input: "monday",
			},
			{
				name:  "dayname capitalized",
				input: "Monday",
			},
			{
				name:  "dayname short",
				input: "mon",
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				if diff := cmp.Diff(monday, item.NewDateFromString(tc.input)); diff != "" {
					t.Errorf("(-exp, +got)%s\n", diff)
				}
			})
		}
	})

	t.Run("relative days", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			exp  item.Date
		}{
			{
				name: "today",
				exp:  item.Today(),
			},
			{
				name: "tod",
				exp:  item.Today(),
			},
			{
				name: "tomorrow",
				exp:  item.Today().Add(1),
			},
			{
				name: "tom",
				exp:  item.Today().Add(1),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				if diff := cmp.Diff(tc.exp, item.NewDateFromString(tc.name)); diff != "" {
					t.Errorf("(-exp, +got)%s\n", diff)
				}
			})
		}
	})
}

func TestDateDaysBetween(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		d1   item.Date
		d2   item.Date
		exp  int
	}{
		{
			name: "same",
			d1:   item.NewDate(2021, 6, 23),
			d2:   item.NewDate(2021, 6, 23),
		},
		{
			name: "one",
			d1:   item.NewDate(2021, 6, 23),
			d2:   item.NewDate(2021, 6, 24),
			exp:  1,
		},
		{
			name: "many",
			d1:   item.NewDate(2021, 6, 23),
			d2:   item.NewDate(2024, 3, 7),
			exp:  988,
		},
		{
			name: "edge",
			d1:   item.NewDate(2020, 12, 30),
			d2:   item.NewDate(2021, 1, 3),
			exp:  4,
		},
		{
			name: "reverse",
			d1:   item.NewDate(2021, 6, 23),
			d2:   item.NewDate(2021, 5, 23),
			exp:  31,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.exp != tc.d1.DaysBetween(tc.d2) {
				t.Errorf("exp %v, got %v", tc.exp, tc.d1.DaysBetween(tc.d2))
			}
		})
	}
}

func TestDateString(t *testing.T) {
	for _, tc := range []struct {
		name string
		date item.Date
		exp  string
	}{
		{
			name: "zero",
			date: item.NewDate(0, 0, 0),
			exp:  "",
		},
		{
			name: "normal",
			date: item.NewDate(2021, 5, 30),
			exp:  "2021-05-30",
		},
		{
			name: "normalize",
			date: item.NewDate(2021, 5, 32),
			exp:  "2021-06-01",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.exp != tc.date.String() {
				t.Errorf("exp %v, got %v", tc.exp, tc.date.String())
			}
		})
	}
}

// func TestDateHuman(t *testing.T) {
// 	for _, tc := range []struct {
// 		name string
// 		date task.Date
// 		exp  string
// 	}{
// 		{
// 			name: "zero",
// 			date: task.NewDate(0, 0, 0),
// 			exp:  "-",
// 		},
// 		{
// 			name: "default",
// 			date: task.NewDate(2020, 1, 1),
// 			exp:  "2020-01-01 (wednesday)",
// 		},
// 		{
// 			name: "today",
// 			date: task.Today(),
// 			exp:  "today",
// 		},
// 		{
// 			name: "tomorrow",
// 			date: task.Today().Add(1),
// 			exp:  "tomorrow",
// 		},
// 	} {
// 		t.Run(tc.name, func(t *testing.T) {
// 			test.Equals(t, tc.exp, tc.date.Human())
// 		})
// 	}
// }

func TestDateIsZero(t *testing.T) {
	t.Parallel()

	if !(item.Date{}.IsZero()) {
		t.Errorf("exp true, got false")
	}
	if item.NewDate(2021, 6, 24).IsZero() {
		t.Errorf("exp false, got true")
	}
}

func TestDateAfter(t *testing.T) {
	t.Parallel()

	day := item.NewDate(2021, 1, 31)
	for _, tc := range []struct {
		name string
		tDay item.Date
		exp  bool
	}{
		{
			name: "after",
			tDay: item.NewDate(2021, 1, 30),
			exp:  true,
		},
		{
			name: "on",
			tDay: day,
		},
		{
			name: "before",
			tDay: item.NewDate(2021, 2, 1),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if diff := cmp.Diff(tc.exp, day.After(tc.tDay)); diff != "" {
				t.Errorf("(-exp, +got)%s\n", diff)
			}
		})
	}
}
