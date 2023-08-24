package validate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/go-chi/chi/v5"
	"github.com/podhmo/quickapi/shared"
)

type Config struct {
	Debug bool

	EnablePathVarValidation bool // returns 404 if invalid pathvar is passed

	EnableRequestValidation  bool // returns 400 if invalid request
	EnableResponseValidation bool // returns error if invalid response

	RequestErrorStatusCode  int
	ResponseErrorStatusCode int
	NewErrorResponseFunc    func(int, error) any
}

func NewBuilder(doc *openapi3.T, debug bool) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		Doc: doc,
		Config: &Config{
			Debug:                    debug,
			EnablePathVarValidation:  true,
			EnableRequestValidation:  true,
			EnableResponseValidation: false,
			RequestErrorStatusCode:   http.StatusBadRequest,
			ResponseErrorStatusCode:  http.StatusServiceUnavailable,
			NewErrorResponseFunc: func(code int, err error) any {
				errRes := shared.ErrorResponse{
					Code:   code,
					Detail: strings.Split(fmt.Sprintf("%+v", err), "\n"),
				}
				if len(errRes.Detail) > 0 {
					errRes.Error = errRes.Detail[0]
				}
				return errRes
			},
		},
	}
}

func NewBuilderForTest(doc *openapi3.T, debug bool) *MiddlewareBuilder {
	b := NewBuilder(doc, debug)
	b.Config.EnableResponseValidation = true
	return b
}

type MiddlewareBuilder struct {
	Doc    *openapi3.T
	Config *Config
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
			Config:    b.Config,
		}
	}
}

type Middleware struct {
	BaseRoute *routers.Route
	Next      http.Handler
	Config    *Config
}

func (v *Middleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	config := v.Config
	debug := config.Debug

	ctx := req.Context()
	logger := shared.GetLoggerOrNil(ctx)

	route := *v.BaseRoute // shallow copy
	route.Method = req.Method
	route.Path = req.URL.Path

	chiURLParams := chi.RouteContext(req.Context()).URLParams
	pathParams := make(map[string]string, len(chiURLParams.Keys))
	for i, k := range chiURLParams.Keys {
		pathParams[k] = chiURLParams.Values[i]
	}

	input := &openapi3filter.RequestValidationInput{
		Route:       &route,
		Request:     req,
		PathParams:  pathParams,
		QueryParams: req.URL.Query(),
	}

	// path vars validation
	if config.EnablePathVarValidation {
		if err := RequestOnlyPathValiables(ctx, input); err != nil {
			code := http.StatusNotFound
			w.WriteHeader(code)

			enc := json.NewEncoder(w)
			enc.SetIndent("", "\t")
			if err := enc.Encode(config.NewErrorResponseFunc(code, err)); err != nil {
				if debug && logger != nil {
					logger.Printf("unexpected json encode error: %+v", err)
				}
			}

			if debug && logger != nil {
				logger.Printf("path vars validation: %+v", err)
			}
			return
		}
	}

	// request validation
	if config.EnableRequestValidation {
		if err := openapi3filter.ValidateRequest(ctx, input); err != nil {
			code := v.Config.RequestErrorStatusCode
			w.WriteHeader(code)

			enc := json.NewEncoder(w)
			enc.SetIndent("", "\t")
			if err := enc.Encode(config.NewErrorResponseFunc(code, err)); err != nil {
				if debug && logger != nil {
					logger.Printf("unexpected json encode error: %+v", err)
				}
			}

			if debug && logger != nil {
				logger.Printf("request validation: %+v", err)
			}

			return
		}
	}

	if !config.EnableResponseValidation {
		if debug && logger != nil {
			logger.Printf("skip response validation") // todo: route
		}
		v.Next.ServeHTTP(w, req)
		return
	}

	// response validation
	bw := &nethttpBodyWriterProxy{ResponseWriter: w, body: new(bytes.Buffer)}
	v.Next.ServeHTTP(bw, req)

	// if v := w.Header().Get("Content-Type"); v == "" {
	// 	w.Header().Set("Content-Type", "application/json")
	// }

	code := bw.status
	if 200 <= code && code < 300 && code != http.StatusNoContent { // 2xx
		body := bw.body.Bytes()
		responseValidationInput := &openapi3filter.ResponseValidationInput{
			RequestValidationInput: input,
			Status:                 code,
			Header:                 w.Header(),
			Body:                   io.NopCloser(bytes.NewBuffer(body)),
		}
		if err := openapi3filter.ValidateResponse(ctx, responseValidationInput); err != nil {
			code := config.ResponseErrorStatusCode
			w.WriteHeader(code)

			enc := json.NewEncoder(w)
			enc.SetIndent("", "\t")
			if err := enc.Encode(config.NewErrorResponseFunc(code, err)); err != nil {
				if debug && logger != nil {
					logger.Printf("unexpected json encode error: %+v", err)
				}
			}

			if debug && logger != nil {
				logger.Printf("response validation: %+v", err)
				logger.Printf("\tresponse body: %s", body)
			}
			return
		}
	}

	// strict response
	if code != 0 {
		w.WriteHeader(code)
	}
	w.Write(bw.body.Bytes())
}

func RequestOnlyPathValiables(ctx context.Context, input *openapi3filter.RequestValidationInput) error {
	route := input.Route
	pathParameters := make([]*openapi3.Parameter, 0, 8)
	for _, p := range route.PathItem.Parameters {
		if p.Value.In != openapi3.ParameterInPath {
			continue
		}
		if route.Operation.Parameters != nil {
			if override := route.Operation.Parameters.GetByInAndName(p.Value.In, p.Value.Name); override != nil {
				continue
			}
		}
		pathParameters = append(pathParameters, p.Value)

	}
	for _, p := range route.Operation.Parameters {
		if p.Value.In != openapi3.ParameterInPath {
			continue
		}
		pathParameters = append(pathParameters, p.Value)
	}
	for _, p := range pathParameters {
		if err := openapi3filter.ValidateParameter(ctx, input, p); err != nil {
			return err
		}
	}
	return nil
}

type nethttpBodyWriterProxy struct {
	http.ResponseWriter

	status int
	body   *bytes.Buffer
}

func (p nethttpBodyWriterProxy) Write(b []byte) (int, error) {
	return p.body.Write(b)
}

func (p *nethttpBodyWriterProxy) WriteHeader(code int) {
	p.status = code
}
