package quickapi_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/podhmo/quickapi"
)

func TestLift_Empty(t *testing.T) {
	action := func(context.Context, quickapi.Empty) ([]int, error) { return nil, nil }
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	quickapi.Lift(action)(rec, req)

	if want, got := 200, rec.Result().StatusCode; want != got {
		t.Errorf("status, want=%d, but got=%d", want, got)
	}
	var got []int
	if err := json.NewDecoder(rec.Result().Body).Decode(&got); err != nil {
		t.Errorf("unexpected error (decode): %+v", err)
	}
	if got == nil {
		t.Errorf("status, want=%#+v, but got=%#+v", nil, got)
	}
}
