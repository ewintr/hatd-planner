package item

import (
	"fmt"
	"strconv"
	"strings"
)

type Recurrer interface {
	RecursOn(date Date) bool
	First() Date
	String() string
}

func NewRecurrer(recurStr string) Recurrer {
	terms := strings.Split(recurStr, ",")
	if len(terms) < 2 {
		return nil
	}

	start := NewDateFromString(terms[0])
	if start.IsZero() {
		return nil
	}

	terms = terms[1:]
	for i, t := range terms {
		terms[i] = strings.TrimSpace(t)
	}

	for _, parseFunc := range []func(Date, []string) (Recurrer, bool){
		ParseDaily, ParseEveryNDays, ParseWeekly,
		ParseEveryNWeeks, ParseEveryNMonths,
	} {
		if recur, ok := parseFunc(start, terms); ok {
			return recur
		}
	}

	return nil
}

func FirstRecurAfter(r Recurrer, d Date) Date {
	lim := NewDate(2050, 1, 1)
	for {
		d = d.Add(1)
		if r.RecursOn(d) || d.Equal(lim) {
			return d
		}
	}
}

type Daily struct {
	Start Date
}

// yyyy-mm-dd, daily
func ParseDaily(start Date, terms []string) (Recurrer, bool) {
	if len(terms) < 1 {
		return nil, false
	}

	if terms[0] != "daily" {
		return nil, false
	}

	return Daily{
		Start: start,
	}, true
}

func (d Daily) RecursOn(date Date) bool {
	return date.Equal(d.Start) || date.After(d.Start)
}

func (d Daily) First() Date { return FirstRecurAfter(d, d.Start.Add(-1)) }

func (d Daily) String() string {
	return fmt.Sprintf("%s, daily", d.Start.String())
}

type EveryNDays struct {
	Start Date
	N     int
}

// yyyy-mm-dd, every 3 days
func ParseEveryNDays(start Date, terms []string) (Recurrer, bool) {
	if len(terms) != 1 {
		return EveryNDays{}, false
	}

	terms = strings.Split(terms[0], " ")
	if len(terms) != 3 || terms[0] != "every" || terms[2] != "days" {
		return EveryNDays{}, false
	}

	n, err := strconv.Atoi(terms[1])
	if err != nil {
		return EveryNDays{}, false
	}

	return EveryNDays{
		Start: start,
		N:     n,
	}, true
}

func (nd EveryNDays) RecursOn(date Date) bool {
	if nd.Start.After(date) {
		return false
	}

	testDate := nd.Start
	for {
		switch {
		case testDate.Equal(date):
			return true
		case testDate.After(date):
			return false
		default:
			testDate = testDate.Add(nd.N)
		}
	}
}

func (nd EveryNDays) First() Date { return FirstRecurAfter(nd, nd.Start.Add(-1)) }

func (nd EveryNDays) String() string {
	return fmt.Sprintf("%s, every %d days", nd.Start.String(), nd.N)
}

type Weekly struct {
	Start    Date
	Weekdays Weekdays
}

// yyyy-mm-dd, weekly, wednesday & saturday & sunday
func ParseWeekly(start Date, terms []string) (Recurrer, bool) {
	if len(terms) < 2 {
		return nil, false
	}

	if terms[0] != "weekly" {
		return nil, false
	}

	wds := Weekdays{}
	for _, wdStr := range strings.Split(terms[1], "&") {
		wd, ok := ParseWeekday(wdStr)
		if !ok {
			continue
		}
		wds = append(wds, wd)
	}
	if len(wds) == 0 {
		return nil, false
	}

	return Weekly{
		Start:    start,
		Weekdays: wds.Unique(),
	}, true
}

func (w Weekly) RecursOn(date Date) bool {
	if w.Start.After(date) {
		return false
	}

	for _, wd := range w.Weekdays {
		if wd == date.Weekday() {
			return true
		}
	}

	return false
}

func (w Weekly) First() Date { return FirstRecurAfter(w, w.Start.Add(-1)) }

func (w Weekly) String() string {
	weekdayStrs := []string{}
	for _, wd := range w.Weekdays {
		weekdayStrs = append(weekdayStrs, wd.String())
	}
	weekdayStr := strings.Join(weekdayStrs, " & ")

	return fmt.Sprintf("%s, weekly, %s", w.Start.String(), strings.ToLower(weekdayStr))
}

type EveryNWeeks struct {
	Start Date
	N     int
}

// yyyy-mm-dd, every 3 weeks
func ParseEveryNWeeks(start Date, terms []string) (Recurrer, bool) {
	if len(terms) != 1 {
		return nil, false
	}

	terms = strings.Split(terms[0], " ")
	if len(terms) != 3 || terms[0] != "every" || terms[2] != "weeks" {
		return nil, false
	}
	n, err := strconv.Atoi(terms[1])
	if err != nil || n < 1 {
		return nil, false
	}

	return EveryNWeeks{
		Start: start,
		N:     n,
	}, true
}

func (enw EveryNWeeks) RecursOn(date Date) bool {
	if enw.Start.After(date) {
		return false
	}
	if enw.Start.Equal(date) {
		return true
	}

	intervalDays := enw.N * 7
	return enw.Start.DaysBetween(date)%intervalDays == 0
}

func (enw EveryNWeeks) First() Date { return FirstRecurAfter(enw, enw.Start.Add(-1)) }

func (enw EveryNWeeks) String() string {
	return fmt.Sprintf("%s, every %d weeks", enw.Start.String(), enw.N)
}

type EveryNMonths struct {
	Start Date
	N     int
}

// yyyy-mm-dd, every 3 months
func ParseEveryNMonths(start Date, terms []string) (Recurrer, bool) {
	if len(terms) != 1 {
		return nil, false
	}

	terms = strings.Split(terms[0], " ")
	if len(terms) != 3 || terms[0] != "every" || terms[2] != "months" {
		return nil, false
	}
	n, err := strconv.Atoi(terms[1])
	if err != nil {
		return nil, false
	}

	return EveryNMonths{
		Start: start,
		N:     n,
	}, true
}

func (enm EveryNMonths) RecursOn(date Date) bool {
	if enm.Start.After(date) {
		return false
	}

	tDate := enm.Start
	for {
		if tDate.Equal(date) {
			return true
		}
		if tDate.After(date) {
			return false
		}
		tDate = tDate.AddMonths(enm.N)
	}

}

func (enm EveryNMonths) First() Date { return FirstRecurAfter(enm, enm.Start.Add(-1)) }

func (enm EveryNMonths) String() string {
	return fmt.Sprintf("%s, every %d months", enm.Start.String(), enm.N)
}
