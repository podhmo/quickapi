package define

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/qopenapi/validate"
	"github.com/podhmo/quickapi/shared"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	"github.com/podhmo/reflect-openapi/docgen"
	"github.com/podhmo/reflect-openapi/dochandler"
	"github.com/podhmo/reflect-openapi/info"
)

type BuildContext struct {
	m *reflectopenapi.Manager
	c *reflectopenapi.Config

	r chi.Router

	mb         *validate.MiddlewareBuilder
	commitFunc func(context.Context) error
}

func NewBuildContext(docM DocModifier, r chi.Router, options ...func(c *reflectopenapi.Config)) (*BuildContext, error) {
	doc, loaded, err := docM()
	if err != nil {
		return nil, err
	}
	c := &reflectopenapi.Config{
		TagNameOption: &reflectopenapi.TagNameOption{
			NameTag:        "json",
			RequiredTag:    "required",
			ParamTypeTag:   "in",
			DescriptionTag: "description",
			OverrideTag:    "openapi-override",
			XNewTypeTag:    "x-go-type",
		},
		Doc:              doc,
		Loaded:           loaded,
		DefaultError:     shared.ErrorResponse{},
		StrictSchema:     true,
		EnableAutoTag:    true,
		Info:             info.New(),
		DisableInputRef:  true,
		DisableOutputRef: true,
	}

	if shared.FORCE {
		c.SkipValidation = true
	}
	for _, opt := range options {
		opt(c)
	}

	m, commit, err := c.NewManager()
	if err != nil {
		return nil, err
	}
	return &BuildContext{
		r:          r,
		c:          c,
		m:          m,
		mb:         validate.NewBuilder(doc, shared.DEBUG),
		commitFunc: commit,
	}, nil
}

func MustBuildContext(docM DocModifier, r chi.Router) *BuildContext {
	bc, err := NewBuildContext(docM, r)
	if err != nil {
		panic(err)
	}
	return bc
}

// Router returns internal github.com/go-chi/chi.Router
func (bc *BuildContext) Router() chi.Router {
	return bc.r
}

// Doc returns internal github.com/getkin/kin-openapi/openapi3.T
func (bc *BuildContext) Doc() *openapi3.T {
	return bc.m.Doc
}

// ReflectOpenAPIManager returns internal github.com/podhmo/reflect-openapi/Manager
func (bc *BuildContext) ReflectOpenAPIManager() *reflectopenapi.Manager {
	return bc.m
}

func (bc *BuildContext) EmitDoc(ctx context.Context, w io.Writer) error {
	if err := bc.commit(ctx); err != nil {
		return fmt.Errorf("EmitDoc (commit): %w", err)
	}

	// pathvars validation. (e.g. todo_id != id,  r.Get("/todos/{todo_id}", ...) with struct { ID string `in:"path" path:"id"`})
	if err := quickapi.WalkRoute(bc.r, func(item quickapi.RouteItem) error { return item.ValidatePathVars() }); err != nil {
		log.Printf("[WARN]  pathvars validation: %+v", err)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(bc.m.Doc); err != nil {
		return fmt.Errorf("emitDoc (json encode): %w", err)
	}
	return nil
}

func (bc *BuildContext) EmitMDDoc(ctx context.Context, w io.Writer) error {
	if err := bc.commit(ctx); err != nil {
		return fmt.Errorf("EmitMDDoc (commit): %w", err)
	}

	doc := docgen.Generate(bc.Doc(), bc.c.Info)
	if err := docgen.WriteDoc(w, doc); err != nil {
		return fmt.Errorf("emitMDDoc (json encode): %w", err)
	}
	return nil
}

func (bc *BuildContext) BuildHandler(ctx context.Context) (http.Handler, error) {
	if err := bc.commit(ctx); err != nil {
		return nil, fmt.Errorf("BuildHandler (commit): %w", err)
	}

	// pathvars validation. (e.g. todo_id != id,  r.Get("/todos/{todo_id}", ...) with struct { ID string `in:"path" path:"id"`})
	if err := quickapi.WalkRoute(bc.r, func(item quickapi.RouteItem) error { return item.ValidatePathVars() }); err != nil {
		return bc.r, fmt.Errorf("BuildHandler pathvars validation: %w", err)
	}
	return bc.r, nil
}

func (bc *BuildContext) BuildOpenAPIDoc(ctx context.Context) (*openapi3.T, error) {
	if err := bc.commit(ctx); err != nil {
		return nil, fmt.Errorf("BuildOpenAPIDoc (commit): %w", err)
	}
	return bc.m.Doc, nil
}

func (bc *BuildContext) BuildDocHandler(ctx context.Context, path string) (http.Handler, error) {
	if err := bc.commit(ctx); err != nil {
		return nil, fmt.Errorf("BuildDocHandler (commit): %w", err)
	}
	return dochandler.New(bc.Doc(), path, bc.c.Info), nil
}

func (bc *BuildContext) commit(ctx context.Context) error {
	if bc.commitFunc == nil {
		shared.GetLogger(ctx).Printf("[WARN]  already committed")
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
