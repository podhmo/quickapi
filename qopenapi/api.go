package main

import (
	"context"
	_ "embed"
	"log"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

type Action[I any, O any] func(context.Context, I) (O, error)

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

////////////////////////////////////////
type APIError struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

//go:embed skeleton.json
var docSkeleton []byte

func main() {
	doc, err := reflectopenapi.NewDocFromSkeleton(docSkeleton)
	if err != nil {
		log.Fatalf("!! %+v", err)
	}

	c := reflectopenapi.Config{
		Doc:          doc,
		DefaultError: APIError{},
		StrictSchema: true,
		IsRequiredCheckFunction: func(tag reflect.StructTag) bool {
			ok := true
			if _, isOptional := tag.Lookup("optional"); isOptional {
				ok = false
			}
			return ok
		},
	}
	c.EmitDoc(func(m *reflectopenapi.Manager) {
		m.RegisterFunc(ListTodo).After(func(op *openapi3.Operation) {
			m.Doc.AddOperation("/Todo", "GET", op)
		})
	})
}
