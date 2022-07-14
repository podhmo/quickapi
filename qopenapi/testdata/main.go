package main

import (
	"context"
	"log"

	"github.com/podhmo/quickapi/qopenapi"
)

type Todo struct {
	ID    string `json:"id"`
	Title string `json:"title" optional:"true"`
	Done  bool   `json:"done"`

	ParentID *string `json:"parentId" optional:"true"` // todo: nullable?
}
type TodoInput struct {
	Sort string `json:"Name"` // id, -id
}
type ListTodoOutput struct {
	Items []Todo `json:"items"`
}

func ListTodo(ctx context.Context, input TodoInput) (output ListTodoOutput, err error) {
	return
}

func main() {
	r, err := qopenapi.NewRouter()
	if err != nil {
		log.Fatalf("!! %+v", err)
	}

	qopenapi.Get(r, "/todo", ListTodo)

	ctx := context.Background()
	if err := qopenapi.EmitDoc(ctx, r); err != nil {
		log.Fatalf("!! %+v", err)
	}
}
