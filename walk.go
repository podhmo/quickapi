package quickapi

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/podhmo/quickapi/internal/pathutil"
	"github.com/podhmo/quickapi/qbind"
)

type RouteItem struct {
	Method      string
	Route       string
	Metadata    qbind.Metadata
	PathVars    map[string]string // pathvar -> regex
	Handler     http.Handler
	Middlewares []func(http.Handler) http.Handler
}

func (item *RouteItem) ValidatePathVars() error {
	if len(item.Metadata.PathVars) == 0 {
		return nil
	}

	for _, taggedName := range item.Metadata.PathVars {
		if _, ok := item.PathVars[taggedName]; !ok {
			return fmt.Errorf("tagged name %q is not found (input is %v)", taggedName, item.Metadata.Input)
		}
	}
	for pathVar := range item.PathVars {
		found := false
		for _, taggedName := range item.Metadata.PathVars {
			if pathVar == taggedName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("pathvar %q is not found (input is %v)", pathVar, item.Metadata.Input)
		}
	}
	return nil
}

func WalkRoute(r chi.Router, fn func(RouteItem) error) error {
	return chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		_, _, vars := pathutil.NormalizeTemplatedPath(route)
		item := RouteItem{Method: method, Route: route, Handler: handler, Middlewares: middlewares, PathVars: vars}

		switch h := handler.(type) {
		case http.HandlerFunc:
			if metadata, tracked := metadataFromHandlerFunc(h); tracked {
				item.Metadata = metadata
				if err := fn(item); err != nil {
					return fmt.Errorf("route %s %s: %w", method, route, err)
				}
				return nil
			}
		case interface{ Metadata() qbind.Metadata }:
			item.Metadata = h.Metadata()
			if err := fn(item); err != nil {
				return fmt.Errorf("route %s %s: %w", method, route, err)
			}
			return nil
		}

		// outside of quickapi
		// log.Printf("ignored: %s %s", method, route)
		return nil
	})
}
