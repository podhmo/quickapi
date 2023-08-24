package definerpc

import (
	"fmt"
	"net/http"

	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/qopenapi/define"
)

func Action[I, O any](bc *define.BuildContext, action quickapi.Action[I, O], middlewares ...func(http.Handler) http.Handler) *define.EndpointModifier[I, O] {
	method := "POST"
	shape := bc.ReflectOpenAPIManager().Visitor.Extractor.Extract(action)
	pkgname := shape.Package.Name
	name := shape.Name
	return define.Method(bc, method, fmt.Sprintf("/%s.%s", pkgname, name), action, middlewares...)
}
