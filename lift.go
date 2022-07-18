package quickapi

import (
	"context"
	"net/http"

	"github.com/podhmo/quickapi/qbind"
	"github.com/podhmo/quickapi/qdump"
)

func NewAPIError(err error, code int) interface {
	error
	qdump.StatusCoder
} {
	return qdump.NewAPIError(err, code)
}

type Action[I any, O any] func(ctx context.Context, input I) (output O, err error)

// Empty is zero Input
type Empty struct{}

// Lift transforms Action to http.Handler
func Lift[I any, O any](action Action[I, O]) http.HandlerFunc {
	metadata := qbind.Scan(action)
	return func(w http.ResponseWriter, req *http.Request) {
		req = req.WithContext(qbind.SetRequest(req.Context(), req)) // for qbind.GetRequest() in action

		// binding request body and query-string and headers to input.
		input, err := qbind.Bind[I](req, metadata)
		if err != nil {
			qdump.DumpError(w, req, err, 400)
			return
		}

		output, err := action(req.Context(), input)

		// dumping result or error as json response
		qdump.Dump(w, req, output, err)
	}
}
