package item

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-cmp/cmp"
)

type ScheduleBody struct {
	Title string `json:"title"`
}

type Schedule struct {
	ID        string   `json:"id"`
	Date      Date     `json:"date"`
	Recurrer  Recurrer `json:"recurrer"`
	RecurNext Date     `json:"recurNext"`
	ScheduleBody
}

func NewSchedule(i Item) (Schedule, error) {
	if i.Kind != KindSchedule {
		return Schedule{}, ErrInvalidKind
	}

	var s Schedule
	if err := json.Unmarshal([]byte(i.Body), &s); err != nil {
		return Schedule{}, fmt.Errorf("could not unmarshal item body: %v", err)
	}

	s.ID = i.ID
	s.Date = i.Date

	return s, nil
}

func (s Schedule) Item() (Item, error) {
	body, err := json.Marshal(s.ScheduleBody)
	if err != nil {
		return Item{}, fmt.Errorf("could not marshal schedule body: %v", err)
	}

	return Item{
		ID:   s.ID,
		Kind: KindSchedule,
		Date: s.Date,
		Body: string(body),
	}, nil
}

func ScheduleDiff(a, b Schedule) string {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)

	return cmp.Diff(string(aJSON), string(bJSON))
}

func ScheduleDiffs(a, b []Schedule) string {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)

	return cmp.Diff(string(aJSON), string(bJSON))
}
