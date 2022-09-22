package qbind

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/podhmo/quickapi/shared"
)

var (
	pathDecoder   = schema.NewDecoder()
	queryDecoder  = schema.NewDecoder()
	headerDecoder = schema.NewDecoder()
)

func init() {
	pathDecoder.SetAliasTag("path")
	queryDecoder.SetAliasTag("query")
	headerDecoder.SetAliasTag("header")
}

var (
	ErrNotFound = fmt.Errorf("not found")

	ErrNoBody            = fmt.Errorf("no body")
	ErrCannotReceiveBody = fmt.Errorf("cannot receive body")
)

func Bind[I any](ctx context.Context, req *http.Request, metadata Metadata) (I, error) {
	var input I
	if len(metadata.PathVars) > 0 {
		pathparams := chi.RouteContext(req.Context()).URLParams
		if len(metadata.PathVars) != len(pathparams.Keys) {
			return input, shared.NewAPIError(fmt.Errorf("route is not found"), http.StatusNotFound)
		}

		m := make(map[string][]string, len(metadata.PathVars))
		for i, k := range pathparams.Keys {
			v := pathparams.Values[i]
			m[k] = []string{v}
		}
		if err := pathDecoder.Decode(&input, m); err != nil {
			if shared.INFO {
				shared.GetLogger(ctx).Printf("[DEBUG] route path is broken: %v, params=%+v", err, m) // TODO: structured logging
			}
			return input, shared.NewAPIError(ErrNotFound, http.StatusNotFound)
		}
	}

	if metadata.HasData {
		switch req.Method {
		case http.MethodGet, http.MethodHead, http.MethodConnect, http.MethodOptions, http.MethodTrace:
			if shared.INFO {
				shared.GetLogger(ctx).Printf("[INFO ] request method %q cannot receive request body, metadata=%+v, on %T", req.Method, metadata, input) // TODO: structured logging
			}
			return input, shared.NewAPIError(ErrCannotReceiveBody, http.StatusBadRequest)
		default:
			if req.Body == nil || req.Body == http.NoBody {
				if shared.INFO {
					shared.GetLogger(ctx).Printf("[INFO ] decode json is neaded, but request body is nil, metadata=%+v, on %T", metadata, input) // TODO: structured logging
				}
				return input, shared.NewAPIError(ErrNoBody, http.StatusBadRequest)
			} else if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
				if shared.INFO {
					shared.GetLogger(ctx).Printf("[ERROR] unexpected error (json.Decode): %+v, on %T", err, input) // TODO: structured logging
				}
				return input, err
			}
		}
	}

	if len(metadata.Queries) > 0 {
		m := make(map[string][]string, len(metadata.Queries))
		v := req.URL.Query()
		for _, k := range metadata.Queries {
			m[k] = []string{v.Get(k)}
		}
		if err := queryDecoder.Decode(&input, m); err != nil {
			if shared.INFO {
				shared.GetLogger(ctx).Printf("[ERROR] unexpected query string: %+v, on %T", err, input)
			}
		}
	}

	if len(metadata.Headers) > 0 {
		m := make(map[string][]string, len(metadata.Headers))
		for _, k := range metadata.Headers {
			m[k] = []string{req.Header.Get(k)}
		}
		if err := headerDecoder.Decode(&input, m); err != nil {
			if shared.INFO {
				shared.GetLogger(ctx).Printf("[ERROR] unexpected header: %+v, on %T", err, input)
			}
		}

	}

	if t, ok := any(input).(Validator); ok {
		if err := t.Validate(req.Context()); err != nil {
			if shared.INFO {
				shared.GetLogger(ctx).Printf("[ERROR] validation is failed: %+v, on %T", err, input)
			}
			return input, err
		}
	}

	return input, nil
}

// // TODO: omit gorilla/schema
// type Field struct {
// 	TagName   string
// 	FieldName string
// 	Set       func(reflect.Value, Field) error
// }

type Metadata struct {
	HasData bool // Action is empty

	JSONFields []string
	Queries    []string // query string keys (recursive structure is not supported, also embedded)
	Headers    []string // header keys (recursive structure is not supported, also embedded)

	PathVars []string // path variables
}

// TODO: cache

func Scan[I any, O any](action func(context.Context, I) (O, error)) Metadata {
	var iz I
	rt := reflect.TypeOf(iz)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	var queries []string
	var headers []string
	var jsonfields []string
	var pathvars []string
	for i, n := 0, rt.NumField(); i < n; i++ {
		field := rt.Field(i)
		if v, ok := field.Tag.Lookup("query"); ok {
			queries = append(queries, v)
			continue
		}
		if v, ok := field.Tag.Lookup("header"); ok {
			headers = append(headers, v)
			continue
		}
		if v, ok := field.Tag.Lookup("path"); ok {
			pathvars = append(pathvars, v)
			continue
		}
		name := field.Name
		if v, ok := field.Tag.Lookup("json"); ok {
			name = v
		}
		jsonfields = append(jsonfields, name)
	}

	metadata := Metadata{
		HasData:    len(jsonfields) > 0,
		JSONFields: jsonfields,
		Queries:    queries,
		Headers:    headers,
		PathVars:   pathvars,
	}
	if shared.DEBUG {
		log.Printf("[DEBUG] Scan %T, metadata=%+v", iz, metadata)
	}
	return metadata
}

type Validator interface {
	Validate(context.Context) error
}
