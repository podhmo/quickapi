package shared

import (
	"context"
	"net/http"

	"golang.org/x/exp/slog"
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

const (
	loggerKey contextKey = "logger"
)

func GetLogger(ctx context.Context) *slog.Logger {
	l := ctx.Value(loggerKey)
	if l == nil {
		return slog.Default()
	}
	return l.(*slog.Logger)
}
func GetLoggerOrNil(ctx context.Context) *slog.Logger {
	l := ctx.Value(loggerKey)
	if l == nil {
		return nil
	}
	return l.(*slog.Logger)
}
func SetLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}
