package quickapi

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"github.com/podhmo/quickapi/qdump"
)

type Action[I any, O any] func(ctx context.Context, input I) (output O, err error)

// Empty is zero Input
type Empty struct{}

func Bind[I any](req *http.Request, isEmpty bool) (I, error) {
	var input I
	if !isEmpty {
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			log.Printf("[ERROR] unexpected error (json.Decode): %+v", err) // TODO: structured logging
			return input, err
		}
	}
	return input, nil

}

// Lift transforms Action to http.Handler
func Lift[I any, O any](action Action[I, O]) http.HandlerFunc {
	var iz I
	isEmpty := reflect.TypeOf(iz).NumField() == 0

	return func(w http.ResponseWriter, req *http.Request) {
		input, err := Bind[I](req, isEmpty)
		if err != nil {
			qdump.DumpError(w, req, err, 400)
		}

		output, err := action(req.Context(), input)
		qdump.Dump(w, req, output, err)
	}
}
