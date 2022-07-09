package qbind

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/schema"
)

var (
	queryDecoder  = schema.NewDecoder()
	headerDecoder = schema.NewDecoder()
)

func init() {
	queryDecoder.SetAliasTag("query")
	headerDecoder.SetAliasTag("header")
}

func Bind[I any](req *http.Request, metadata Metadata) (I, error) {
	var input I
	if metadata.HasData {
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			log.Printf("[ERROR] unexpected error (json.Decode): %+v", err) // TODO: structured logging
			return input, err
		}
		// TODO: validation
	}

	// TODO: handling metadata (query tag)
	if err := queryDecoder.Decode(&input, req.URL.Query()); err != nil {
		log.Printf("[DEBUG] unexpected query string: %+v", err)
	}

	// TODO: handling metadata (header tag)
	if err := queryDecoder.Decode(&input, req.Header); err != nil {
		log.Printf("[DEBUG] unexpected query string: %+v", err)
	}
	return input, nil
}

type Metadata struct {
	HasData bool // Action is empty
}

func Scan[I any, O any](action func(context.Context, I) (O, error)) Metadata {
	var iz I
	isEmpty := reflect.TypeOf(iz).NumField() == 0
	return Metadata{
		HasData: !isEmpty,
	}
}
