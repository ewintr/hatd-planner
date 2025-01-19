package main

import (
	"fmt"
	"os"
	"path/filepath"

	"go-mod.ewintr.nl/planner/plan/cli"
	"go-mod.ewintr.nl/planner/plan/storage/sqlite"
	"go-mod.ewintr.nl/planner/sync/client"
	"gopkg.in/yaml.v3"
)

func main() {
	confPath := os.Getenv("PLAN_CONFIG_PATH")
	if confPath == "" {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			fmt.Printf("could not get config path: %s\n", err)
			os.Exit(1)
		}
		confPath = filepath.Join(userConfigDir, "planner", "plan", "config.yaml")
	}
	conf, err := LoadConfig(confPath)
	if err != nil {
		fmt.Printf("could not open config file: %s\n", err)
		os.Exit(1)
	}

	repos, err := sqlite.NewSqlites(conf.DBPath)
	if err != nil {
		fmt.Printf("could not open db file: %s\n", err)
		os.Exit(1)
	}

	syncClient := client.New(conf.SyncURL, conf.ApiKey)

	cli := cli.NewCLI(repos, syncClient)
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Configuration struct {
	DBPath  string `yaml:"db_path"`
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
