package qdump

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"github.com/podhmo/quickapi/shared"
)

func Dump[O any](ctx context.Context, w http.ResponseWriter, req *http.Request, output O, err error) {
	if err != nil {
		select {
		case <-req.Context().Done():
			// [http headers - Possible reason for NGINX 499 error codes - Stack Overflow](https://stackoverflow.com/questions/12973304/possible-reason-for-nginx-499-error-codes)
			DumpError(w, req, err, 499)
			return
		default:
		}

		if t, ok := err.(shared.Redirector); ok {
			t.Redirect(w, req)
			return
		}

		code := shared.StatusCodeOf(err)
		if code == 500 {
			shared.GetLogger(ctx).Printf("[ERROR] unexpected error: %+v", err) // TODO: structured logging
		}
		DumpError(w, req, err, code)
		return
	}

	if t, ok := any(output).(shared.StatusCoder); ok {
		render.Status(req, t.StatusCode())
	}

	dumpFunc := shared.GetDumpFunc(ctx, dumpJSON[O])
	dumpFunc(ctx, w, req, output)
}

func DumpError(w http.ResponseWriter, req *http.Request, err error, code int) {
	v := shared.ErrorResponse{Error: "internal server error", Code: code}
	if code != 500 {
		v.Error = err.Error()
	}

	render.Status(req, code)
	render.JSON(w, req, v)
}

func dumpJSON[O any](ctx context.Context, w http.ResponseWriter, req *http.Request, output O) {
	// Force to return empty JSON array [] instead of null in case of zero slice.
	output = FillNil(ctx, output)

	render.JSON(w, req, output)
}
