package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
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

func mount(r chi.Router) {
	r.Get("/todos", quickapi.Lift(ListTodo))
}

func main() {
	ctx := context.Background()
	r := quickapi.DefaultRouter()
	mount(r)

	port := 8080
	addr := fmt.Sprintf(":%d", port)
	if err := quickapi.NewServer(addr, r, 5*time.Second).ListenAndServe(ctx); err != nil {
		log.Printf("[Error] !! %+v", err)
	}
}
