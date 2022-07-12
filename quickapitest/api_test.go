package quickapitest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDecodeResponse(t *testing.T) {
	type person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	want := person{Name: "foo", Age: 20}
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("unexpected error (json.Encode): %+v", err)
		}
	})

	path := "/users/1"
	req := httptest.NewRequest("GET", path, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	got := DecodeResponse[person](t, path, rec.Result(), http.StatusOK)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("DecodeResponse() mismatch (-want +got):\n%s", diff)
	}
}
