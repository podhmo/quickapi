package quickapi

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"github.com/go-chi/render"
)

type Action[I any, O any] func(ctx context.Context, input I) (output O, err error)

// Empty is zero Input
type Empty struct{}

// Lift transforms Action to http.Handler
func Lift[I any, O any](action Action[I, O]) http.HandlerFunc {
	var iz I
	isEmpty := reflect.TypeOf(iz).NumField() == 0

	return func(w http.ResponseWriter, req *http.Request) {
		var input I
		if !isEmpty {
			if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
				log.Printf("[ERROR] unexpected error (json.Decode): %+v", err) // TODO: structured logging
				writeJSONError(w, req, err, 500)
				return
			}
		}

		output, err := action(req.Context(), input)
		if err != nil {
			code := StatusCodeOf(err)
			if code == 500 {
				log.Printf("[ERROR] unexpected error: %+v", err) // TODO: structured logging
			}
			writeJSONError(w, req, err, code)
		}

		// TODO: support recursive structure (for openAPI)
		// Force to return empty JSON array [] instead of null in case of zero slice.
		if val := reflect.ValueOf(output); val.Kind() == reflect.Slice && val.IsNil() {
			output = reflect.MakeSlice(val.Type(), 0, 0).Interface().(O)
		}
		render.JSON(w, req, output)
	}
}

type errorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func writeJSONError(w http.ResponseWriter, req *http.Request, err error, code int) {
	v := errorResponse{Error: "internal server error", Code: code}
	if code != 500 {
		v.Error = err.Error()
	}

	render.Status(req, code)
	render.JSON(w, req, v)
}
