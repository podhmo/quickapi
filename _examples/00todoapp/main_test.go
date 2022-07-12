package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/podhmo/quickapi/quickapitest"
)

type TodoListOutput struct {
	Items []Todo `json:"items"`
}

func TestListTodo(t *testing.T) {
	router := chi.NewRouter()
	mount(router)

	path := "/todos?sort=-id"
	req := httptest.NewRequest("GET", path, nil)
	got := quickapitest.DoRequest[TodoListOutput](t, req, http.StatusOK, router)

	want := TodoListOutput{Items: []Todo{
		{ID: 3, Title: "byebye", Done: false},
		{ID: 1, Title: "hello", Done: false},
	}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("GET %q mismatch (-want +got):\n%s", path, diff)
	}
}
