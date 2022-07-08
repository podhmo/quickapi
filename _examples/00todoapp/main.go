package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Todo struct {
	Title string `json:"title"`
	Done  bool   `json:"bool"`
}

var todos []Todo = []Todo{
	{Title: "hello", Done: false},
	{Title: "boo", Done: true},
	{Title: "byebye", Done: false},
}

func ListTodo(
	ctx context.Context,
	input Empty,
) (output struct {
	Items []Todo `json:"items"`
}, err error) {
	items := make([]Todo, 0, len(todos))
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

	r.Get("/todos", Lift(ListTodo))

	port := 8080
	log.Printf("[Info]  listening: :%d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		log.Printf("[Error] !! %+v", err)
	}
}

// ----------------------------------------

type Action[I any, O any] func(ctx context.Context, input I) (output O, err error)

type Empty struct{}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"Code"`
}

var (
	Pretty bool
)

func Lift[I any, O any](action Action[I, O]) http.HandlerFunc {
	var iz I
	isEmpty := reflect.TypeOf(iz).NumField() == 0

	return func(w http.ResponseWriter, req *http.Request) {
		var input I
		if !isEmpty {
			if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
				log.Printf("[ERROR] unexpected error (json.Decode): %+v", err) // TODO: structured logging
				writeJSONError(w, req, err, 500)
				return
			}
		}

		code := 500
		output, err := action(req.Context(), input)
		if err != nil {
			// TODO: handling status code
			if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
				if code == 500 {
					log.Printf("[ERROR] unexpected error: %+v", err) // TODO: structured logging
				}
				writeJSONError(w, req, err, code)
			}
		}

		ctx := context.WithValue(req.Context(), render.ContentTypeCtxKey, code)
		render.JSON(w, req.WithContext(ctx), output)
	}
}

func writeJSONError(w http.ResponseWriter, req *http.Request, err error, code int) {
	v := ErrorResponse{Error: "internal server error", Code: code}
	if code != 500 {
		v.Error = err.Error()
	}

	ctx := context.WithValue(req.Context(), render.ContentTypeCtxKey, code)
	render.JSON(w, req.WithContext(ctx), v)
}
