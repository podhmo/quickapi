package define

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/go-chi/chi/v5"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

type APIError struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type BuildContext struct {
	m *reflectopenapi.Manager
	c *reflectopenapi.Config

	r chi.Router

	commit func(context.Context) error
}

func NewBuildContext(docM DocModifier, r chi.Router) (*BuildContext, error) {
	doc := docM()
	c := reflectopenapi.Config{
		Doc:          doc,
		DefaultError: APIError{},
		StrictSchema: true,
		IsRequiredCheckFunction: func(tag reflect.StructTag) bool {
			required := true
			if val, ok := tag.Lookup("openapi"); ok && val != "body" {
				required = false
			}
			if _, isOptional := tag.Lookup("optional"); isOptional {
				required = false
			}
			if val, ok := tag.Lookup("json"); ok && strings.Contains(val, "omitempty") {
				required = false
			}
			return required
		},
	}

	m, commit, err := c.NewManager()
	if err != nil {
		return nil, err
	}
	return &BuildContext{r: r, c: &c, m: m, commit: commit}, nil
}

func MustBuildContext(docM DocModifier, r chi.Router) *BuildContext {
	bc, err := NewBuildContext(docM, r)
	if err != nil {
		panic(err)
	}
	return bc
}

func (bc *BuildContext) EmitDoc(ctx context.Context) error {
	if err := bc.commit(ctx); err != nil {
		return fmt.Errorf("emitDoc (commit): %w", err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(bc.m.Doc); err != nil {
		return fmt.Errorf("emitDoc (json encode): %w", err)
	}
	return nil
}
func (bc *BuildContext) Handler() chi.Router {
	return bc.r
}

// ----------------------------------------
func DefaultError(bc *BuildContext, typ interface{}) {
	bc.c.DefaultError = typ
}
