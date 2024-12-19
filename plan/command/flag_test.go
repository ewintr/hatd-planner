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

	valid := item.NewDate(2024, 11, 20)
	f := command.FlagDate{}
	if f.IsSet() {
		t.Errorf("exp false, got true")
	}

	if err := f.Set(valid.String()); err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	if !f.IsSet() {
		t.Errorf("exp true, got false")
	}

	act, ok := f.Get().(item.Date)
	if !ok {
		t.Errorf("exp true, got false")
	}
	if act.String() != valid.String() {
		t.Errorf("exp %v, got %v", valid.String(), act.String())
	}
}

func TestFlagTime(t *testing.T) {
	t.Parallel()

	valid := item.NewTime(12, 30)
	f := command.FlagTime{}
	if f.IsSet() {
		t.Errorf("exp false, got true")
	}

	if err := f.Set(valid.String()); err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	if !f.IsSet() {
		t.Errorf("exp true, got false")
	}

	act, ok := f.Get().(item.Time)
	if !ok {
		t.Errorf("exp true, got false")
	}
	if act.String() != valid.String() {
		t.Errorf("exp %v, got %v", valid.String(), act.String())
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

func TestFlagRecurrer(t *testing.T) {
	t.Parallel()

	validStr := "2024-12-23, daily"
	valid := item.NewRecurrer(validStr)
	f := command.FlagRecurrer{}
	if f.IsSet() {
		t.Errorf("exp false, got true")
	}

	if err := f.Set(validStr); err != nil {
		t.Errorf("exp nil, got %v", err)
	}

	if !f.IsSet() {
		t.Errorf("exp true, got false")
	}

	act, ok := f.Get().(item.Recurrer)
	if !ok {
		t.Errorf("exp true, got false")
	}
	if act != valid {
		t.Errorf("exp %v, got %v", valid, act)
	}
}
