package command

import (
	"fmt"
	"time"

	"go-mod.ewintr.nl/planner/item"
)

type ArgSet struct {
	Main  string
	Flags map[string]Flag
}

func (as *ArgSet) Set(name, val string) error {
	f, ok := as.Flags[name]
	if !ok {
		return fmt.Errorf("unknown flag %s", name)
	}
	return f.Set(val)
}

func (as *ArgSet) IsSet(name string) bool {
	f, ok := as.Flags[name]
	if !ok {
		return false
	}
	return f.IsSet()
}

func (as *ArgSet) GetString(name string) string {
	flag, ok := as.Flags[name]
	if !ok {
		return ""
	}
	val, ok := flag.Get().(string)
	if !ok {
		return ""
	}
	return val
}

func (as *ArgSet) GetDate(name string) item.Date {
	flag, ok := as.Flags[name]
	if !ok {
		return item.Date{}
	}
	val, ok := flag.Get().(item.Date)
	if !ok {
		return item.Date{}
	}
	return val
}

func (as *ArgSet) GetTime(name string) item.Time {
	flag, ok := as.Flags[name]
	if !ok {
		return item.Time{}
	}
	val, ok := flag.Get().(item.Time)
	if !ok {
		return item.Time{}
	}
	return val
}

func (as *ArgSet) GetDuration(name string) time.Duration {
	flag, ok := as.Flags[name]
	if !ok {
		return time.Duration(0)
	}
	val, ok := flag.Get().(time.Duration)
	if !ok {
		return time.Duration(0)
	}
	return val
}

func (as *ArgSet) GetRecurrer(name string) item.Recurrer {
	flag, ok := as.Flags[name]
	if !ok {
		return nil
	}
	val, ok := flag.Get().(item.Recurrer)
	if !ok {
		return nil
	}
	return val
}

func (as *ArgSet) GetInt(name string) int {
	flag, ok := as.Flags[name]
	if !ok {
		return 0
	}
	val, ok := flag.Get().(int)
	if !ok {
		return 0
	}
	return val
}
