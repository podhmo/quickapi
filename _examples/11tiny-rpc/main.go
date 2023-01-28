package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/qopenapi/define"
	"github.com/podhmo/quickapi/qopenapi/definerpc"
)

func main() {
	doc := define.Doc().Title("tinyrpc")
	r := quickapi.DefaultRouter()
	bc := define.MustBuildContext(doc, r)

	definerpc.Action(bc, Hello)

	ctx := context.Background()
	if err := bc.EmitDoc(ctx, os.Stdout); err != nil {
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
