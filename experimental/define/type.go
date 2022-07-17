package define

import (
	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

type TypeModifier reflectopenapi.RegisterTypeAction

func Type(bc *BuildContext, typ interface{}) *TypeModifier {
	return (*TypeModifier)(bc.m.RegisterType(typ))
}

func Enum[T any](bc *BuildContext, defaultValue T, values ...T) *TypeModifier {
	dst := make([]interface{}, len(values)+1)
	typedValue := T(defaultValue)
	dst[0] = typedValue
	for i, v := range values {
		dst[i+1] = T(v)
	}
	return (*TypeModifier)(bc.m.RegisterType(typedValue, func(ref *openapi3.Schema) {
		ref.Default = dst[0]
		ref.Enum = dst
	}))
}
func StringEnum[T ~string](bc *BuildContext, defaultValue T, values ...T) *TypeModifier {
	return Enum(bc, defaultValue, values...)
}
func IntEnum[T ~int](bc *BuildContext, defaultValue T, values ...T) *TypeModifier {
	return Enum(bc, defaultValue, values...)
}

func (m *TypeModifier) After(f func(ref *openapi3.SchemaRef)) *TypeModifier {
	return (*TypeModifier)((*reflectopenapi.RegisterTypeAction)(m).After(f))
}
func (m *TypeModifier) Before(f func(s *openapi3.Schema)) *TypeModifier {
	return (*TypeModifier)((*reflectopenapi.RegisterTypeAction)(m).Before(f))
}
