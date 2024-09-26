package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

func main() {
	confPath, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("could not get config path: %s\n", err)
		os.Exit(1)
	}
	conf, err := LoadConfig(filepath.Join(confPath, "planner", "cal", "config.yaml"))
	if err != nil {
		fmt.Printf("could not open config file: %s\n", err)
		os.Exit(1)
	}

	repo, err := NewSqlite(conf.DBPath)
	if err != nil {
		fmt.Printf("could not open db file: %s\n", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Println("need a sub command")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "list":
		all, err := repo.FindAll()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, e := range all {
			fmt.Printf("%s\t%s\t%s\t%s\n", e.ID, e.Title, e.Start.Format(time.DateTime), e.Duration.String())
		}
	default:
		fmt.Printf("unkown command %s\n", os.Args[1])
	}

	// one := item.Event{
	// 	ID: "a",
	// 	EventBody: item.EventBody{
	// 		Title: "title",
	// 		Start: time.Now(),
	// 		End:   time.Now().Add(-5 * time.Second),
	// 	},
	// }
	// if err := repo.Store(one); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

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
