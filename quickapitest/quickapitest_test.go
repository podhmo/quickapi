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

func TestDoRequest(t *testing.T) {
	want := person{Name: "foo", Age: 20}

	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("unexpected error (json.Encode): %+v", err)
		}
	})

	path := "/users/1"
	req := httptest.NewRequest("GET", path, nil)
	got := DoHandler[person](t, h, req, http.StatusOK)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("DoRequest(), path=%q mismatch (-want +got):\n%s", path, diff)
	}
}
