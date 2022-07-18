package quickapi_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/qdump"
	"github.com/podhmo/quickapi/quickapitest"
)

func TestLift_OK(t *testing.T) {
	action := func(context.Context, quickapi.Empty) ([]int, error) { return []int{1, 2, 3}, nil }
	handler := quickapi.Lift(action)
	req := httptest.NewRequest("GET", "/", nil)

	got := quickapitest.DoRequest[[]int](t, req, 200, handler)
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("data, want=%#+v, but got=%#+v", want, got)
	}
}

func TestLift_OK_NilAsEmptySlice(t *testing.T) {
	action := func(context.Context, quickapi.Empty) ([]int, error) { return nil, nil }
	handler := quickapi.Lift(action)
	req := httptest.NewRequest("GET", "/", nil)

	got := quickapitest.DoRequest[[]int](t, req, 200, handler)
	want := []int{}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("data, want=%#+v, but got=%#+v", want, got)
	}
}

type errorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func TestLift_NotFound(t *testing.T) {
	code := 404
	action := func(context.Context, quickapi.Empty) ([]int, error) {
		return nil, quickapi.NewAPIError(fmt.Errorf("hmm"), code)
	}

	handler := quickapi.Lift(action)
	req := httptest.NewRequest("GET", "/", nil)

	got := quickapitest.DoRequest[errorResponse](t, req, code, handler)
	want := errorResponse{Code: code, Error: "api-error: hmm"}
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
	req := httptest.NewRequest("GET", "/", strings.NewReader(`{"name": "foo"}`))

	got := quickapitest.DoRequest[errorResponse](t, req, code, handler)
	want := errorResponse{Code: code, Error: "api-error: ill"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Lift() mismatch (-want +got):\n%s", diff)
	}
}

func TestLift_Found_Redirect(t *testing.T) {
	code := 302
	action := func(ctx context.Context, input quickapi.Empty) ([]int, error) {
		return nil, qdump.Redirect(http.StatusFound, "http://example.net")
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
