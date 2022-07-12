package quickapitest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestDecodeResponse(t *testing.T) {
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

	got := DecodeResponse[person](t, rec.Result(), http.StatusOK)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("DecodeResponse(), path=%q mismatch (-want +got):\n%s", path, diff)
	}
}

func TestDoRequest(t *testing.T) {
	want := person{Name: "foo", Age: 20}

	path := "/users/1"
	req := httptest.NewRequest("GET", path, nil)
	got := DoRequest[person](
		t, req, http.StatusOK,
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if err := json.NewEncoder(w).Encode(want); err != nil {
				t.Fatalf("unexpected error (json.Encode): %+v", err)
			}
		}))
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("DoRequest(), path=%q mismatch (-want +got):\n%s", path, diff)
	}
}
