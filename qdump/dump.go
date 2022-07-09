package qdump

import (
	"log"
	"net/http"
	"reflect"

	"github.com/go-chi/render"
)

func Dump[O any](w http.ResponseWriter, req *http.Request, output O, err error) {
	if err != nil {
		code := StatusCodeOf(err)
		if code == 500 {
			log.Printf("[ERROR] unexpected error: %+v", err) // TODO: structured logging
		}
		DumpError(w, req, err, code)
		return
	}

	// TODO: support recursive structure (for openAPI)
	// Force to return empty JSON array [] instead of null in case of zero slice.
	if val := reflect.ValueOf(output); val.Kind() == reflect.Slice && val.IsNil() {
		output = reflect.MakeSlice(val.Type(), 0, 0).Interface().(O)
	}
	render.JSON(w, req, output)
}

type errorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func DumpError(w http.ResponseWriter, req *http.Request, err error, code int) {
	v := errorResponse{Error: "internal server error", Code: code}
	if code != 500 {
		v.Error = err.Error()
	}

	render.Status(req, code)
	render.JSON(w, req, v)
}
