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
	PathVars    map[string]string // pathvar -> regex
	Handler     http.Handler
	Middlewares []func(http.Handler) http.Handler
}

func (item *RouteItem) ValidatePathVars() error {
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

// arranged version of getkin/kin-openapi/openapi3/paths.go
func normalizeTemplatedPath(path string) (string, uint, map[string]string) {
	if strings.IndexByte(path, '{') < 0 {

		// add hoc support of 'foo/*'
		if strings.HasSuffix(path, "*") {
			return path[:len(path)-1] + "{STAR*}", 1, map[string]string{"STAR*": ""}
		}
		return path, 0, nil
	}

	var buffTpl strings.Builder
	buffTpl.Grow(len(path))

	var (
		count      uint
		isVariable bool
		isPattern  bool
		vars       = make(map[string]string)
		pattern    strings.Builder
		buffVar    strings.Builder
		lv         int
	)

	for i, c := range path {
		// log.Printf("c:%v\tisVariable:%v\tisPattern:%v\tlv:%d", string(c), isVariable, isPattern, lv)
		if isVariable {
			if c == '}' {
				lv--
				if lv == 0 {
					// End path variable
					isVariable = false
					isPattern = false

					vars[buffVar.String()] = pattern.String()
					buffVar = strings.Builder{}

				} else if isPattern {
					pattern.WriteRune(c)
					continue
				} else if c == ':' {
					isPattern = true
					continue
				} else {
					buffVar.WriteRune(c)
				}
			} else if isPattern {
				if c == '{' {
					lv++
				}
				pattern.WriteRune(c)
				continue
			} else if c == ':' {
				isPattern = true
				continue
			} else if c == ' ' {
				continue
			} else {
				buffVar.WriteRune(c)
			}
		} else if c == '{' {
			lv++
			if lv == 1 {
				// Begin path variable
				isVariable = true
				isPattern = false
				pattern = strings.Builder{}
				// The character '{' will be appended
				count++
			} else {
				buffVar.WriteRune(c)
			}
		} else if c == '*' && i == len(path)-1 {
			buffTpl.WriteString("{STAR*}")
			vars["STAR*"] = ""
			count++
			continue
		}

		// Append the character
		//log.Printf("\tP:%v\t%q", string(c), buffTpl.String())
		buffTpl.WriteRune(c)
	}
	return buffTpl.String(), count, vars
}
