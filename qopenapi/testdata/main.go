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
type TodoInputSort string // enum
type TodoInput struct {
	Sort TodoInputSort `openapi:"query" query:"sort"`
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

	qopenapi.Get(r, "/todo", ListTodo).Description("List")
	qopenapi.DefineStringEnum[TodoInputSort](r, "id", "-id")

	ctx := context.Background()
	if err := qopenapi.EmitDoc(ctx, r); err != nil {
		log.Fatalf("!! %+v", err)
	}
}
