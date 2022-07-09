package qbind

import (
	"context"
	"encoding/json"
	"fmt"
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
		fmt.Println("ok")
		DEBUG = ok
	}

}

func Bind[I any](req *http.Request, metadata Metadata) (I, error) {
	var input I
	if metadata.HasData {
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			log.Printf("[ERROR] unexpected error (json.Decode): %+v, on %T", err, input) // TODO: structured logging
			return input, err
		}
		// TODO: validation
	}

	// TODO: handling metadata (query tag)
	if err := queryDecoder.Decode(&input, req.URL.Query()); err != nil {
		if DEBUG {
			log.Printf("[DEBUG] unexpected query string: %+v, on %T", err, input)
		}
	}

	// TODO: handling metadata (header tag)
	if err := queryDecoder.Decode(&input, req.Header); err != nil {
		if DEBUG {
			log.Printf("[DEBUG] unexpected header: %+v, on %T", err, input)
		}
	}
	return input, nil
}

type Metadata struct {
	HasData bool     // Action is empty
	Queries []string // query string keys (recursive structure is not supported, also embedded)
	Headers []string // header keys (recursive structure is not supported, also embedded)
}

func Scan[I any, O any](action func(context.Context, I) (O, error)) Metadata {
	var iz I
	rt := reflect.TypeOf(iz)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	var queries []string
	var headers []string
	for i, n := 0, rt.NumField(); i < n; i++ {
		if v, ok := rt.Field(i).Tag.Lookup("query"); ok {
			queries = append(queries, v)
			continue
		}
		if v, ok := rt.Field(i).Tag.Lookup("header"); ok {
			headers = append(headers, v)
			continue
		}
	}
	return Metadata{
		HasData: rt.NumField()-len(queries)-len(headers) == 0,
		Queries: queries,
		Headers: headers,
	}
}
