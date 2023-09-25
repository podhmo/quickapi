package quickapi

import (
	"context"
	"net/http"

	"github.com/podhmo/quickapi/shared"
)

var NewAPIError = shared.NewAPIError
var NoContent = shared.NoContent
var GetRequest = shared.GetRequest
var Redirect = shared.Redirect

// Empty is zero Input
type Empty = shared.Empty

// ErrorResponse represents a normal error response type
type ErrorResponse = shared.ErrorResponse

func Inject[V any](v V, set func(context.Context, V) context.Context) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return set(ctx, v)
	}
}
func InjectMiddleware(injects ...func(context.Context) context.Context) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			for _, inj := range injects {
				ctx = inj(ctx)
			}
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
