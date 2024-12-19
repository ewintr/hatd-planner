package command_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go-mod.ewintr.nl/planner/plan/command"
)

func TestParseArgs(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		args     []string
		expMain  []string
		expFlags map[string]string
		expErr   bool
	}{
		{
			name:     "empty",
			expMain:  []string{},
			expFlags: map[string]string{},
		},
		{
			name:     "just main",
			args:     []string{"one", "two three", "four"},
			expMain:  []string{"one", "two three", "four"},
			expFlags: map[string]string{},
		},
		{
			name:    "with flags",
			args:    []string{"-flag1", "value1", "one", "two", "-flag2", "value2", "-flag3", "value3"},
			expMain: []string{"one", "two"},
			expFlags: map[string]string{
				"flag1": "value1",
				"flag2": "value2",
				"flag3": "value3",
			},
		},
		{
			name:   "flag without value",
			args:   []string{"one", "two", "-flag1"},
			expErr: true,
		},
		{
			name:   "split main",
			args:   []string{"one", "-flag1", "value1", "two"},
			expErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actMain, actFlags, actErr := command.ParseFlags(tc.args)
			if tc.expErr != (actErr != nil) {
				t.Errorf("exp %v, got %v", tc.expErr, actErr)
			}
			if tc.expErr {
				return
			}
			if diff := cmp.Diff(tc.expMain, actMain); diff != "" {
				t.Errorf("(exp +, got -)\n%s", diff)
			}
			if diff := cmp.Diff(tc.expFlags, actFlags); diff != "" {
				t.Errorf("(exp +, got -)\n%s", diff)
			}
		})
	}
}
