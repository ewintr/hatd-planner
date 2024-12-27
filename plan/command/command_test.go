package command_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go-mod.ewintr.nl/planner/plan/command"
)

func TestFindFields(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name      string
		args      []string
		expMain   []string
		expFields map[string]string
	}{
		{
			name:      "empty",
			expMain:   []string{},
			expFields: map[string]string{},
		},
		{
			name:      "just main",
			args:      []string{"one", "two three", "four"},
			expMain:   []string{"one", "two three", "four"},
			expFields: map[string]string{},
		},
		{
			name:    "with flags",
			args:    []string{"flag1:value1", "one", "two", "flag2:value2", "flag3:value3"},
			expMain: []string{"one", "two"},
			expFields: map[string]string{
				"flag1": "value1",
				"flag2": "value2",
				"flag3": "value3",
			},
		},
		{
			name:    "flag without value",
			args:    []string{"one", "two", "flag1:"},
			expMain: []string{"one", "two"},
			expFields: map[string]string{
				"flag1": "",
			},
		},
		{
			name:    "split main",
			args:    []string{"one", "flag1:value1", "two"},
			expMain: []string{"one", "two"},
			expFields: map[string]string{
				"flag1": "value1",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actMain, actFields := command.FindFields(tc.args)
			if diff := cmp.Diff(tc.expMain, actMain); diff != "" {
				t.Errorf("(exp +, got -)\n%s", diff)
			}
			if diff := cmp.Diff(tc.expFields, actFields); diff != "" {
				t.Errorf("(exp +, got -)\n%s", diff)
			}
		})
	}
}

func TestResolveFields(t *testing.T) {
	t.Parallel()

	tmpl := map[string][]string{
		"one": []string{"a", "b"},
		"two": []string{"c", "d"},
	}
	for _, tc := range []struct {
		name   string
		fields map[string]string
		expRes map[string]string
		expErr bool
	}{
		{
			name:   "empty",
			expRes: map[string]string{},
		},
		{
			name: "unknown",
			fields: map[string]string{
				"unknown": "value",
			},
			expErr: true,
		},
		{
			name: "duplicate",
			fields: map[string]string{
				"a": "val1",
				"b": "val2",
			},
			expErr: true,
		},
		{
			name: "valid",
			fields: map[string]string{
				"a": "val1",
				"d": "val2",
			},
			expRes: map[string]string{
				"one": "val1",
				"two": "val2",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actRes, actErr := command.ResolveFields(tc.fields, tmpl)
			if tc.expErr != (actErr != nil) {
				t.Errorf("exp %v, got %v", tc.expErr, actErr != nil)
			}
			if tc.expErr {
				return
			}
			if diff := cmp.Diff(tc.expRes, actRes); diff != "" {
				t.Errorf("(+exp, -got)%s\n", diff)
			}
		})
	}
}
