package qbind

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"sync"

	"github.com/gorilla/schema"
)

var (
	// for parameters binding
	mu      sync.Mutex
	decoder = schema.NewDecoder()
)

func Bind[I any](req *http.Request, metadata Metadata) (I, error) {
	var input I
	if !metadata.IsEmpty {
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			log.Printf("[ERROR] unexpected error (json.Decode): %+v", err) // TODO: structured logging
			return input, err
		}
	}
	return input, nil
}

type Metadata struct {
	IsEmpty bool // Action is empty
}

func Scan[I any, O any](action func(context.Context, I) (O, error)) Metadata {
	var iz I
	isEmpty := reflect.TypeOf(iz).NumField() == 0
	return Metadata{
		IsEmpty: isEmpty,
	}
}
