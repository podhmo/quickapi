package qdump

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
)

func Dump[O any](w http.ResponseWriter, req *http.Request, output O, err error) {
	if err != nil {
		select {
		case <-req.Context().Done():
			// [http headers - Possible reason for NGINX 499 error codes - Stack Overflow](https://stackoverflow.com/questions/12973304/possible-reason-for-nginx-499-error-codes)
			DumpError(w, req, err, 499)
			return
		default:
		}

		code := StatusCodeOf(err)
		if code == 500 {
			log.Printf("[ERROR] unexpected error: %+v", err) // TODO: structured logging
		}
		DumpError(w, req, err, code)
		return
	}

	if t, ok := any(output).(Redirector); ok {
		t.Redirect(w, req)
		return
	}

	// Force to return empty JSON array [] instead of null in case of zero slice.
	output = FillNil(output)

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
