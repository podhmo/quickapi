package quickapi

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	http.Server

	GracefulShutdownTimeout time.Duration
}

func ListenAndServeWithGracefulShutdown(
	ctx context.Context,
	addr string, handler http.Handler,
	gracefulShutdownTimeout time.Duration,
	options ...func(*Server),
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	server := &Server{
		Server: http.Server{
			Addr:    addr,
			Handler: handler,

			// hmm
			ReadHeaderTimeout: 2 * time.Second,
			IdleTimeout:       30 * time.Second,
		},
		GracefulShutdownTimeout: gracefulShutdownTimeout,
	}
	for _, opt := range options {
		opt(server)
	}

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), server.GracefulShutdownTimeout)
		defer cancel()
		server.Shutdown(ctx)
	}()
	log.Printf("[INFO]  listening: %s", server.Addr)
	return server.ListenAndServe()
}
