package main

import (
	"fmt"
	"os"
	"path/filepath"

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

	cli := command.CLI{
		Commands: []command.Command{
			command.NewAdd(localIDRepo, eventRepo, syncRepo),
			command.NewList(localIDRepo, eventRepo),
			command.NewUpdate(localIDRepo, eventRepo, syncRepo),
			command.NewDelete(localIDRepo, eventRepo, syncRepo),
			command.NewSync(syncClient, syncRepo, localIDRepo, eventRepo),
		},
	}

	if err := cli.Run(os.Args); err != nil {
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
