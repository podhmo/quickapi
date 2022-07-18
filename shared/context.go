package shared

import (
	"context"
	"net/http"
)

type contextKey string

const (
	requestKey contextKey = "request"
)

func GetRequest(ctx context.Context) *http.Request {
	return ctx.Value(requestKey).(*http.Request)
}

func SetRequest(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, requestKey, req)
}
