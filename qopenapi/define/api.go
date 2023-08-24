package define

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

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

type Config struct {
	ReflectOpenAPI *reflectopenapi.Config
	Validation     *validate.Config
}

type BuildContext struct {
	Config *Config

	m *reflectopenapi.Manager
	r chi.Router

	mb         *validate.MiddlewareBuilder
	onceCommit *onceCommit
}

type onceCommit struct {
	commitFunc func(context.Context) error
	err        error
	once       sync.Once
}

func (c *onceCommit) Do(ctx context.Context) error {
	c.once.Do(func() {
		c.err = c.commitFunc(ctx)
	})
	return c.err
}

func NewBuildContext(docM DocModifier, r chi.Router, options ...func(c *Config)) (*BuildContext, error) {
	doc, loaded, err := docM()
	if err != nil {
		return nil, err
	}

	mb := validate.NewBuilder(doc, shared.DEBUG)
	c := &Config{
		ReflectOpenAPI: &reflectopenapi.Config{
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
		},
		Validation: mb.Config,
	}

	for _, opt := range options {
		opt(c)
	}

	m, commit, err := c.ReflectOpenAPI.NewManager()
	if err != nil {
		return nil, err
	}
	return &BuildContext{
		Config:     c,
		r:          r,
		m:          m,
		mb:         mb,
		onceCommit: &onceCommit{commitFunc: commit},
	}, nil
}

func MustBuildContext(docM DocModifier, r chi.Router, options ...func(c *reflectopenapi.Config)) *BuildContext {
	bc, err := NewBuildContext(docM, r, options...)
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

func (bc *BuildContext) EmitDoc(ctx context.Context, filename string) error {
	if err := bc.onceCommit.Do(ctx); err != nil {
		return fmt.Errorf("EmitDoc (commit): %w", err)
	}

	// pathvars validation. (e.g. todo_id != id,  r.Get("/todos/{todo_id}", ...) with struct { ID string `in:"path" path:"id"`})
	if err := quickapi.WalkRoute(bc.r, func(item quickapi.RouteItem) error { return item.ValidatePathVars() }); err != nil {
		log.Printf("[WARN]  pathvars validation: %+v", err)
	}

	return writeFileOrStdout(ctx, filename, func(ctx context.Context, w io.Writer) error {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(bc.m.Doc); err != nil {
			return fmt.Errorf("emitDoc (json encode): %w", err)
		}
		return nil
	})
}

func (bc *BuildContext) EmitMDDoc(ctx context.Context, filename string) error {
	if err := bc.onceCommit.Do(ctx); err != nil {
		return fmt.Errorf("EmitMDDoc (commit): %w", err)
	}

	return writeFileOrStdout(ctx, filename, func(ctx context.Context, w io.Writer) error {
		c := bc.Config.ReflectOpenAPI
		doc := docgen.Generate(bc.Doc(), c.Info)
		if err := docgen.WriteDoc(w, doc); err != nil {
			return fmt.Errorf("emitMDDoc (writeDoc): %w", err)
		}
		return nil
	})
}

func (bc *BuildContext) BuildHandler(ctx context.Context) (http.Handler, error) {
	if err := bc.onceCommit.Do(ctx); err != nil {
		return nil, fmt.Errorf("BuildHandler (commit): %w", err)
	}

	// pathvars validation. (e.g. todo_id != id,  r.Get("/todos/{todo_id}", ...) with struct { ID string `in:"path" path:"id"`})
	if err := quickapi.WalkRoute(bc.r, func(item quickapi.RouteItem) error { return item.ValidatePathVars() }); err != nil {
		return bc.r, fmt.Errorf("BuildHandler pathvars validation: %w", err)
	}
	return bc.r, nil
}

func (bc *BuildContext) BuildOpenAPIDoc(ctx context.Context) (*openapi3.T, error) {
	if err := bc.onceCommit.Do(ctx); err != nil {
		return nil, fmt.Errorf("BuildOpenAPIDoc (commit): %w", err)
	}
	return bc.m.Doc, nil
}

func (bc *BuildContext) BuildDocHandler(ctx context.Context, path string, mdtext []byte) (http.Handler, error) {
	if err := bc.onceCommit.Do(ctx); err != nil {
		return nil, fmt.Errorf("BuildDocHandler (commit): %w", err)
	}
	c := bc.Config.ReflectOpenAPI
	return dochandler.New(bc.Doc(), path, c.Info, string(mdtext)), nil
}

func writeFileOrStdout(ctx context.Context, filename string, writeFunc func(context.Context, io.Writer) error) error {
	var w io.Writer = os.Stdout
	if filename != "" {
		f, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("create file %q: %w", filename, err)
		}
		defer f.Close()
		w = f
	}
	if err := writeFunc(ctx, w); err != nil {
		return fmt.Errorf("write file %q: %w", filename, err)
	}
	return nil
}

// ----------------------------------------
func DefaultError(bc *BuildContext, typ interface{}) {
	bc.Config.ReflectOpenAPI.DefaultError = typ
}
