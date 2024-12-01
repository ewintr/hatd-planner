package command_test

import (
	"testing"
	"time"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/command"
)

func TestFlagString(t *testing.T) {
	t.Parallel()

	valid := "test"
	f := command.FlagString{}
	if f.IsSet() {
		t.Errorf("exp false, got true")
	}

	if err := f.Set(valid); err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	if !f.IsSet() {
		t.Errorf("exp true, got false")
	}

	act, ok := f.Get().(string)
	if !ok {
		t.Errorf("exp true, got false")
	}
	if act != valid {
		t.Errorf("exp %v, got %v", valid, act)
	}
}

func TestFlagDate(t *testing.T) {
	t.Parallel()

	valid := time.Date(2024, 11, 20, 0, 0, 0, 0, time.UTC)
	validStr := "2024-11-20"
	f := command.FlagDate{}
	if f.IsSet() {
		t.Errorf("exp false, got true")
	}

	if err := f.Set(validStr); err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	if !f.IsSet() {
		t.Errorf("exp true, got false")
	}

	act, ok := f.Get().(time.Time)
	if !ok {
		t.Errorf("exp true, got false")
	}
	if act != valid {
		t.Errorf("exp %v, got %v", valid, act)
	}
}

func TestFlagTime(t *testing.T) {
	t.Parallel()

	valid := time.Date(0, 1, 1, 12, 30, 0, 0, time.UTC)
	validStr := "12:30"
	f := command.FlagTime{}
	if f.IsSet() {
		t.Errorf("exp false, got true")
	}

	if err := f.Set(validStr); err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	if !f.IsSet() {
		t.Errorf("exp true, got false")
	}

	act, ok := f.Get().(time.Time)
	if !ok {
		t.Errorf("exp true, got false")
	}
	if act != valid {
		t.Errorf("exp %v, got %v", valid, act)
	}
}

func TestFlagDurationTime(t *testing.T) {
	t.Parallel()

	valid := time.Hour
	validStr := "1h"
	f := command.FlagDuration{}
	if f.IsSet() {
		t.Errorf("exp false, got true")
	}

	if err := f.Set(validStr); err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	if !f.IsSet() {
		t.Errorf("exp true, got false")
	}

	act, ok := f.Get().(time.Duration)
	if !ok {
		t.Errorf("exp true, got false")
	}
	if act != valid {
		t.Errorf("exp %v, got %v", valid, act)
	}
}

func TestFlagPeriod(t *testing.T) {
	t.Parallel()

	valid := item.PeriodMonth
	validStr := "month"
	f := command.FlagPeriod{}
	if f.IsSet() {
		t.Errorf("exp false, got true")
	}

	if err := f.Set(validStr); err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	if !f.IsSet() {
		t.Errorf("exp true, got false")
	}

	act, ok := f.Get().(item.RecurPeriod)
	if !ok {
		t.Errorf("exp true, got false")
	}
	if act != valid {
		t.Errorf("exp %v, got %v", valid, act)
	}
}
