package quickapi_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/quickapitest"
)

type ints struct{ XS []int }

func TestLift_OK(t *testing.T) {
	action := func(context.Context, quickapi.Empty) ([]int, error) { return []int{1, 2, 3}, nil }
	handler := quickapi.Lift(action)
	req := httptest.NewRequest("GET", "/", nil)

	got := quickapitest.DoRequest[[]int](t, req, 200, handler)
	want := []int{1, 2, 3}
	if diff := cmp.Diff(ints{want}, ints{got}); diff != "" {
		t.Errorf("Lift() mismatch (-want +got):\n%s", diff)
	}
}

func TestLift_OK_NilAsEmptySlice(t *testing.T) {
	action := func(context.Context, quickapi.Empty) ([]int, error) { return nil, nil }
	handler := quickapi.Lift(action)
	req := httptest.NewRequest("GET", "/", nil)

	got := quickapitest.DoRequest[[]int](t, req, 200, handler)
	want := []int{}
	if diff := cmp.Diff(ints{want}, ints{got}); diff != "" {
		t.Errorf("Lift() mismatch (-want +got):\n%s", diff)
	}
}

func TestLift_OK_WithDefault(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}

	action := func(ctx context.Context, i person) (person, error) { return i, nil }
	handler := quickapi.NewHandler(action)
	handler.Default = func() person { return person{Name: "foo"} }

	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"age": 20}`))

	got := quickapitest.DoRequest[person](t, req, 200, handler)
	want := person{Name: "foo", Age: 20}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Lift() mismatch (-want +got):\n%s", diff)
	}
}

func TestLift_NotFound(t *testing.T) {
	code := 404
	action := func(context.Context, quickapi.Empty) ([]int, error) {
		return nil, quickapi.NewAPIError(fmt.Errorf("hmm"), code)
	}

	handler := quickapi.Lift(action)
	req := httptest.NewRequest("GET", "/", nil)

	got := quickapitest.DoRequest[quickapi.ErrorResponse](t, req, code, handler)
	want := quickapi.ErrorResponse{Code: code, Error: "api-error: hmm"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Lift() mismatch (-want +got):\n%s", diff)
	}
}

type Ill struct {
	Name string `json:"name"`
}

func (ob Ill) Validate(ctx context.Context) error {
	return quickapi.NewAPIError(fmt.Errorf("ill"), http.StatusUnprocessableEntity)
}

func TestLift_UnprocessableEntity_withValidation(t *testing.T) {
	code := 422
	action := func(ctx context.Context, input Ill) ([]int, error) {
		return nil, nil
	}

	handler := quickapi.Lift(action)
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name": "foo"}`))

	got := quickapitest.DoRequest[quickapi.ErrorResponse](t, req, code, handler)
	want := quickapi.ErrorResponse{Code: code, Error: "api-error: ill"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Lift() mismatch (-want +got):\n%s", diff)
	}
}

func TestLift_Found_Redirect(t *testing.T) {
	code := 302
	action := func(ctx context.Context, input quickapi.Empty) ([]int, error) {
		return nil, quickapi.Redirect(http.StatusFound, "http://example.net")
	}

	handler := quickapi.Lift(action)
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	res := rec.Result()

	if want, got := code, res.StatusCode; want != got {
		t.Errorf("Lift() status-code, want=%d != got=%d", want, got)
	}
	if want, got := "http://example.net", res.Header.Get("Location"); want != got {
		t.Errorf("Lift() header location, want=%q != got=%q", want, got)
	}
}

func TestLift_NoContent(t *testing.T) {
	code := 204
	action := func(ctx context.Context, input quickapi.Empty) (any, error) {
		return quickapi.NoContent(204), nil
	}

	handler := quickapi.Lift(action)
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	res := rec.Result()

	if want, got := code, res.StatusCode; want != got {
		t.Errorf("Lift() status-code, want=%d != got=%d", want, got)
	}
}

func TestLift_BindPathVars(t *testing.T) {
	type Input struct {
		ID int `path:"id"`
	}
	type Output struct {
		InputID int `json:"inputId"`
	}
	action := func(ctx context.Context, input Input) (Output, error) { return Output{InputID: input.ID}, nil }
	r := chi.NewRouter() // need chi.RouteContext for pathvar binding
	r.Get("/{id}", quickapi.Lift(action))

	t.Run("200", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/10", nil)
		got := quickapitest.DoRequest[Output](t, req, 200, r)

		want := Output{InputID: 10}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("data, want=%#+v, but got=%#+v", want, got)
		}
	})

	t.Run("404", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/foo", nil)
		quickapitest.DoRequest[quickapi.ErrorResponse](t, req, 404, r)
	})
}

func TestLift_BindData(t *testing.T) {
	type Input struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	type Output struct {
		Message string `json:"message"`
	}
	action := func(ctx context.Context, input Input) (Output, error) {
		return Output{Message: fmt.Sprintf("%s(%d): hello", input.Name, input.Age)}, nil
	}
	h := quickapi.Lift(action)

	t.Run("200", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name": "foo", "age": 20}`))
		got := quickapitest.DoRequest[Output](t, req, 200, h)

		want := Output{Message: "foo(20): hello"}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("data, want=%#+v, but got=%#+v", want, got)
		}
	})

	t.Run("400-no-body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", nil)
		quickapitest.DoRequest[quickapi.ErrorResponse](t, req, 400, h)
	})
	t.Run("400-invalid-type", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/foo", strings.NewReader(`{"name": "foo", "age": "20"}`))
		quickapitest.DoRequest[Output](t, req, 400, h)
	})
}

func TestLift_BindQueryVars(t *testing.T) {
	type Input struct {
		ID   int    `path:"id"`
		Sort string `query:"sort"`
	}
	type Output struct {
		InputID   int    `json:"inputId"`
		InputSort string `json:"inputSort"`
	}
	action := func(ctx context.Context, input Input) (Output, error) {
		return Output{InputID: input.ID, InputSort: input.Sort}, nil
	}
	r := chi.NewRouter() // need chi.RouteContext for pathvar binding
	r.Get("/{id}", quickapi.Lift(action))

	t.Run("200", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/10", nil)
		got := quickapitest.DoRequest[Output](t, req, 200, r)

		want := Output{InputID: 10}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("data, want=%#+v, but got=%#+v", want, got)
		}
	})
	t.Run("200-with-query", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/10?sort=-id", nil)
		got := quickapitest.DoRequest[Output](t, req, 200, r)

		want := Output{InputID: 10, InputSort: "-id"}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("data, want=%#+v, but got=%#+v", want, got)
		}
	})
}
