package command

import (
	"errors"
	"fmt"
	"strings"
)

const (
	FlagTitle = "title"
	FlagOn    = "on"
	FlagAt    = "at"
	FlagFor   = "for"
)

type Command interface {
	Execute([]string, map[string]string) error
}

type CLI struct {
	Commands []Command
}

func (cli *CLI) Run(args []string) error {
	main, flags, err := ParseFlags(args)
	if err != nil {
		return err
	}
	for _, c := range cli.Commands {
		err := c.Execute(main, flags)
		switch {
		case errors.Is(err, ErrWrongCommand):
			continue
		case err != nil:
			return err
		}
	}

	return fmt.Errorf("could not find matching command")
}

func ParseFlags(args []string) ([]string, map[string]string, error) {
	flags := make(map[string]string)
	main := make([]string, 0)
	var inMain bool
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			inMain = false
			if i+1 >= len(args) {
				return nil, nil, fmt.Errorf("flag wihout value")
			}
			flags[strings.TrimPrefix(args[i], "-")] = args[i+1]
			i++
			continue
		}

		if !inMain && len(main) > 0 {
			return nil, nil, fmt.Errorf("two mains")
		}
		inMain = true
		main = append(main, args[i])
	}

	return main, flags, nil
}
