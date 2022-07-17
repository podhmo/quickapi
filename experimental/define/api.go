package define

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
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

	commitFunc func(context.Context) error
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
	return &BuildContext{r: r, c: &c, m: m, commitFunc: commit}, nil
}

func MustBuildContext(docM DocModifier, r chi.Router) *BuildContext {
	bc, err := NewBuildContext(docM, r)
	if err != nil {
		panic(err)
	}
	return bc
}

func (bc *BuildContext) EmitDoc(ctx context.Context, w io.Writer) error {
	if err := bc.commit(ctx); err != nil {
		return fmt.Errorf("EmitDoc (commit): %w", err)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(bc.m.Doc); err != nil {
		return fmt.Errorf("emitDoc (json encode): %w", err)
	}
	return nil
}

func (bc *BuildContext) BuildHandler(ctx context.Context) (chi.Router, error) {
	if err := bc.commit(ctx); err != nil {
		return nil, fmt.Errorf("BuildHandler (commit): %w", err)
	}
	return bc.r, nil
}

func (bc *BuildContext) BuildOpenAPIDoc(ctx context.Context) (*openapi3.T, error) {
	if err := bc.commit(ctx); err != nil {
		return nil, fmt.Errorf("BuildOpenAPIDoc (commit): %w", err)
	}
	return bc.m.Doc, nil
}

func (bc *BuildContext) commit(ctx context.Context) error {
	if bc.commitFunc == nil {
		log.Printf("[WARN]  already committed")
		return nil
	}
	defer func() { bc.commitFunc = nil }()
	commit := bc.commitFunc
	if err := commit(ctx); err != nil {
		return err
	}
	return nil
}

// ----------------------------------------
func DefaultError(bc *BuildContext, typ interface{}) {
	bc.c.DefaultError = typ
}
