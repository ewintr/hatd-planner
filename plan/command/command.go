package command

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"go-mod.ewintr.nl/planner/plan/format"
	"go-mod.ewintr.nl/planner/plan/storage"
	"go-mod.ewintr.nl/planner/sync/client"
)

const (
	DateFormat = "2006-01-02"
	TimeFormat = "15:04"
)

var (
	ErrWrongCommand = errors.New("wrong command")
	ErrInvalidArg   = errors.New("invalid argument")
)

type Dependencies struct {
	LocalIDRepo storage.LocalID
	TaskRepo    storage.Task
	SyncRepo    storage.Sync
	SyncClient  client.Client
}

type CommandArgs interface {
	Parse(main []string, fields map[string]string) (Command, error)
}

type Command interface {
	Do(deps Dependencies) ([][]string, error)
}

type CLI struct {
	deps    Dependencies
	cmdArgs []CommandArgs
}

func NewCLI(deps Dependencies) *CLI {
	return &CLI{
		deps: deps,
		cmdArgs: []CommandArgs{
			NewShowArgs(),
			NewAddArgs(), NewDeleteArgs(), NewListArgs(),
			NewSyncArgs(), NewUpdateArgs(),
		},
	}
}

func (cli *CLI) Run(args []string) error {
	main, fields := FindFields(args)
	for _, ca := range cli.cmdArgs {
		cmd, err := ca.Parse(main, fields)
		switch {
		case errors.Is(err, ErrWrongCommand):
			continue
		case err != nil:
			return err
		}

		data, err := cmd.Do(cli.deps)
		if err != nil {
			return err
		}

		switch {
		case len(data) == 0:
		case len(data) == 1 && len(data[0]) == 1:
			fmt.Println(data[0][0])
		default:
			fmt.Printf("\n%s\n", format.Table(data))
		}
		return nil
	}

	return fmt.Errorf("could not find matching command")
}

func FindFields(args []string) ([]string, map[string]string) {
	fields := make(map[string]string)
	main := make([]string, 0)
	for i := 0; i < len(args); i++ {
		if k, v, ok := strings.Cut(args[i], ":"); ok && !strings.Contains(k, " ") {
			fields[k] = v
			continue
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
					return nil, fmt.Errorf("%w: duplicate field: %v", ErrInvalidArg, tk)
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
		return nil, fmt.Errorf("%w: unknown field(s): %v", ErrInvalidArg, strings.Join(ks, ","))
	}

	return res, nil
}
