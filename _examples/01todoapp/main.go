package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/qbind"
	"github.com/podhmo/quickapi/qdump"
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

func GetTodo(
	ctx context.Context,
	input struct {
		ID int `path:"id"`
	},
) (output Todo, err error) {
	for _, x := range todos {
		if x.ID == input.ID {
			return x, nil
		}
	}
	return Todo{}, quickapi.NewAPIError(fmt.Errorf("not found"), 404)
}

func mount(r chi.Router) {
	r.Get("/todos", quickapi.Lift(ListTodo))
	r.Get("/todos/{id}", quickapi.Lift(GetTodo))
	r.Method("GET", "/todos/foo/{id}", quickapi.NewHandler(GetTodo, qdump.Dump[Todo]))
}

func main() {
	ctx := context.Background()
	r := quickapi.DefaultRouter()
	mount(r)

	// TODO: validation typo, something like  r.Get("/todos/{todo_id}", ...) // id != todo_id
	if err := chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		switch h := handler.(type) {
		case http.HandlerFunc:
			fmt.Println("@@", method, route, quickapi.MetadataFromHandlerFunc(h).PathVars)
		case interface{ Metadata() qbind.Metadata }:
			fmt.Println("**", method, route, h.Metadata())
		}
		return nil
	}); err != nil {
		log.Printf("[Error] ! %+v", err)
		return
	}

	port := 8080
	addr := fmt.Sprintf(":%d", port)
	if err := quickapi.NewServer(addr, r, 5*time.Second).ListenAndServe(ctx); err != nil {
		log.Printf("[Error] !! %+v", err)
	}
}
