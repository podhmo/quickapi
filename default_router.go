package quickapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
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

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		render.Status(req, 404)
		v := errorResponse{Error: http.StatusText(404), Code: 404}
		render.JSON(w, req, v)
	})
	return r
}
