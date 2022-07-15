package define

import (
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/podhmo/quickapi"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

type EndpointModifier reflectopenapi.RegisterFuncAction

func (m *EndpointModifier) After(f func(op *openapi3.Operation)) *EndpointModifier {
	return (*EndpointModifier)((*reflectopenapi.RegisterFuncAction)(m).After(f))
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

func Method[I any, O any](bc *BuildContext, method, path string, action quickapi.Action[I, O]) *EndpointModifier {
	h := quickapi.Lift(action)
	bc.r.Method(method, path, h)
	m := bc.m
	return (*EndpointModifier)(m.RegisterFunc(action).After(func(op *openapi3.Operation) {
		m.Doc.AddOperation(path, method, op)
	}))
}

func Get[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O]) *EndpointModifier {
	return Method(bc, "GET", path, action)
}
func Post[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O]) *EndpointModifier {
	return Method(bc, "POST", path, action)
}
func Put[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O]) *EndpointModifier {
	return Method(bc, "PUT", path, action)
}
func Patch[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O]) *EndpointModifier {
	return Method(bc, "PATCH", path, action)
}
func Delete[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O]) *EndpointModifier {
	return Method(bc, "DELETE", path, action)
}
func Head[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O]) *EndpointModifier {
	return Method(bc, "HEAD", path, action)
}
func Options[I any, O any](bc *BuildContext, path string, action quickapi.Action[I, O]) *EndpointModifier {
	return Method(bc, "OPTIONS", path, action)
}
