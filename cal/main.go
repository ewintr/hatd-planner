package main

import (
	"fmt"
	"os"
	"time"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/sync/client"
)

func main() {
	fmt.Println("cal")

	c := client.NewClient("http://localhost:8092", "testKey")
	items, err := c.Updated([]item.Kind{item.KindEvent}, time.Time{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", items)

	i := item.Item{
		ID:      "id-1",
		Kind:    item.KindEvent,
		Updated: time.Now(),
		Body:    "body",
	}
	if err := c.Update([]item.Item{i}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	items, err = c.Updated([]item.Kind{item.KindEvent}, time.Time{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", items)
}
