package main

import (
	"fmt"
	"os"
	"time"

	"go-mod.ewintr.nl/planner/item"
)

func main() {
	fmt.Println("cal")

	repo, err := NewSqlite("test.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	one := item.Event{
		ID: "a",
		EventBody: item.EventBody{
			Title: "title",
			Start: time.Now(),
			End:   time.Now().Add(-5 * time.Second),
		},
	}
	if err := repo.Delete(one.ID); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	all, err := repo.FindAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("all: %+v\n", all)

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
