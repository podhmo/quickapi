package main

import (
	"context"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/podhmo/quickapi/experimental/define"
)

type Todo struct {
	ID    string `json:"id"`
	Title string `json:"title,omitempty"`
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
	bc, err := define.NewBuildContext(chi.NewRouter())
	if err != nil {
		log.Fatalf("!! %+v", err)
	}

	define.Get(bc, "/todo", ListTodo).Description("List")
	define.StringEnum[TodoInputSort](bc, "id", "-id")

	ctx := context.Background()
	if err := define.EmitDoc(ctx, bc); err != nil {
		log.Fatalf("!! %+v", err)
	}
}
