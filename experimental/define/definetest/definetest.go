package definetest

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/podhmo/or"
	"github.com/podhmo/quickapi/experimental/define"
	"github.com/podhmo/quickapi/quickapitest"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

func NewHandler(t *testing.T, options ...func(*define.BuildContext)) http.Handler {
	t.Helper()
	ctx := quickapitest.NewContext(t)
	bc := or.Fatal(define.NewBuildContext(define.Doc(nil), chi.NewRouter(), func(c *reflectopenapi.Config) {
		c.SkipExtractComments = true
	}))(t)

	for _, opt := range options {
		opt(bc)
	}
	return or.Fatal(bc.BuildHandler(ctx))(t)
}
