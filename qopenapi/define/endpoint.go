package define

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/internal/pathutil"
	"github.com/podhmo/quickapi/shared"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

type EndpointModifier[I any, O any] struct {
	Handler  *quickapi.LiftedHandler[I, O]
	register *reflectopenapi.RegisterFuncAction
}

func Method[I any, O any](bc *BuildContext, method, path string, action quickapi.Action[I, O], middlewares ...func(http.Handler) http.Handler) *EndpointModifier[I, O] {
	h := quickapi.NewHandler(action)
	m := bc.m

	normalizedPath, _, pathvars := pathutil.NormalizeTemplatedPath(path)
	if bc.c.Loaded {
		op := findOpenapi3Operation(m.Doc, method, normalizedPath)
		if op == nil {
			msg := fmt.Sprintf("path not found: %s %s", method, path)
			if !shared.FORCE {
				panic(msg)
			} else {
				log.Printf("[INFO]  %s", msg)
			}
		} else {
			middleware := bc.mb.BuildMiddleware(path, op)
			middlewares = append([]func(http.Handler) http.Handler{middleware}, middlewares...)
		}
		bc.r.With(middlewares...).Method(method, path, h)
	}

	// if c.Loaded is true, this thunk is ignored.
	register := m.RegisterFunc(action).After(func(op *openapi3.Operation) {
		m.Doc.AddOperation(normalizedPath, method, op)

		// add pattern if type==string and regex is existed
		for _, p := range op.Parameters {
			if p.Value == nil {
				continue
			}
			if p.Value.In != "path" {
				continue
			}
			if regex, ok := pathvars[p.Value.Name]; ok && regex != "" {
				if value := p.Value.Schema.Value; value != nil && value.Type == "string" {
					if !strings.HasSuffix(regex, "$") {
						regex = regex + "$"
					}
					if !strings.HasPrefix(regex, "^") {
						regex = "^" + regex
					}
					value.Pattern = regex
				}
			}
		}

		middleware := bc.mb.BuildMiddleware(path, op)
		middlewares := append([]func(http.Handler) http.Handler{middleware}, middlewares...)
		bc.r.With(middlewares...).Method(method, path, h)
	})
	return &EndpointModifier[I, O]{
		Handler:  h,
		register: register,
	}
}

func Get[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O], middlewares ...func(http.Handler) http.Handler) *EndpointModifier[I, O] {
	return Method(bc, "GET", path, action, middlewares...)
}
func Post[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O], middlewares ...func(http.Handler) http.Handler) *EndpointModifier[I, O] {
	return Method(bc, "POST", path, action, middlewares...)
}
func Put[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O], middlewares ...func(http.Handler) http.Handler) *EndpointModifier[I, O] {
	return Method(bc, "PUT", path, action, middlewares...)
}
func Patch[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O], middlewares ...func(http.Handler) http.Handler) *EndpointModifier[I, O] {
	return Method(bc, "PATCH", path, action, middlewares...)
}
func Delete[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O], middlewares ...func(http.Handler) http.Handler) *EndpointModifier[I, O] {
	return Method(bc, "DELETE", path, action, middlewares...)
}
func Head[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O], middlewares ...func(http.Handler) http.Handler) *EndpointModifier[I, O] {
	return Method(bc, "HEAD", path, action, middlewares...)
}
func Options[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O], middlewares ...func(http.Handler) http.Handler) *EndpointModifier[I, O] {
	return Method(bc, "OPTIONS", path, action, middlewares...)
}

func (m *EndpointModifier[I, O]) After(f func(op *openapi3.Operation)) *EndpointModifier[I, O] {
	return &EndpointModifier[I, O]{
		Handler:  m.Handler,
		register: m.register.After(f),
	}
}

func (m *EndpointModifier[I, O]) OperationID(operationID string) *EndpointModifier[I, O] {
	return m.After(func(op *openapi3.Operation) {
		op.OperationID = strings.TrimSpace(operationID)
	})
}

func (m *EndpointModifier[I, O]) Description(description string) *EndpointModifier[I, O] {
	return m.After(func(op *openapi3.Operation) {
		op.Description = strings.TrimSpace(description)
	})
}
func (m *EndpointModifier[I, O]) Status(code int) *EndpointModifier[I, O] {
	return m.After(func(op *openapi3.Operation) {
		def, ok := op.Responses["200"]
		if ok {
			delete(op.Responses, "200")
			op.Responses[strconv.Itoa(code)] = def
		}
	})
}
func (m *EndpointModifier[I, O]) Tags(tags ...string) *EndpointModifier[I, O] {
	return m.After(func(op *openapi3.Operation) {
		added := make([]string, 0, len(tags))
		for _, x := range tags {
			found := false
			for _, y := range op.Tags {
				if x == y {
					found = true
					break
				}
			}
			if !found {
				added = append(added, x)
			}
		}
		op.Tags = append(op.Tags, added...)
	})
}
func (m *EndpointModifier[I, O]) AnotherError(bc *BuildContext, code int, typ interface{}, description string) *EndpointModifier[I, O] {
	return m.After(func(op *openapi3.Operation) {
		bc.m.RegisterType(typ, func(ref *openapi3.SchemaRef) {
			val := openapi3.NewResponse().WithDescription(description).WithJSONSchemaRef(ref)
			op.Responses[strconv.Itoa(code)] = &openapi3.ResponseRef{Value: val}
		})
	})
}
func (m *EndpointModifier[I, O]) Example(code int, description string, value interface{}) *EndpointModifier[I, O] {
	return &EndpointModifier[I, O]{
		Handler:  m.Handler,
		register: m.register.Example(code, "application/json", "", description, value),
	}
}
func (m *EndpointModifier[I, O]) DefaultInput(fn func() I) *EndpointModifier[I, O] {
	// side effect!
	m.Handler.Default = fn
	m.register = m.register.DefaultInput(fn())
	return m
}

func GetHTML[I any](bc *BuildContext, path string, action quickapi.Action[I, string], dump quickapi.DumpFunc[string], middlewares ...func(http.Handler) http.Handler) *EndpointModifier[I, string] {
	h := quickapi.NewHandlerWithCustomDump(action, dump)
	m := bc.m
	method := "GET"

	if bc.c.Loaded {
		op := findOpenapi3Operation(m.Doc, method, path)
		if op == nil {
			msg := fmt.Sprintf("path not found: %s %s", method, path)
			if !shared.FORCE {
				panic(msg)
			} else {
				log.Printf("[INFO]  %s", msg)
			}
		} else {
			middleware := bc.mb.BuildMiddleware(path, op)
			middlewares = append([]func(http.Handler) http.Handler{middleware}, middlewares...)
		}
		bc.r.With(middlewares...).Method(method, path, h)
	}

	// if c.Loaded is true, this thunk is ignored.
	return &EndpointModifier[I, string]{
		Handler: h,
		register: m.RegisterFuncText(action, "text/html").After(func(op *openapi3.Operation) {
			m.Doc.AddOperation(path, method, op)
			middleware := bc.mb.BuildMiddleware(path, op)
			middlewares := append([]func(http.Handler) http.Handler{middleware}, middlewares...)
			bc.r.With(middlewares...).Method(method, path, h)
		}),
	}
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
