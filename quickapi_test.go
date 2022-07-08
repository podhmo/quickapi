package quickapi_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/podhmo/quickapi"
)

func doRequest[I any, O any](t *testing.T, action quickapi.Action[I, O], expectedStatus int) O {
	t.Helper()

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler := quickapi.Lift(action)
	handler(rec, req)

	if want, got := expectedStatus, rec.Result().StatusCode; want != got {
		t.Errorf("status, want=%d, but got=%d", want, got)
	}

	var got O
	if err := json.NewDecoder(rec.Result().Body).Decode(&got); err != nil {
		t.Errorf("unexpected error (decode): %+v", err)
	}
	return got
}

func TestLift_OK(t *testing.T) {
	action := func(context.Context, quickapi.Empty) ([]int, error) { return []int{1, 2, 3}, nil }
	got := doRequest(t, action, 200)

	if want := []int{1, 2, 3}; !reflect.DeepEqual(want, got) {
		t.Errorf("data, want=%#+v, but got=%#+v", want, got)
	}
}

func TestLift_OK_NilAsEmptySlice(t *testing.T) {
	action := func(context.Context, quickapi.Empty) ([]int, error) { return nil, nil }
	got := doRequest(t, action, 200)

	if want := []int{}; !reflect.DeepEqual(want, got) {
		t.Errorf("data, want=%#+v, but got=%#+v", want, got)
	}
}
