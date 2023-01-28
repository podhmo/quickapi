package quickapi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/podhmo/quickapi/qbind"
)

type RouteItem struct {
	Method      string
	Route       string
	Metadata    qbind.Metadata
	PathVars    map[string]struct{}
	Handler     http.Handler
	Middlewares []func(http.Handler) http.Handler
}

func WalkRoute(r chi.Router, fn func(RouteItem) error) error {
	return chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		_, _, vars := normalizeTemplatedPath(route)
		item := RouteItem{Method: method, Route: route, Handler: handler, Middlewares: middlewares, PathVars: vars}
		switch h := handler.(type) {
		case http.HandlerFunc:
			item.Metadata = metadataFromHandlerFunc(h)
		case interface{ Metadata() qbind.Metadata }:
			item.Metadata = h.Metadata()
		}
		if err := fn(item); err != nil {
			return fmt.Errorf("route %s %s: %w", method, route, err)
		}
		return nil
	})
}

// copy from getkin/kin-openapi/openapi3/paths.go
func normalizeTemplatedPath(path string) (string, uint, map[string]struct{}) {
	if strings.IndexByte(path, '{') < 0 {
		return path, 0, nil
	}

	var buffTpl strings.Builder
	buffTpl.Grow(len(path))

	var (
		cc         rune
		count      uint
		isVariable bool
		vars       = make(map[string]struct{})
		buffVar    strings.Builder
	)
	for i, c := range path {
		if isVariable {
			if c == '}' {
				// End path variable
				isVariable = false

				vars[buffVar.String()] = struct{}{}
				buffVar = strings.Builder{}

				// First append possible '*' before this character
				// The character '}' will be appended
				if i > 0 && cc == '*' {
					buffTpl.WriteRune(cc)
				}
			} else {
				buffVar.WriteRune(c)
				continue
			}

		} else if c == '{' {
			// Begin path variable
			isVariable = true

			// The character '{' will be appended
			count++
		}

		// Append the character
		buffTpl.WriteRune(c)
		cc = c
	}
	return buffTpl.String(), count, vars
}
