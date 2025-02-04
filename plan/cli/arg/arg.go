package arg

import (
	"fmt"
	"slices"
	"strings"

	"go-mod.ewintr.nl/planner/plan/command"
)

func FindFields(args []string) ([]string, map[string]string) {
	fields := make(map[string]string)
	main := make([]string, 0)
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "http://") || strings.HasPrefix(args[i], "https://") {
			main = append(main, args[i])
			continue
		}
		// normal key:value
		if k, v, ok := strings.Cut(args[i], ":"); ok && !strings.Contains(k, " ") {
			fields[k] = v
			continue
		}
		// empty key:
		if !strings.Contains(args[i], " ") && strings.HasSuffix(args[i], ":") {
			k := strings.TrimSuffix(args[i], ":")
			fields[k] = ""
		}
		main = append(main, args[i])
	}

	return main, fields
}
func ResolveFields(fields map[string]string, tmpl map[string][]string) (map[string]string, error) {
	res := make(map[string]string)
	for k, v := range fields {
		for tk, tv := range tmpl {
			if slices.Contains(tv, k) {
				if _, ok := res[tk]; ok {
					return nil, fmt.Errorf("%w: duplicate field: %v", command.ErrInvalidArg, tk)
				}
				res[tk] = v
				delete(fields, k)
			}
		}
	}
	if len(fields) > 0 {
		ks := make([]string, 0, len(fields))
		for k := range fields {
			ks = append(ks, k)
		}
		return nil, fmt.Errorf("%w: unknown field(s): %v", command.ErrInvalidArg, strings.Join(ks, ","))
	}

	return res, nil
}
