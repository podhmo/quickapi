package main

import (
	"context"
	"fmt"
	"log"

	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/qopenapi/define"
	"github.com/podhmo/quickapi/qopenapi/definerpc"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

func main() {
	doc := define.Doc().Title("tinyrpc")
	r := quickapi.DefaultRouter()
	bc := define.MustBuildContext(doc, r, func(c *reflectopenapi.Config) {
		c.DisableInputRef = true
		c.DisableOutputRef = true
	})

	definerpc.Action(bc, Hello)

	ctx := context.Background()
	if err := bc.EmitDoc(ctx, ""); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

type HelloInput struct {
	Name string `json:"name"`
}
type HelloOutput struct {
	Message string `json:"message"`
}

// hello, greeting message
func Hello(ctx context.Context, input HelloInput) (*HelloOutput, error) {
	return &HelloOutput{Message: fmt.Sprintf("hello %s", input.Name)}, nil
}
