package definerpc

import (
	"net/http"

	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/qopenapi/define"
)

func Action[I, O any](bc *define.BuildContext, action quickapi.Action[I, O], middlewares ...func(http.Handler) http.Handler) *define.EndpointModifier[I, O] {
	method := "POST"
	shape := bc.ReflectOpenAPIManager().Visitor.Extractor.Extract(action)
	return define.Method(bc, method, "/"+shape.FullName(), action, middlewares...)
}
