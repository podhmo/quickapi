package quickapi_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/podhmo/quickapi"
)

func doRequest[I any, O any](t *testing.T, action quickapi.Action[I, O], wantCode int) O {
	t.Helper()

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler := quickapi.Lift(action)
	handler(rec, req)

	gotCode := rec.Result().StatusCode
	if want, got := wantCode, gotCode; want != got {
		t.Errorf("status-code, want=%d, but got=%d", want, got)
	}

	var got O
	if gotCode == 200 {
		if err := json.NewDecoder(rec.Result().Body).Decode(&got); err != nil {
			t.Errorf("unexpected error (decode): %+v", err)
		}
	} else {
		buf := new(strings.Builder)
		io.Copy(buf, rec.Result().Body)
		t.Logf("response: %s", buf.String())
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

func TestLift_NotFound(t *testing.T) {
	code := 404
	action := func(context.Context, quickapi.Empty) ([]int, error) {
		return nil, quickapi.NewAPIError(fmt.Errorf("hmm"), code)
	}
	doRequest(t, action, code)
}
