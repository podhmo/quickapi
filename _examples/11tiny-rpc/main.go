//go:generate go run ./ -gendoc -docfile openapi.json -mdfile apidoc.md
package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"go/token"
	"log"
	"time"

	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/qopenapi/define"
	"github.com/podhmo/quickapi/qopenapi/definerpc"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	reflectshape "github.com/podhmo/reflect-shape"
)

//go:embed openapi.json
var openapiDocData []byte

//go:embed apidoc.md
var mdDocData []byte

var options struct {
	gendoc  bool
	docfile string
	mdfile  string
	port    int
}

func main() {
	flag.BoolVar(&options.gendoc, "gendoc", false, "generate openapi doc")
	flag.IntVar(&options.port, "port", 8080, "port")
	flag.StringVar(&options.docfile, "docfile", "", "file name of openapi doc. if this value is empty output to stdout.")
	flag.StringVar(&options.mdfile, "mdfile", "", "")
	flag.Parse()
	if err := run(); err != nil {
		log.Fatalf("!! %+v", err)
	}

}

func run() error {
	ctx := context.Background()

	doc := define.Doc().
		Title("hello example")

	if !options.gendoc {
		doc = doc.LoadFromData(openapiDocData)
	}

	router := quickapi.DefaultRouter()
	bc, err := define.NewBuildContext(doc, router, func(c *reflectopenapi.Config) {
		c.GoPositionFunc = func(fset *token.FileSet, f *reflectshape.Func) string {
			// TODO: multiple package
			fpos := fset.Position(f.Pos())
			return fmt.Sprintf("https://github.com/podhmo/quickapi/blob/main/_examples/11tiny-rpc/main.go#L%d", fpos.Line)
		}
	})
	if err != nil {
		return err
	}

	definerpc.Action(bc, Hello)

	if options.gendoc {
		if err := bc.EmitDoc(ctx, options.docfile); err != nil {
			return err
		}

		if options.mdfile != "" {
			if err := bc.EmitMDDoc(ctx, options.mdfile); err != nil {
				return err
			}
		}
		return nil
	}

	handler, err := bc.BuildHandler(ctx)
	if err != nil {
		return err
	}
	dochandler, err := bc.BuildDocHandler(ctx, "/_doc", mdDocData)
	if err != nil {
		return err
	}
	bc.Router().Mount("/_doc", dochandler)

	if err := quickapi.NewServer(fmt.Sprintf(":%d", options.port), handler, 5*time.Second).ListenAndServe(ctx); err != nil {
		log.Printf("[Error] !! %+v", err)
	}
	return nil
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
