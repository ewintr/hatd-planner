package command

import (
	"fmt"
)

type ListArgs struct {
}

func NewListArgs() ListArgs {
	return ListArgs{}
}

func (la ListArgs) Parse(main []string, flags map[string]string) (Command, error) {
	if len(main) > 0 && main[0] != "list" {
		return nil, ErrWrongCommand
	}

	return &List{}, nil
}

type List struct {
}

func (list *List) Do(deps Dependencies) error {
	localIDs, err := deps.LocalIDRepo.FindAll()
	if err != nil {
		return fmt.Errorf("could not get local ids: %v", err)
	}
	all, err := deps.TaskRepo.FindAll()
	if err != nil {
		return err
	}
	for _, e := range all {
		lid, ok := localIDs[e.ID]
		if !ok {
			return fmt.Errorf("could not find local id for %s", e.ID)
		}
		fmt.Printf("%s\t%d\t%s\t%s\t%s\n", e.ID, lid, e.Title, e.Date.String(), e.Duration.String())
	}

	return nil
}
