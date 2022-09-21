package quickapi

import (
	"context"
	"net/http"

	"github.com/podhmo/quickapi/qbind"
	"github.com/podhmo/quickapi/qdump"
	"github.com/podhmo/quickapi/shared"
)

type Action[I any, O any] func(ctx context.Context, input I) (output O, err error)
type DumpFunc[O any] func(ctx context.Context, w http.ResponseWriter, req *http.Request, output O, err error)
type LiftedHandler[I any, O any] struct {
	Action   Action[I, O]
	Metadata qbind.Metadata
	Dump     DumpFunc[O]
}

func (h *LiftedHandler[I, O]) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := shared.SetRequest(req.Context(), req) // for shared.GetRequest() in action
	req = req.WithContext(ctx)

	// binding request body and query-string and headers to input.
	input, err := qbind.Bind[I](ctx, req, h.Metadata)
	if err != nil {
		code := shared.StatusCodeOfOrDefault(err, 400)
		qdump.DumpError(w, req, err, code)
		return
	}

	output, err := h.Action(req.Context(), input)

	// dumping result or error as json response
	h.Dump(ctx, w, req, output, err)
}

// Lift transforms Action to http.Handler
func Lift[I any, O any](action Action[I, O]) *LiftedHandler[I, O] {
	metadata := qbind.Scan(action)
	return &LiftedHandler[I, O]{
		Action:   action,
		Metadata: metadata,
		Dump:     qdump.Dump[O],
	}
}
