package command_test

// func TestArgSet(t *testing.T) {
// 	t.Parallel()

// 	as := command.ArgSet{
// 		Main: "main",
// 		Flags: map[string]string{
// 			"name 1": "value 1",
// 			"name 2": "value 2",
// 			"name 3": "value 3",
// 		},
// 	}

// 	t.Run("hasflag", func(t *testing.T) {
// 		t.Run("true", func(t *testing.T) {
// 			if has := as.HasFlag("name 1"); !has {
// 				t.Errorf("exp true, got %v", has)
// 			}
// 		})
// 		t.Run("false", func(t *testing.T) {
// 			if has := as.HasFlag("unknown"); has {
// 				t.Errorf("exp false, got %v", has)
// 			}
// 		})
// 	})

// 	t.Run("flag", func(t *testing.T) {
// 		t.Run("known", func(t *testing.T) {
// 			if val := as.Flag("name 1"); val != "value 1" {
// 				t.Errorf("exp value 1, got %v", val)
// 			}
// 		})
// 		t.Run("unknown", func(t *testing.T) {
// 			if val := as.Flag("unknown"); val != "" {
// 				t.Errorf(`exp "", got %v`, val)
// 			}
// 		})
// 	})

// 	t.Run("setflag", func(t *testing.T) {
// 		exp := "new value"
// 		as.SetFlag("new name", exp)
// 		if act := as.Flag("new name"); exp != act {
// 			t.Errorf("exp %v, got %v", exp, act)
// 		}
// 	})
// }

// func TestParseArgs(t *testing.T) {
// 	t.Parallel()

// 	for _, tc := range []struct {
// 		name   string
// 		args   []string
// 		expAS  *command.ArgSet
// 		expErr bool
// 	}{
// 		{
// 			name: "empty",
// 			expAS: &command.ArgSet{
// 				Flags: map[string]string{},
// 			},
// 		},
// 		{
// 			name: "just main",
// 			args: []string{"one", "two three", "four"},
// 			expAS: &command.ArgSet{
// 				Main:  "one two three four",
// 				Flags: map[string]string{},
// 			},
// 		},
// 		{
// 			name: "with flags",
// 			args: []string{"-flag1", "value1", "one", "two", "-flag2", "value2", "-flag3", "value3"},
// 			expAS: &command.ArgSet{
// 				Main: "one two",
// 				Flags: map[string]string{
// 					"flag1": "value1",
// 					"flag2": "value2",
// 					"flag3": "value3",
// 				},
// 			},
// 		},
// 		{
// 			name:   "flag without value",
// 			args:   []string{"one", "two", "-flag1"},
// 			expErr: true,
// 		},
// 		{
// 			name:   "split main",
// 			args:   []string{"one", "-flag1", "value1", "two"},
// 			expErr: true,
// 		},
// 	} {
// 		t.Run(tc.name, func(t *testing.T) {
// 			actAS, actErr := command.ParseArgs(tc.args)
// 			if tc.expErr != (actErr != nil) {
// 				t.Errorf("exp %v, got %v", tc.expErr, actErr)
// 			}
// 			if tc.expErr {
// 				return
// 			}
// 			if diff := cmp.Diff(tc.expAS, actAS); diff != "" {
// 				t.Errorf("(exp +, got -)\n%s", diff)
// 			}
// 		})
// 	}
// }
