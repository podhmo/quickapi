package quickapi_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/podhmo/quickapi"
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

func TestLift_NotFound(t *testing.T) {
	code := 404
	action := func(context.Context, quickapi.Empty) ([]int, error) {
		return nil, quickapi.NewAPIError(fmt.Errorf("hmm"), code)
	}

	handler := quickapi.Lift(action)
	req := httptest.NewRequest("GET", "/", nil)

	type errorResponse struct {
		Code  int    `json:"code"`
		Error string `json:"error"`
	}

	got := quickapitest.DoRequest[errorResponse](t, req, code, handler)
	want := errorResponse{Code: code, Error: "api-error: hmm"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("FillNil() mismatch (-want +got):\n%s", diff)
	}
}
