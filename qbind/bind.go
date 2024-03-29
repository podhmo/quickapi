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
	"golang.org/x/exp/slog"
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

type Validator interface {
	Validate(context.Context) error
}

func Bind[I any](ctx context.Context, req *http.Request, metadata Metadata, input *I) error {
	if len(metadata.PathVars) > 0 {
		pathparams := chi.RouteContext(req.Context()).URLParams
		if len(metadata.PathVars) != len(pathparams.Keys) {
			return shared.NewAPIError(fmt.Errorf("route is not found"), http.StatusNotFound)
		}

		m := make(map[string][]string, len(metadata.PathVars))
		for i, k := range pathparams.Keys {
			v := pathparams.Values[i]
			m[k] = []string{v}
		}
		if err := pathDecoder.Decode(input, m); err != nil {
			if shared.INFO {
				shared.GetLogger(ctx).Debug("route path is broken", slog.Any("error", err), slog.Any("params", m))
			}
			return shared.NewAPIError(ErrNotFound, http.StatusNotFound)
		}
	}

	if metadata.HasData {
		switch req.Method {
		case http.MethodGet, http.MethodHead, http.MethodConnect, http.MethodOptions, http.MethodTrace:
			if shared.INFO {
				shared.GetLogger(ctx).Info("cannot receive request body", slog.String("method", req.Method), slog.Any("metadata", metadata), slog.Any("input", reflect.TypeOf(input)))
			}
			return shared.NewAPIError(ErrCannotReceiveBody, http.StatusBadRequest)
		default:
			if req.Body == nil || req.Body == http.NoBody {
				if shared.INFO {
					shared.GetLogger(ctx).Info("decode json is needed, but request body is nil", slog.String("method", req.Method), slog.Any("metadata", metadata), slog.Any("input", reflect.TypeOf(input)))
				}
				return shared.NewAPIError(ErrNoBody, http.StatusBadRequest)
			} else if err := json.NewDecoder(req.Body).Decode(input); err != nil {
				if shared.INFO {
					shared.GetLogger(ctx).Error("unexpected error (json.Decode)", slog.Any("error", err), slog.Any("input", reflect.TypeOf(input)))
				}
				return err
			}
		}
	}

	if len(metadata.Queries) > 0 {
		m := make(map[string][]string, len(metadata.Queries))
		v := req.URL.Query()
		for _, k := range metadata.Queries {
			m[k] = []string{v.Get(k)}
		}
		if err := queryDecoder.Decode(input, m); err != nil {
			if shared.INFO {
				shared.GetLogger(ctx).Error("unexpected error (query string)", slog.Any("error", err), slog.Any("input", reflect.TypeOf(input)))
			}
		}
	}

	if len(metadata.Headers) > 0 {
		m := make(map[string][]string, len(metadata.Headers))
		for _, k := range metadata.Headers {
			m[k] = []string{req.Header.Get(k)}
		}
		if err := headerDecoder.Decode(input, m); err != nil {
			if shared.INFO {
				shared.GetLogger(ctx).Error("unexpected error (header)", slog.Any("error", err), slog.Any("input", reflect.TypeOf(input)))
			}
		}

	}

	if t, ok := any(input).(Validator); ok {
		if err := t.Validate(req.Context()); err != nil {
			if shared.INFO {
				shared.GetLogger(ctx).Error("validation is failed", slog.Any("error", err), slog.Any("input", reflect.TypeOf(input)))
			}
			return err
		}
	}

	return nil
}

// // TODO: omit gorilla/schema
// type Field struct {
// 	TagName   string
// 	FieldName string
// 	Set       func(reflect.Value, Field) error
// }

type Metadata struct {
	HasData bool // Action is empty

	Input reflect.Type

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

	var r result
	scan(&r, rt)

	metadata := Metadata{
		HasData:    len(r.jsonfields) > 0,
		Input:      rt,
		JSONFields: r.jsonfields,
		Queries:    r.queries,
		Headers:    r.headers,
		PathVars:   r.pathvars,
	}
	if shared.DEBUG {
		log.Printf("[DEBUG] Scan %T, metadata=%+v", iz, metadata)
	}
	return metadata
}

type result struct {
	queries    []string
	headers    []string
	jsonfields []string
	pathvars   []string
}

func scan(r *result, rt reflect.Type) {
	for i, n := 0, rt.NumField(); i < n; i++ {
		field := rt.Field(i)
		if field.Anonymous { // embedded support by recursive call
			embedded := field.Type
			for embedded.Kind() == reflect.Ptr {
				embedded = embedded.Elem()
			}
			scan(r, embedded)
			continue
		}
		if v, ok := field.Tag.Lookup("query"); ok {
			r.queries = append(r.queries, v)
			continue
		}
		if v, ok := field.Tag.Lookup("header"); ok {
			r.headers = append(r.headers, v)
			continue
		}
		if v, ok := field.Tag.Lookup("path"); ok {
			r.pathvars = append(r.pathvars, v)
			continue
		}
		name := field.Name
		if v, ok := field.Tag.Lookup("json"); ok {
			name = v
		}
		r.jsonfields = append(r.jsonfields, name)
	}
}
