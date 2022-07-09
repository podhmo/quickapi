package quickapi

import (
	"context"
	"net/http"

	"github.com/podhmo/quickapi/qbind"
	"github.com/podhmo/quickapi/qdump"
)

func NewAPIError(err error, code int) interface {
	error
	qdump.HasStatusCode
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
		input, err := qbind.Bind[I](req, metadata)
		if err != nil {
			qdump.DumpError(w, req, err, 400)
		}

		output, err := action(req.Context(), input)
		qdump.Dump(w, req, output, err)
	}
}
