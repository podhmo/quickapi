package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/quickapitest"
)

func TestOK(t *testing.T) {
	h := quickapi.Lift(Hello)

	req := httptest.NewRequest("GET", "/", nil)
	got := quickapitest.DoRequest[HelloOutput](t, req, http.StatusOK, h)

	want := HelloOutput{Hello: "world"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("data, want=%#+v, but got=%#+v", want, got)
	}
}
