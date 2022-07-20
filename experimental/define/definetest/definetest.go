package definetest

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/podhmo/or"
	"github.com/podhmo/quickapi/experimental/define"
	"github.com/podhmo/quickapi/quickapitest"
	"github.com/podhmo/quickapi/shared"
)

func NewHandler(t *testing.T, options ...func(*define.BuildContext)) http.Handler {
	t.Helper()
	ctx := quickapitest.NewContext(t)
	bc := or.Fatal(define.NewBuildContext(define.Doc(), chi.NewRouter()))(t)

	// skip extract comments
	m := bc.ReflectOpenAPIManager()
	m.Visitor.CommentLookup = nil

	// silent logger
	bc.Router().Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			req = req.WithContext(shared.SetLogger(req.Context(), shared.GetLogger(ctx)))
			next.ServeHTTP(w, req)
		})
	})

	for _, opt := range options {
		opt(bc)
	}
	return or.Fatal(bc.BuildHandler(ctx))(t)
}
