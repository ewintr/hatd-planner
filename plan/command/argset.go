package command

import (
	"fmt"
	"time"
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

func (as *ArgSet) GetTime(name string) time.Time {
	flag, ok := as.Flags[name]
	if !ok {
		return time.Time{}
	}
	val, ok := flag.Get().(time.Time)
	if !ok {
		return time.Time{}
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
