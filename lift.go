package quickapi

import (
	"context"
	"net/http"
	"reflect"
	"sync"

	"github.com/podhmo/quickapi/qbind"
	"github.com/podhmo/quickapi/qdump"
	"github.com/podhmo/quickapi/shared"
)

type Action[I any, O any] func(ctx context.Context, input I) (output O, err error)
type DumpFunc[O any] func(ctx context.Context, w http.ResponseWriter, req *http.Request, output O, err error)

type LiftedHandler[I any, O any] struct {
	Action   Action[I, O]
	metadata qbind.Metadata
	Dump     DumpFunc[O]
	Default  func() I
}

func (h *LiftedHandler[I, O]) Metadata() qbind.Metadata {
	return h.metadata
}

func (h *LiftedHandler[I, O]) ToFunc() http.HandlerFunc {
	// for chi.Router.Get()
	fn := http.HandlerFunc(h.ServeHTTP)
	mu.Lock()
	defer mu.Unlock()
	funcToHandler[reflect.ValueOf(fn).Pointer()] = h.metadata
	return fn
}

func (h *LiftedHandler[I, O]) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := shared.SetRequest(req.Context(), req) // for shared.GetRequest() in action
	req = req.WithContext(ctx)

	// binding request body and query-string and headers to input.
	var input I
	if h.Default != nil {
		input = h.Default()
	}
	if err := qbind.Bind(ctx, req, h.metadata, &input); err != nil {
		code := shared.StatusCodeOfOrDefault(err, 400)
		qdump.DumpError(w, req, err, code)
		return
	}

	output, err := h.Action(req.Context(), input)

	// dumping result or error as json response
	h.Dump(ctx, w, req, output, err)
}

// NewHandler create new lifted handler
func NewHandler[I any, O any](action Action[I, O], dump DumpFunc[O]) *LiftedHandler[I, O] {
	metadata := qbind.Scan(action)
	h := &LiftedHandler[I, O]{
		Action:   action,
		metadata: metadata,
		Dump:     dump,
	}
	return h
}

var _ http.Handler = (*LiftedHandler[any, any])(nil)

// LiftHandler transforms Action to http.Handler
func LiftHandler[I any, O any](action Action[I, O]) *LiftedHandler[I, O] {
	return NewHandler(action, qdump.Dump[O])
}

// Lift transforms Action to http.HandlerFunc (for go.chi's router methods)
func Lift[I any, O any](action Action[I, O]) http.HandlerFunc {
	h := LiftHandler(action)
	return h.ToFunc()
}

var (
	funcToHandler = map[uintptr]qbind.Metadata{}
	mu            sync.RWMutex
)

// metadataFromHandlerFunc return metadata from handler func, this is the adapter for chi.Router. chi.Router.Get() receives http.HandlerFunc instead of http.Handler.
func metadataFromHandlerFunc(fn http.HandlerFunc) (qbind.Metadata, bool) {
	mu.RLock()
	defer mu.RUnlock()
	v, ok := funcToHandler[reflect.ValueOf(fn).Pointer()]
	return v, ok
}
