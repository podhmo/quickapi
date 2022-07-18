package qbind

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"

	"github.com/gorilla/schema"
)

var (
	queryDecoder  = schema.NewDecoder()
	headerDecoder = schema.NewDecoder()
	DEBUG         = false
)

func init() {
	queryDecoder.SetAliasTag("query")
	headerDecoder.SetAliasTag("header")
	if ok, _ := strconv.ParseBool(os.Getenv("DEBUG")); ok {
		DEBUG = ok
	}

}

func Bind[I any](req *http.Request, metadata Metadata) (I, error) {
	var input I
	if metadata.HasData {
		if req.Body == nil {
			log.Printf("[INFO] decode json is neaded, but request body is nil, metadata=%+v, on %T", metadata, input) // TODO: structured logging
		} else if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			log.Printf("[ERROR] unexpected error (json.Decode): %+v, on %T", err, input) // TODO: structured logging
			return input, err
		}
		// TODO: validation
	}

	if len(metadata.Queries) > 0 {
		m := make(map[string][]string, len(metadata.Queries))
		v := req.URL.Query()
		for _, k := range metadata.Queries {
			m[k] = []string{v.Get(k)}
		}
		if err := queryDecoder.Decode(&input, m); err != nil {
			if DEBUG {
				log.Printf("[DEBUG] unexpected query string: %+v, on %T", err, input)
			}
		}
	}

	if len(metadata.Headers) > 0 {
		m := make(map[string][]string, len(metadata.Headers))
		for _, k := range metadata.Headers {
			m[k] = []string{req.Header.Get(k)}
		}
		if err := headerDecoder.Decode(&input, m); err != nil {
			if DEBUG {
				log.Printf("[DEBUG] unexpected header: %+v, on %T", err, input)
			}
		}

	}
	return input, nil
}

type Metadata struct {
	HasData bool // Action is empty

	JSONFields []string
	Queries    []string // query string keys (recursive structure is not supported, also embedded)
	Headers    []string // header keys (recursive structure is not supported, also embedded)
	PathVars   []string // path variables
}

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
	if DEBUG {
		log.Printf("[DEBUG] on %T, metadata=%+v", iz, metadata)
	}
	return metadata
}
