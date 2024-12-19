package command_test

import (
	"testing"
	"time"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/command"
)

func TestArgSet(t *testing.T) {
	for _, tt := range []struct {
		name     string
		flags    map[string]command.Flag
		flagName string
		setValue string
		exp      interface{}
		expErr   bool
	}{
		{
			name: "string flag success",
			flags: map[string]command.Flag{
				"title": &command.FlagString{Name: "title"},
			},
			flagName: "title",
			setValue: "test title",
			exp:      "test title",
		},
		{
			name: "date flag success",
			flags: map[string]command.Flag{
				"date": &command.FlagDate{Name: "date"},
			},
			flagName: "date",
			setValue: "2024-01-02",
			exp:      time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "time flag success",
			flags: map[string]command.Flag{
				"time": &command.FlagTime{Name: "time"},
			},
			flagName: "time",
			setValue: "15:04",
			exp:      time.Date(0, 1, 1, 15, 4, 0, 0, time.UTC),
		},
		{
			name: "duration flag success",
			flags: map[string]command.Flag{
				"duration": &command.FlagDuration{Name: "duration"},
			},
			flagName: "duration",
			setValue: "2h30m",
			exp:      2*time.Hour + 30*time.Minute,
		},
		{
			name: "recur period flag success",
			flags: map[string]command.Flag{
				"recur": &command.FlagRecurrer{Name: "recur"},
			},
			flagName: "recur",
			setValue: "2024-12-23, daily",
			exp:      item.NewRecurrer("2024-12-23, daily"),
		},
		{
			name:     "unknown flag error",
			flags:    map[string]command.Flag{},
			flagName: "unknown",
			setValue: "value",
			expErr:   true,
		},
		{
			name: "invalid date format error",
			flags: map[string]command.Flag{
				"date": &command.FlagDate{Name: "date"},
			},
			flagName: "date",
			setValue: "invalid",
			expErr:   true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			as := &command.ArgSet{
				Main:  "test",
				Flags: tt.flags,
			}

			err := as.Set(tt.flagName, tt.setValue)
			if (err != nil) != tt.expErr {
				t.Errorf("ArgSet.Set() error = %v, expErr %v", err, tt.expErr)
				return
			}

			if tt.expErr {
				return
			}

			if !as.IsSet(tt.flagName) {
				t.Errorf("ArgSet.IsSet() = false, want true for flag %s", tt.flagName)
			}
		})
	}
}
