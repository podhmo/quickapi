package define

import (
	"context"
	_ "embed"
	"encoding/json"
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

//go:embed skeleton.json
var docSkeleton []byte

type BuildContext struct {
	m *reflectopenapi.Manager
	c *reflectopenapi.Config

	r chi.Router

	commit func(context.Context) error
}

func NewBuildContext(r chi.Router) (*BuildContext, error) {
	doc, err := reflectopenapi.NewDocFromSkeleton(docSkeleton)
	if err != nil {
		return nil, err
	}

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

func EmitDoc(ctx context.Context, bc *BuildContext) error {
	if err := bc.commit(ctx); err != nil {
		return err
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(bc.m.Doc); err != nil {
		return err
	}
	return nil
}

// ----------------------------------------
func DefaultError(bc *BuildContext, typ interface{}) {
	bc.c.DefaultError = typ
}
