package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"go-mod.ewintr.nl/planner/plan/command"
	"go-mod.ewintr.nl/planner/plan/storage/sqlite"
	"gopkg.in/yaml.v3"
)

func main() {
	confPath, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("could not get config path: %s\n", err)
		os.Exit(1)
	}
	conf, err := LoadConfig(filepath.Join(confPath, "planner", "plan", "config.yaml"))
	if err != nil {
		fmt.Printf("could not open config file: %s\n", err)
		os.Exit(1)
	}

	localIDRepo, eventRepo, err := sqlite.NewSqlites(conf.DBPath)
	if err != nil {
		fmt.Printf("could not open db file: %s\n", err)
		os.Exit(1)
	}

	app := &cli.App{
		Name:  "plan",
		Usage: "Plan your day with events",
		Commands: []*cli.Command{
			command.NewAddCmd(localIDRepo, eventRepo),
			command.NewListCmd(localIDRepo, eventRepo),
			command.NewUpdateCmd(localIDRepo, eventRepo),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// all, err := repo.FindAll()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// fmt.Printf("all: %+v\n", all)

	// c := client.NewClient("http://localhost:8092", "testKey")
	// items, err := c.Updated([]item.Kind{item.KindEvent}, time.Time{})
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// fmt.Printf("%+v\n", items)

	// i := item.Item{
	// 	ID:      "id-1",
	// 	Kind:    item.KindEvent,
	// 	Updated: time.Now(),
	// 	Body:    "body",
	// }
	// if err := c.Update([]item.Item{i}); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// items, err = c.Updated([]item.Kind{item.KindEvent}, time.Time{})
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// fmt.Printf("%+v\n", items)
}

type Configuration struct {
	DBPath string `yaml:"dbpath"`
}

func LoadConfig(path string) (Configuration, error) {
	confFile, err := os.ReadFile(path)
	if err != nil {
		return Configuration{}, fmt.Errorf("could not open file: %s", err)
	}
	var conf Configuration
	if err := yaml.Unmarshal(confFile, &conf); err != nil {
		return Configuration{}, fmt.Errorf("could not unmarshal config: %s", err)
	}

	return conf, nil
}
