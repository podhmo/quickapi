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
func SetLogger(ctx context.Context, l logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

const (
	dumpFuncKey contextKey = "dumpFunc"
)

type DumpFunc[O any] func(context.Context, http.ResponseWriter, *http.Request, O)

func GetDumpFunc[O any](ctx context.Context, defaultFunc DumpFunc[O]) DumpFunc[O] {
	fn := ctx.Value(dumpFuncKey)
	if fn == nil {
		return defaultFunc
	}
	return fn.(DumpFunc[O])
}
func SetDumpFunc[O any](ctx context.Context, dump DumpFunc[O]) context.Context {
	return context.WithValue(ctx, dumpFuncKey, dump)
}
