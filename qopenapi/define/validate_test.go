package define_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "embed"

	"github.com/google/go-cmp/cmp"
	"github.com/podhmo/or"
	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/qopenapi/define"
	"github.com/podhmo/quickapi/qopenapi/define/definetest"
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

//go:embed testdata/post-400-response.json
var ngresponseBody []byte

func TestIt(t *testing.T) {
	items := []Output{{Name: "foo"}, {Name: "bar"}}

	handler := definetest.NewHandler(t, func(bc *define.BuildContext) {
		define.Get(bc, "/", func(context.Context, shared.Empty) ([]Output, error) {
			return items, nil
		})

		define.Post(bc, "/", func(ctx context.Context, input Input) (any, error) {
			items = append(items, Output(input))
			return shared.NoContent(http.StatusCreated), nil
		}).Status(http.StatusCreated)
	})

	t.Run("GET", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		wantCode := 200
		got := quickapitest.DoHandler[[]Output](t, handler, req, wantCode)

		want := []Output{{Name: "foo"}, {Name: "bar"}}
		if diff := cmp.Diff(ref{want}, ref{got}); diff != "" {
			t.Errorf("%s %s, response mismatch (-want +got):\n%s", req.Method, req.URL.Path, diff)
		}
	})

	t.Run("POST", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name": "moo"}`))
		wantCode := 201
		quickapitest.DoHandler[any](t, handler, req, wantCode)
	})

	t.Run("POST-invalid", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", strings.NewReader(`{}`))
		wantCode := 400
		got := quickapitest.DoHandler[quickapi.ErrorResponse](t, handler, req, wantCode)

		var want quickapi.ErrorResponse
		or.Fatal(t, json.NewDecoder(bytes.NewBuffer(ngresponseBody)).Decode(&want))(t)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("%s %s, response mismatch (-want +got):\n%s", req.Method, req.URL.Path, diff)
		}
	})

	t.Run("POST-manually", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name": "moo2"}`))
		req.Header.Set("Content-type", "application/json")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		wantCode := 201
		if want, got := wantCode, rec.Result().StatusCode; want != got {
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
