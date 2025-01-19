package cli

import (
	"errors"
	"fmt"

	"go-mod.ewintr.nl/planner/plan/cli/arg"
	"go-mod.ewintr.nl/planner/plan/command"
	"go-mod.ewintr.nl/planner/plan/command/schedule"
	"go-mod.ewintr.nl/planner/plan/command/task"
	"go-mod.ewintr.nl/planner/sync/client"
)

type CLI struct {
	repos   command.Repositories
	client  client.Client
	cmdArgs []command.CommandArgs
}

func NewCLI(repos command.Repositories, client client.Client) *CLI {
	return &CLI{
		repos:  repos,
		client: client,
		cmdArgs: []command.CommandArgs{
			command.NewSyncArgs(),
			// task
			task.NewShowArgs(), task.NewProjectsArgs(),
			task.NewAddArgs(), task.NewDeleteArgs(), task.NewListArgs(),
			task.NewUpdateArgs(),
			// schedule
			schedule.NewAddArgs(),
		},
	}
}

func (cli *CLI) Run(args []string) error {
	main, fields := arg.FindFields(args)
	for _, ca := range cli.cmdArgs {
		cmd, err := ca.Parse(main, fields)
		switch {
		case errors.Is(err, command.ErrWrongCommand):
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
