package definetest_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/podhmo/quickapi/experimental/define"
	"github.com/podhmo/quickapi/experimental/define/definetest"
	"github.com/podhmo/quickapi/quickapitest"
	"github.com/podhmo/quickapi/shared"
)

type Input struct {
	Name string `json:"name"`
}
type Output struct {
	Name string `json:"name"`
}

type ref struct {
	V []Output
}

func TestIt(t *testing.T) {
	items := []Output{{Name: "foo"}, {Name: "bar"}}

	handler := definetest.NewHandler(t, func(bc *define.BuildContext) {
		define.Get(bc, "/", func(context.Context, shared.Empty) ([]Output, error) {
			return items, nil
		})

		define.Post(bc, "/", func(ctx context.Context, input Input) (any, error) {
			items = append(items, Output(input))
			return shared.NoContent(http.StatusCreated), nil
		})
	})

	t.Run("GET", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		got := quickapitest.DoRequest[[]Output](t, req, 200, handler)

		want := []Output{{Name: "foo"}, {Name: "bar"}}
		if diff := cmp.Diff(ref{want}, ref{got}); diff != "" {
			t.Errorf("%s %s, response mismatch (-want +got):\n%s", req.Method, req.URL.Path, diff)
		}
	})

	t.Run("POST", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name": "moo"}`))
		quickapitest.DoRequest[any](t, req, 201, handler)
	})

	t.Run("POST-manually", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name": "moo2"}`))
		req.Header.Set("Content-type", "application/json")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if want, got := 201, rec.Result().StatusCode; want != got {
			t.Errorf("%s %s, status code: want=%d != got=%d", req.Method, req.URL.Path, want, got)
		}
	})

	// db check
	want := []Output{{Name: "foo"}, {Name: "bar"}, {Name: "moo"}, {Name: "moo2"}}
	got := items
	if diff := cmp.Diff(ref{want}, ref{got}); diff != "" {
		t.Errorf("db mismatch (-want +got):\n%s", diff)
	}
}
