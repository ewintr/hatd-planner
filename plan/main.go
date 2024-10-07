package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"go-mod.ewintr.nl/planner/plan/command"
	"go-mod.ewintr.nl/planner/plan/storage/sqlite"
	"go-mod.ewintr.nl/planner/sync/client"
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

	localIDRepo, eventRepo, syncRepo, err := sqlite.NewSqlites(conf.DBPath)
	if err != nil {
		fmt.Printf("could not open db file: %s\n", err)
		os.Exit(1)
	}

	syncClient := client.New(conf.SyncURL, conf.ApiKey)

	app := &cli.App{
		Name:  "plan",
		Usage: "Plan your day with events",
		Commands: []*cli.Command{
			command.NewAddCmd(localIDRepo, eventRepo, syncRepo),
			command.NewListCmd(localIDRepo, eventRepo),
			command.NewUpdateCmd(localIDRepo, eventRepo, syncRepo),
			command.NewDeleteCmd(localIDRepo, eventRepo, syncRepo),
			command.NewSyncCmd(syncClient, syncRepo, localIDRepo, eventRepo),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Configuration struct {
	DBPath  string `yaml:"dbpath"`
	SyncURL string `yaml:"sync_url"`
	ApiKey  string `yaml:"api_key"`
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
