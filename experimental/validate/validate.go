package validate

import (
	"context"
	"log"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/podhmo/quickapi/qdump"
)

func Middleware(doc *openapi3.T, op *openapi3.Operation, pattern string) func(http.Handler) http.Handler {
	pathItem := doc.Paths.Find(pattern)
	route := &routers.Route{
		Spec:      doc,
		Server:    doc.Servers[0], // xxx
		PathItem:  pathItem,
		Operation: op,
	}
	return func(next http.Handler) http.Handler {
		return &Validator{
			BaseRoute: route,
			Next:      next,
			Extractor: &Extractor{Debug: true},
		}
	}
}

type Validator struct {
	BaseRoute *routers.Route
	Next      http.Handler
	Extractor *Extractor
}

func (v *Validator) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	reqResult := v.Extractor.ExtractRequestValidation(ctx, req, *v.BaseRoute)
	if reqResult.Error != nil {
		code := http.StatusBadRequest
		render.Status(req, code)
		render.JSON(w, req, qdump.NewAPIError(reqResult.Error, code))
		return
	}
	v.Next.ServeHTTP(w, req) // after qdump.Dump()
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
			log.Printf("[DEBUG]  validate request is failed: %T\n%+v", err, err)
		}
		return RequestValidation{Route: &route, Input: input, Error: err}
	}
	if e.Debug {
		log.Printf("[DEBUG] request is OK") // todo: path and parameters
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