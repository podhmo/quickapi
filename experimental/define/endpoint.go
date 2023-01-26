package define

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/qbind"
	"github.com/podhmo/quickapi/shared"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

type EndpointModifier reflectopenapi.RegisterFuncAction

func Method[I any, O any](bc *BuildContext, method, path string, action quickapi.Action[I, O], middlewares ...func(http.Handler) http.Handler) *EndpointModifier {
	h := quickapi.Lift(action)
	m := bc.m

	if bc.c.Loaded {
		op := findOpenapi3Operation(m.Doc, method, path)
		if op == nil {
			if !shared.FORCE {
				panic(fmt.Sprintf("path not found: %s %s", method, path))
			}
		} else {
			middleware := bc.mb.BuildMiddleware(path, op)
			middlewares = append([]func(http.Handler) http.Handler{middleware}, middlewares...)
		}
		bc.r.With(middlewares...).Method(method, path, h)
	}

	// if c.Loaded is true, this thunk is ignored.
	return (*EndpointModifier)(m.RegisterFunc(action).After(func(op *openapi3.Operation) {
		m.Doc.AddOperation(path, method, op)
		middleware := bc.mb.BuildMiddleware(path, op)
		middlewares := append([]func(http.Handler) http.Handler{middleware}, middlewares...)
		bc.r.With(middlewares...).Method(method, path, h)
	}))
}

func Get[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O], middlewares ...func(http.Handler) http.Handler) *EndpointModifier {
	return Method(bc, "GET", path, action, middlewares...)
}
func Post[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O], middlewares ...func(http.Handler) http.Handler) *EndpointModifier {
	return Method(bc, "POST", path, action, middlewares...)
}
func Put[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O], middlewares ...func(http.Handler) http.Handler) *EndpointModifier {
	return Method(bc, "PUT", path, action, middlewares...)
}
func Patch[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O], middlewares ...func(http.Handler) http.Handler) *EndpointModifier {
	return Method(bc, "PATCH", path, action, middlewares...)
}
func Delete[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O], middlewares ...func(http.Handler) http.Handler) *EndpointModifier {
	return Method(bc, "DELETE", path, action, middlewares...)
}
func Head[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O], middlewares ...func(http.Handler) http.Handler) *EndpointModifier {
	return Method(bc, "HEAD", path, action, middlewares...)
}
func Options[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O], middlewares ...func(http.Handler) http.Handler) *EndpointModifier {
	return Method(bc, "OPTIONS", path, action, middlewares...)
}

func (m *EndpointModifier) After(f func(op *openapi3.Operation)) *EndpointModifier {
	return (*EndpointModifier)((*reflectopenapi.RegisterFuncAction)(m).After(f))
}

func (m *EndpointModifier) OperationID(operationID string) *EndpointModifier {
	return m.After(func(op *openapi3.Operation) {
		op.OperationID = strings.TrimSpace(operationID)
	})
}

func (m *EndpointModifier) Description(description string) *EndpointModifier {
	return m.After(func(op *openapi3.Operation) {
		op.Description = strings.TrimSpace(description)
	})
}
func (m *EndpointModifier) Status(code int) *EndpointModifier {
	return m.After(func(op *openapi3.Operation) {
		def, ok := op.Responses["200"]
		if ok {
			delete(op.Responses, "200")
			op.Responses[strconv.Itoa(code)] = def
		}
	})
}
func (m *EndpointModifier) AnotherError(bc *BuildContext, code int, typ interface{}, description string) *EndpointModifier {
	return m.After(func(op *openapi3.Operation) {
		ref := bc.m.Visitor.VisitType(typ)
		val := openapi3.NewResponse().WithDescription(description).WithJSONSchemaRef(ref)
		op.Responses[strconv.Itoa(code)] = &openapi3.ResponseRef{Value: val}
	})
}
func (a *EndpointModifier) Example(code int, title string, value interface{}) *EndpointModifier {
	fn := (*reflectopenapi.RegisterFuncAction)(a)
	return (*EndpointModifier)(fn.Example(code, "application/json", title, value))
}

func GetHTML[I any](bc *BuildContext, path string, action quickapi.Action[I, string], dump quickapi.DumpFunc[string], middlewares ...func(http.Handler) http.Handler) *EndpointModifier {
	metadata := qbind.Scan(action)
	h := &quickapi.LiftedHandler[I, string]{
		Action:   action,
		Metadata: metadata,
		Dump:     dump,
	}
	m := bc.m
	method := "GET"

	if bc.c.Loaded {
		op := findOpenapi3Operation(m.Doc, method, path)
		if op == nil {
			if !shared.FORCE {
				panic(fmt.Sprintf("path not found: %s %s", method, path))
			}
		} else {
			middleware := bc.mb.BuildMiddleware(path, op)
			middlewares = append([]func(http.Handler) http.Handler{middleware}, middlewares...)
		}
		bc.r.With(middlewares...).Method(method, path, h)
	}

	// if c.Loaded is true, this thunk is ignored.
	return (*EndpointModifier)(m.RegisterFunc(action).After(func(op *openapi3.Operation) {
		// overwrite response/200/content/{application-json -> text/html}
		res := op.Responses.Get(200).Value
		res.Content = openapi3.NewContentWithSchemaRef(res.Content.Get("application/json").Schema, []string{"text/html"})

		m.Doc.AddOperation(path, method, op)
		middleware := bc.mb.BuildMiddleware(path, op)
		middlewares := append([]func(http.Handler) http.Handler{middleware}, middlewares...)
		bc.r.With(middlewares...).Method(method, path, h)
	}))
}

func findOpenapi3Operation(doc *openapi3.T, method, path string) *openapi3.Operation {
	pathItem := doc.Paths.Find(path)
	if pathItem == nil {
		return nil
	}

	switch strings.ToUpper(method) {
	case http.MethodGet:
		return pathItem.Get
	case http.MethodHead:
		return pathItem.Head
	case http.MethodPost:
		return pathItem.Post
	case http.MethodPut:
		return pathItem.Put
	case http.MethodPatch:
		return pathItem.Patch
	case http.MethodDelete:
		return pathItem.Delete
	case http.MethodConnect:
		return pathItem.Connect
	case http.MethodOptions:
		return pathItem.Options
	case http.MethodTrace:
		return pathItem.Trace
	default:
		return nil
	}
}
