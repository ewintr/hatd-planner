package command

import (
	"errors"
	"fmt"
	"slices"
	"strings"

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

type Repositories interface {
	Begin() (*storage.Tx, error)
	LocalID(tx *storage.Tx) storage.LocalID
	Sync(tx *storage.Tx) storage.Sync
	Task(tx *storage.Tx) storage.Task
}

type CommandArgs interface {
	Parse(main []string, fields map[string]string) (Command, error)
}

type Command interface {
	Do(repos Repositories, client client.Client) (CommandResult, error)
}

type CommandResult interface {
	Render() string
}

type CLI struct {
	repos   Repositories
	client  client.Client
	cmdArgs []CommandArgs
}

func NewCLI(repos Repositories, client client.Client) *CLI {
	return &CLI{
		repos:  repos,
		client: client,
		cmdArgs: []CommandArgs{
			NewShowArgs(), NewProjectsArgs(),
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

		result, err := cmd.Do(cli.repos, cli.client)
		if err != nil {
			return err
		}
		fmt.Println(result.Render())

		return nil
	}

	return fmt.Errorf("could not find matching command")
}

func FindFields(args []string) ([]string, map[string]string) {
	fields := make(map[string]string)
	main := make([]string, 0)
	for i := 0; i < len(args); i++ {
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
