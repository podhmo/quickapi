package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/podhmo/quickapi"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"bool"`
}

var todos []Todo

func init() {
	todos = []Todo{
		{ID: 1, Title: "hello", Done: false},
		{ID: 2, Title: "boo", Done: true},
		{ID: 3, Title: "byebye", Done: false},
	}
}

func ListTodo(
	ctx context.Context,
	input struct {
		Sort string `query:"sort"` // enum -id, id
	},
) (output struct {
	Items []Todo `json:"items"`
}, err error) {
	var items []Todo
	for _, x := range todos {
		if x.Done {
			continue
		}
		items = append(items, x)
	}
	if strings.HasPrefix(input.Sort, "-") {
		sort.Slice(items, func(i, j int) bool { return items[i].ID > items[j].ID })
	}
	output.Items = items
	return
}

func main() {
	r := quickapi.DefaultRouter()
	r.Get("/todos", quickapi.Lift(ListTodo))

	port := 8080
	log.Printf("[Info]  listening: :%d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		log.Printf("[Error] !! %+v", err)
	}
}
