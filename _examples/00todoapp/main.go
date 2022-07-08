package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/podhmo/quickapi"
)

type Todo struct {
	Title string `json:"title"`
	Done  bool   `json:"bool"`
}

var todos []Todo

// func init() {
// 	todos = []Todo{
// 		{Title: "hello", Done: false},
// 		{Title: "boo", Done: true},
// 		{Title: "byebye", Done: false},
// 	}
// }

func ListTodo(
	ctx context.Context,
	input quickapi.Empty,
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
	output.Items = items
	return
}

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/todos", quickapi.Lift(ListTodo))

	port := 8080
	log.Printf("[Info]  listening: :%d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		log.Printf("[Error] !! %+v", err)
	}
}
