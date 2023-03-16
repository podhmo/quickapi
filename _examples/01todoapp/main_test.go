package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/quickapitest"
)

type TodoListOutput struct {
	Items []Todo `json:"items"`
}

func TestListTodo(t *testing.T) {
	router := chi.NewRouter()
	mount(router)

	req := httptest.NewRequest("GET", "/todos?sort=-id", nil)
	got := quickapitest.DoRequest[TodoListOutput](t, router, req, http.StatusOK)

	want := TodoListOutput{Items: []Todo{
		{ID: 3, Title: "byebye", Done: false},
		{ID: 1, Title: "hello", Done: false},
	}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("response mismatch (-want +got):\n%s", diff)
	}
}

func TestGetTodo(t *testing.T) {
	router := chi.NewRouter()
	mount(router)

	t.Run("ok", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/todos/1", nil)
		got := quickapitest.DoRequest[Todo](t, router, req, http.StatusOK)

		want := Todo{ID: 1, Title: "hello", Done: false}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("response mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("ng", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/todos/10", nil)
		got := quickapitest.DoRequest[quickapi.ErrorResponse](t, router, req, http.StatusNotFound)

		want := quickapi.ErrorResponse{Code: 404, Error: "api-error: not found"}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("response mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestMethodNotAllowed(t *testing.T) {
	router := chi.NewRouter()
	mount(router)

	req := httptest.NewRequest("PUT", "/todos", nil)
	quickapitest.DoRequest[any](t, req, http.StatusMethodNotAllowed, router)
}
