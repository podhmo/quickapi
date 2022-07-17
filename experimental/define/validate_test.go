package define_test

import (
	"context"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/podhmo/or"
	"github.com/podhmo/quickapi/experimental/define"
)

type Todo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func GetTodo(ctx context.Context, input struct {
	ID string `path:"id" openapi:"path"`
}) (output Todo, err error) {
	output.ID = input.ID
	return
}

func TestValidate(t *testing.T) {
	ctx := context.Background()

	bc := or.Fatal(define.NewBuildContext(define.Doc(), chi.NewRouter()))(t)
	define.Get(bc, "/foo/{id}", GetTodo)
	h := or.Fatal(bc.BuildHandler(ctx))(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/foo/1", nil)
	h.ServeHTTP(rec, req)

	httputil.DumpResponse(rec.Result(), true)
}
