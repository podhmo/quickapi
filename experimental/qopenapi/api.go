package qopenapi

import (
	"context"
	_ "embed"
	"encoding/json"
	"os"
	"reflect"
	"strings"

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

type BuildContext struct {
	m *reflectopenapi.Manager
	c *reflectopenapi.Config

	r chi.Router

	commit func(context.Context) error
}

func New(r chi.Router) (*BuildContext, error) {
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
			if val, ok := tag.Lookup("json"); ok && strings.Contains(val, "omitempty") {
				required = false
			}
			return required
		},
	}

	m, commit, err := c.NewManager()
	if err != nil {
		return nil, err
	}
	return &BuildContext{r: r, c: &c, m: m, commit: commit}, nil
}

func EmitDoc(ctx context.Context, bc *BuildContext) error {
	if err := bc.commit(ctx); err != nil {
		return err
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(bc.m.Doc); err != nil {
		return err
	}
	return nil
}

func Method[I any, O any](bc *BuildContext, method, path string, action quickapi.Action[I, O]) *reflectopenapi.RegisterFuncAction {
	h := quickapi.Lift(action)
	bc.r.Method(method, path, h)
	m := bc.m
	return m.RegisterFunc(action).After(func(op *openapi3.Operation) {
		m.Doc.AddOperation(path, method, op)
	})
}

func Get[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O]) *reflectopenapi.RegisterFuncAction {
	return Method(bc, "GET", path, action)
}
func Post[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O]) *reflectopenapi.RegisterFuncAction {
	return Method(bc, "POST", path, action)
}
func Put[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O]) *reflectopenapi.RegisterFuncAction {
	return Method(bc, "PUT", path, action)
}
func Patch[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O]) *reflectopenapi.RegisterFuncAction {
	return Method(bc, "PATCH", path, action)
}
func Delete[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O]) *reflectopenapi.RegisterFuncAction {
	return Method(bc, "DELETE", path, action)
}
func Head[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O]) *reflectopenapi.RegisterFuncAction {
	return Method(bc, "HEAD", path, action)
}
func Options[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O]) *reflectopenapi.RegisterFuncAction {
	return Method(bc, "OPTIONS", path, action)
}

func DefineType(bc *BuildContext, typ interface{}) *reflectopenapi.RegisterTypeAction {
	return bc.m.RegisterType(typ)
}

func DefineEnum[T any](bc *BuildContext, defaultValue T, values ...T) *reflectopenapi.RegisterTypeAction {
	dst := make([]interface{}, len(values)+1)
	typedValue := T(defaultValue)
	dst[0] = typedValue
	for i, v := range values {
		dst[i+1] = T(v)
	}
	return bc.m.RegisterType(typedValue, func(ref *openapi3.Schema) {
		ref.Default = dst[0]
		ref.Enum = dst
	})
}
func DefineStringEnum[T ~string](bc *BuildContext, defaultValue T, values ...T) *reflectopenapi.RegisterTypeAction {
	return DefineEnum(bc, defaultValue, values...)
}
func DefineIntEnum[T ~int](bc *BuildContext, defaultValue T, values ...T) *reflectopenapi.RegisterTypeAction {
	return DefineEnum(bc, defaultValue, values...)
}
