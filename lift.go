package quickapi

import (
	"context"
	"net/http"

	"github.com/podhmo/quickapi/qbind"
	"github.com/podhmo/quickapi/qdump"
	"github.com/podhmo/quickapi/shared"
)

type Action[I any, O any] func(ctx context.Context, input I) (output O, err error)

// Lift transforms Action to http.Handler
func Lift[I any, O any](action Action[I, O]) http.HandlerFunc {
	metadata := qbind.Scan(action)
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := shared.SetRequest(req.Context(), req) // for shared.GetRequest() in action
		req = req.WithContext(ctx)

		// binding request body and query-string and headers to input.
		input, err := qbind.Bind[I](ctx, req, metadata)
		if err != nil {
			code := shared.StatusCodeOfOrDefault(err, 400)
			qdump.DumpError(w, req, err, code)
			return
		}

		output, err := action(req.Context(), input)

		// dumping result or error as json response
		qdump.Dump(ctx, w, req, output, err)
	}
}
