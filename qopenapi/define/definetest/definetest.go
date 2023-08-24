package definetest

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/podhmo/or"
	"github.com/podhmo/quickapi/qopenapi/define"
	"github.com/podhmo/quickapi/quickapitest"
)

func NewHandler(t *testing.T, options ...func(*define.BuildContext)) http.Handler {
	t.Helper()
	ctx := quickapitest.NewContext(t)
	bc := or.Fatal(define.NewBuildContext(define.Doc(), chi.NewRouter(), func(c *define.Config) {
		c.ReflectOpenAPI.SkipExtractComments = true
		c.Validation.EnableResponseValidation = true
	}))(t)

	for _, opt := range options {
		opt(bc)
	}
	return or.Fatal(bc.BuildHandler(ctx))(t)
}
