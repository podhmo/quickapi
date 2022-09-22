package shared

import (
	"context"
	"log"
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

const (
	loggerKey contextKey = "logger"
)

type logger interface {
	Printf(format string, v ...any)
}

func GetLogger(ctx context.Context) logger {
	l := ctx.Value(loggerKey)
	if l == nil {
		return log.Default()
	}
	return l.(logger)
}
func GetLoggerOrNil(ctx context.Context) logger {
	l := ctx.Value(loggerKey)
	if l == nil {
		return nil
	}
	return l.(logger)
}
func SetLogger(ctx context.Context, l logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}
