package command

import (
	"fmt"
	"sort"

	"go-mod.ewintr.nl/planner/plan/format"
)

type ProjectsArgs struct{}

func NewProjectsArgs() ProjectsArgs {
	return ProjectsArgs{}
}

func (pa ProjectsArgs) Parse(main []string, fields map[string]string) (Command, error) {
	if len(main) != 1 || main[0] != "projects" {
		return nil, ErrWrongCommand
	}

	return Projects{}, nil
}

type Projects struct{}

func (ps Projects) Do(deps Dependencies) (CommandResult, error) {
	projects, err := deps.TaskRepo.Projects()
	if err != nil {
		return nil, fmt.Errorf("could not find projects: %v", err)
	}

	return ProjectsResult{
		Projects: projects,
	}, nil
}

type ProjectsResult struct {
	Projects map[string]int
}

func (psr ProjectsResult) Render() string {
	projects := make([]string, 0, len(psr.Projects))
	for pr := range psr.Projects {
		projects = append(projects, pr)
	}
	sort.Strings(projects)
	data := [][]string{{"projects", "count"}}
	for _, p := range projects {
		data = append(data, []string{p, fmt.Sprintf("%d", psr.Projects[p])})
	}

	return fmt.Sprintf("\n%s\n", format.Table(data))
}
