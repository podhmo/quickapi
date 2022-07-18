package validate

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/podhmo/quickapi/shared"
)

func NewBuilder(doc *openapi3.T, debug bool) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		Doc:       doc,
		Extractor: &Extractor{Debug: debug},
	}
}

type MiddlewareBuilder struct {
	Doc       *openapi3.T
	Extractor *Extractor
}

func (b *MiddlewareBuilder) BuildMiddleware(pattern string, op *openapi3.Operation) func(http.Handler) http.Handler {
	doc := b.Doc
	pathItem := doc.Paths.Find(pattern)
	route := &routers.Route{
		Spec:      doc,
		Server:    doc.Servers[0], // xxx
		PathItem:  pathItem,
		Operation: op,
	}
	return func(next http.Handler) http.Handler {
		return &Middleware{
			BaseRoute: route,
			Next:      next,
			Extractor: b.Extractor,
		}
	}
}

type Middleware struct {
	BaseRoute *routers.Route
	Next      http.Handler
	Extractor *Extractor
}

func (v *Middleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	reqResult := v.Extractor.ExtractRequestValidation(ctx, req, *v.BaseRoute)
	if err := reqResult.Error; err != nil {
		code := http.StatusBadRequest
		detail := strings.Split(fmt.Sprintf("%v", err), "\n")
		value := &shared.APIError{Code: code, Error: detail[0], Detail: detail}

		render.Status(req, code)
		render.JSON(w, req, value)
		return
	}
	v.Next.ServeHTTP(w, req) // after qdump.Dump()
	// TODO: response validation
}

type Extractor struct {
	Debug bool
}

type RequestValidation struct {
	Route *routers.Route
	Input *openapi3filter.RequestValidationInput
	Error error
}

func (e *Extractor) ExtractRequestValidation(ctx context.Context, req *http.Request, base routers.Route) RequestValidation {
	route := base // shallow copy
	route.Method = req.Method
	route.Path = req.URL.Path

	chiURLParams := chi.RouteContext(req.Context()).URLParams
	pathParams := make(map[string]string, len(chiURLParams.Keys))
	for i, k := range chiURLParams.Keys {
		pathParams[k] = chiURLParams.Values[i]
	}

	input := &openapi3filter.RequestValidationInput{
		Request:     req,
		PathParams:  pathParams,
		QueryParams: req.URL.Query(),
		Route:       &route,
	}
	if err := openapi3filter.ValidateRequest(ctx, input); err != nil {
		if e.Debug {
			log.Printf("[DEBUG] request is NG (%T) method=%s, path=%s, operationId=%s\n%+v", err, route.Method, route.Path, route.Operation.OperationID, err)
		}
		return RequestValidation{Route: &route, Input: input, Error: err}
	}
	if e.Debug {
		log.Printf("[DEBUG] request is OK method=%s, path=%s, operationId=%s", route.Method, route.Path, route.Operation.OperationID)
	}
	return RequestValidation{Route: &route, Input: input}
}

type ResponseValidation struct {
	Route *routers.Route
	Input *openapi3filter.ResponseValidationInput
	Error error
}

func (e *Extractor) ExtractResponseValidation(ctx context.Context, validation *RequestValidation, res *http.Response) ResponseValidation {
	input := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: validation.Input,
		Status:                 res.StatusCode,
		Header:                 res.Header,
		Body:                   res.Body,
	}
	if err := openapi3filter.ValidateResponse(ctx, input); err != nil {
		if e.Debug {
			log.Printf("[DEBUG] validate response is failed: %T\n%+v", err, err)
		}
		return ResponseValidation{Route: validation.Route, Input: input, Error: err}
	}
	if e.Debug {
		log.Printf("[DEBUG] response is OK") // todo: path and parameters
	}
	return ResponseValidation{Route: validation.Route, Input: input}
}
