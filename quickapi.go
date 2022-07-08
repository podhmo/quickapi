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

type Empty struct{}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"Code"`
}

var (
	Pretty bool
)

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

		code := 500
		output, err := action(req.Context(), input)
		if err != nil {
			// TODO: handling status code
			if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
				if code == 500 {
					log.Printf("[ERROR] unexpected error: %+v", err) // TODO: structured logging
				}
				writeJSONError(w, req, err, code)
			}
		}

		// TODO: support recursive structure (for openAPI)
		// Force to return empty JSON array [] instead of null in case of zero slice.
		if val := reflect.ValueOf(output); val.Kind() == reflect.Slice && val.IsNil() {
			output = reflect.MakeSlice(val.Type(), 0, 0).Interface().(O)
		}

		ctx := context.WithValue(req.Context(), render.ContentTypeCtxKey, code)
		render.JSON(w, req.WithContext(ctx), output)
	}
}

func writeJSONError(w http.ResponseWriter, req *http.Request, err error, code int) {
	v := ErrorResponse{Error: "internal server error", Code: code}
	if code != 500 {
		v.Error = err.Error()
	}

	ctx := context.WithValue(req.Context(), render.ContentTypeCtxKey, code)
	render.JSON(w, req.WithContext(ctx), v)
}
