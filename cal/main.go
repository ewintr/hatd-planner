package main

import (
	"fmt"
	"os"
	"time"

	"go-mod.ewintr.nl/planner/sync/client"
	"go-mod.ewintr.nl/planner/sync/planner"
)

func main() {
	fmt.Println("cal")

	c := client.NewClient("http://localhost:8092", "testKey")
	items, err := c.Updated([]planner.Kind{planner.KindEvent}, time.Time{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", items)

	i := planner.Item{
		ID:      "id-1",
		Kind:    planner.KindEvent,
		Updated: time.Now(),
		Body:    "body",
	}
	if err := c.Update([]planner.Item{i}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	items, err = c.Updated([]planner.Kind{planner.KindEvent}, time.Time{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", items)
}
