package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"
	"go-mod.ewintr.nl/planner/item"
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

	repo, err := NewSqlite(conf.DBPath)
	if err != nil {
		fmt.Printf("could not open db file: %s\n", err)
		os.Exit(1)
	}

	app := &cli.App{
		Name:  "plan",
		Usage: "Plan your day with events",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List everything",
				Action: func(cCtx *cli.Context) error {
					return List(repo)
				},
			},
			{
				Name:  "add",
				Usage: "Add a new event",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "The event that will happen",
					},
					&cli.StringFlag{
						Name:    "date",
						Aliases: []string{"d"},
						Usage:   "The date, in YYYY-MM-DD format",
					},
					&cli.StringFlag{
						Name:    "time",
						Aliases: []string{"t"},
						Usage:   "The time, in HH:MM format. If omitted, the event will last the whole day",
					},
					&cli.StringFlag{
						Name:    "for",
						Aliases: []string{"f"},
						Usage:   "The duration, in show format (e.g. 1h30m)",
					},
				},
				Action: func(cCtx *cli.Context) error {
					return Add(cCtx, repo)
				},
			},
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

func List(repo EventRepo) error {
	all, err := repo.FindAll()
	if err != nil {
		return err
	}
	for _, e := range all {
		fmt.Printf("%s\t%s\t%s\t%s\n", e.ID, e.Title, e.Start.Format(time.DateTime), e.Duration.String())
	}

	return nil
}

func Add(cCtx *cli.Context, repo EventRepo) error {
	desc := cCtx.String("name")
	date, err := time.Parse("2006-01-02", cCtx.String("date"))
	if err != nil {
		return fmt.Errorf("could not parse date: %v", err)
	}

	one := item.Event{
		ID: "a",
		EventBody: item.EventBody{
			Title: desc,
			Start: date,
		},
	}
	if err := repo.Store(one); err != nil {
		return fmt.Errorf("could not store event: %v", err)
	}

	return nil
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
