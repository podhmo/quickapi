package quickapi

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func ListenAndServeWithGracefulShutdown(
	ctx context.Context,
	addr string, handler http.Handler,
	gracefulShutdownTimeout time.Duration,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	server := http.Server{
		Addr:    addr,
		Handler: handler,
	}
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
		defer cancel()
		server.Shutdown(ctx)
	}()
	log.Printf("[INFO]  listening: %s", server.Addr)
	return server.ListenAndServe()
}
