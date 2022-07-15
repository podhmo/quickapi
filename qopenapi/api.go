package qopenapi

import (
	"context"
	_ "embed"
	"encoding/json"
	"os"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/podhmo/quickapi"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

type APIError struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

//go:embed skeleton.json
var docSkeleton []byte

type Router struct {
	m *reflectopenapi.Manager
	c *reflectopenapi.Config

	r chi.Router

	commit func(context.Context) error
}

func NewRouter() (*Router, error) {
	doc, err := reflectopenapi.NewDocFromSkeleton(docSkeleton)
	if err != nil {
		return nil, err
	}

	c := reflectopenapi.Config{
		Doc:          doc,
		DefaultError: APIError{},
		StrictSchema: true,
		IsRequiredCheckFunction: func(tag reflect.StructTag) bool {
			required := true
			if val, ok := tag.Lookup("openapi"); ok && val != "body" {
				required = false
			}
			if _, isOptional := tag.Lookup("optional"); isOptional {
				required = false
			}
			return required
		},
	}

	m, commit, err := c.NewManager()
	if err != nil {
		return nil, err
	}
	r := &Router{r: chi.NewRouter(), c: &c, m: m, commit: commit}
	return r, nil
}

func EmitDoc(ctx context.Context, r *Router) error {
	if err := r.commit(ctx); err != nil {
		return err
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(r.m.Doc); err != nil {
		return err
	}
	return nil
}

func Method[I any, O any](r *Router, method, path string, action quickapi.Action[I, O]) *reflectopenapi.RegisterFuncAction {
	h := quickapi.Lift(action)
	r.r.Method(method, path, h)
	m := r.m
	return m.RegisterFunc(action).After(func(op *openapi3.Operation) {
		m.Doc.AddOperation(path, method, op)
	})
}

func Get[I any, O any](r *Router, path string, action quickapi.Action[I, O]) *reflectopenapi.RegisterFuncAction {
	return Method(r, "GET", path, action)
}
func Post[I any, O any](r *Router, path string, action quickapi.Action[I, O]) *reflectopenapi.RegisterFuncAction {
	return Method(r, "POST", path, action)
}
func Put[I any, O any](r *Router, path string, action quickapi.Action[I, O]) *reflectopenapi.RegisterFuncAction {
	return Method(r, "PUT", path, action)
}
func Patch[I any, O any](r *Router, path string, action quickapi.Action[I, O]) *reflectopenapi.RegisterFuncAction {
	return Method(r, "PATCH", path, action)
}
func Delete[I any, O any](r *Router, path string, action quickapi.Action[I, O]) *reflectopenapi.RegisterFuncAction {
	return Method(r, "DELETE", path, action)
}
func Head[I any, O any](r *Router, path string, action quickapi.Action[I, O]) *reflectopenapi.RegisterFuncAction {
	return Method(r, "HEAD", path, action)
}
func Options[I any, O any](r *Router, path string, action quickapi.Action[I, O]) *reflectopenapi.RegisterFuncAction {
	return Method(r, "OPTIONS", path, action)
}
