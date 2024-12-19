package item

import (
	"encoding/json"
	"time"
)

const (
	TimeFormat = "15:04"
)

type Time struct {
	t time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *Time) UnmarshalJSON(data []byte) error {
	timeString := ""
	if err := json.Unmarshal(data, &timeString); err != nil {
		return err
	}
	nt := NewTimeFromString(timeString)
	t.t = nt.Time()

	return nil
}

func NewTime(hour, minute int) Time {
	return Time{
		t: time.Date(0, 0, 0, hour, minute, 0, 0, time.UTC),
	}
}

func NewTimeFromString(timeStr string) Time {
	tm, err := time.Parse(TimeFormat, timeStr)
	if err != nil {
		return Time{t: time.Time{}}
	}

	return Time{t: tm}
}

func (t *Time) String() string {
	return t.t.Format(TimeFormat)
}

func (t *Time) Time() time.Time {
	return t.t
}

func (t *Time) IsZero() bool {
	return t.t.IsZero()
}

func (t *Time) Hour() int {
	return t.t.Hour()
}

func (t *Time) Minute() int {
	return t.t.Minute()
}

func (t *Time) Add(d time.Duration) Time {
	return Time{
		t: t.t.Add(d),
	}
}
