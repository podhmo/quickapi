package quickapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// A good base middleware stack https://github.com/go-chi/chi
var DEFAULT_MIDDLEWARES = []func(http.Handler) http.Handler{
	middleware.RequestID,
	middleware.RealIP,
	middleware.Logger,
	middleware.Recoverer,
}

func DefaultRouter() chi.Router {
	r := chi.NewRouter()
	for _, m := range DEFAULT_MIDDLEWARES {
		r.Use(m)
	}
	return r
}
