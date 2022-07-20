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

func NewServer(
	addr string, handler http.Handler,
	gracefulShutdownTimeout time.Duration,
) *Server {
	return &Server{
		Server: http.Server{
			Addr:    addr,
			Handler: handler,

			// hmm
			ReadHeaderTimeout: 2 * time.Second,
			IdleTimeout:       30 * time.Second,
		},
		GracefulShutdownTimeout: gracefulShutdownTimeout,
	}
}

func (server *Server) ListenAndServe(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), server.GracefulShutdownTimeout)
		defer cancel()
		server.Shutdown(ctx)
	}()
	log.Printf("[INFO]  listening: %s", server.Addr)
	return server.Server.ListenAndServe()
}
