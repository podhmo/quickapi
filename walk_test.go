package quickapi_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/podhmo/quickapi"
)

func TestValidatePathVars(t *testing.T) {
	type testcase struct {
		msg     string
		path    string
		handler http.HandlerFunc
	}

	t.Run("ok", func(t *testing.T) {
		cases := []testcase{
			{"one", "/foo/{id}", quickapi.Lift(func(ctx context.Context, input struct {
				ID int `in:"path" path:"id"`
			}) (interface{}, error) {
				return nil, nil
			})},
			{"two", "/foo/{foo_id}/bar/{bar_id}", quickapi.Lift(func(ctx context.Context, input struct {
				FooID int `in:"path" path:"foo_id"`
				BarID int `in:"path" path:"bar_id"`
			}) (interface{}, error) {
				return nil, nil
			})},
			{"regexp", "/articles/{rid:^[0-9]{5,6}}", quickapi.Lift(func(ctx context.Context, input struct {
				RID int `in:"path" path:"rid"`
			}) (interface{}, error) {
				return nil, nil
			})},
		}

		for _, c := range cases {
			c := c
			t.Run(c.msg, func(t *testing.T) {
				r := chi.NewRouter()
				r.Get(c.path, c.handler)

				err := quickapi.WalkRoute(r, func(item quickapi.RouteItem) error {
					return item.ValidatePathVars()
				})
				if err != nil {
					t.Errorf("ValidatePathVars(): unexpected error %+v", err)
				}
			})
		}
	})

	t.Run("ng", func(t *testing.T) {
		cases := []testcase{
			{"missing", "/foo/{id}", quickapi.Lift(func(ctx context.Context, input struct {
			}) (interface{}, error) {
				return nil, nil
			})},
			{"conflict", "/foo/{id}", quickapi.Lift(func(ctx context.Context, input struct {
				ID string `in:"path" path:"_id"`
			}) (interface{}, error) {
				return nil, nil
			})},
			// TODO: type conflict
		}

		for _, c := range cases {
			c := c
			t.Run(c.msg, func(t *testing.T) {
				r := chi.NewRouter()
				r.Get(c.path, c.handler)

				err := quickapi.WalkRoute(r, func(item quickapi.RouteItem) error {
					return item.ValidatePathVars()
				})
				if err == nil {
					t.Errorf("ValidatePathVars(): expect error, but nil")
				}
				t.Logf("validation error: %+v", err)
			})
		}
	})
}
