package main

import (
	"context"
	"log"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/podhmo/quickapi/qopenapi"
)

type Todo struct {
	ID    string `json:"id"`
	Title string `json:"title" optional:"true"`
	Done  bool   `json:"done"`

	ParentID *string `json:"parentId" optional:"true"` // todo: nullable?
}
type TodoInput struct {
	Sort string `openapi:"query" query:"sort"` // id, -id
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

	qopenapi.Get(r, "/todo", ListTodo).After(func(op *openapi3.Operation) {
		op.Description = "List"
	})

	ctx := context.Background()
	if err := qopenapi.EmitDoc(ctx, r); err != nil {
		log.Fatalf("!! %+v", err)
	}
}
