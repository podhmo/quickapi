//go:generate go run ./ -gendoc -docfile openapi.json -mdfile apidoc.md
package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"go/token"
	"log"
	"time"

	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/qopenapi/define"
	reflectshape "github.com/podhmo/reflect-shape"
)

//go:embed openapi.json
var openapiDocData []byte

//go:embed apidoc.md
var mdDocData []byte

var options struct {
	gendoc  bool
	docfile string
	mdfile  string
	port    int
}

func main() {
	flag.BoolVar(&options.gendoc, "gendoc", false, "generate openapi doc")
	flag.IntVar(&options.port, "port", 8080, "port")
	flag.StringVar(&options.docfile, "docfile", "", "file name of openapi doc. if this value is empty output to stdout.")
	flag.StringVar(&options.mdfile, "mdfile", "", "")
	flag.Parse()
	if err := run(); err != nil {
		log.Fatalf("!! %+v", err)
	}

}

func run() error {
	ctx := context.Background()

	doc := define.Doc().
		Title("Swagger Petstore").
		Description("This is a sample server Petstore server.  You can find out more about Swagger at [http://swagger.io](http://swagger.io) or on [irc.freenode.net, #swagger](http://swagger.io/irc/).  For this sample, you can use the api key `special-key` to test the authorization filters.").
		Version("1.0.0").
		Server("http://petstore.swagger.io/api", "").
		Server(fmt.Sprintf("http://localhost:%d", options.port), "local development")

	if !options.gendoc {
		doc = doc.LoadFromData(openapiDocData)
	}

	router := quickapi.DefaultRouter()
	bc, err := define.NewBuildContext(doc, router, func(c *define.Config) {
		c.ReflectOpenAPI.GoPositionFunc = func(fset *token.FileSet, f *reflectshape.Func) string {
			// TODO: multiple package
			fpos := fset.Position(f.Pos())
			return fmt.Sprintf("https://github.com/podhmo/quickapi/blob/main/_examples/10openapi-petstore/main.go#L%d", fpos.Line)
		}
	})
	if err != nil {
		return err
	}

	define.DefaultError(bc, Error{})

	mount(bc)

	if options.gendoc {
		if err := bc.EmitDoc(ctx, options.docfile); err != nil {
			return err
		}

		if options.mdfile != "" {
			if err := bc.EmitMDDoc(ctx, options.mdfile); err != nil {
				return err
			}
		}
		return nil
	}

	handler, err := bc.BuildHandler(ctx)
	if err != nil {
		return err
	}
	dochandler, err := bc.BuildDocHandler(ctx, "/_doc", mdDocData)
	if err != nil {
		return err
	}
	bc.Router().Mount("/_doc", dochandler)

	if err := quickapi.NewServer(fmt.Sprintf(":%d", options.port), handler, 5*time.Second).ListenAndServe(ctx); err != nil {
		log.Printf("[Error] !! %+v", err)
	}
	return nil
}

// see: https://github.com/deepmap/oapi-codegen/blob/master/examples/petstore-expanded/petstore-expanded.yaml

func mount(bc *define.BuildContext) {
	cat := Pet{ID: "1", Name: "foo", Tag: "Cat"}
	// define.Type(bc, cat) // as default
	define.Type(bc, Pet{}).Example(cat)

	{
		api := &PetAPI{}
		longDescription := `Returns all pets from the system that the user has access to
		Nam sed condimentum est. Maecenas tempor sagittis sapien, nec rhoncus sem sagittis sit amet. Aenean at gravida augue, ac iaculis sem. Curabitur odio lorem, ornare eget elementum nec, cursus id lectus. Duis mi turpis, pulvinar ac eros ac, tincidunt varius justo. In hac habitasse platea dictumst. Integer at adipiscing ante, a sagittis ligula. Aenean pharetra tempor ante molestie imperdiet. Vivamus id aliquam diam. Cras quis velit non tortor eleifend sagittis. Praesent at enim pharetra urna volutpat venenatis eget eget mauris. In eleifend fermentum facilisis. Praesent enim enim, gravida ac sodales sed, placerat id erat. Suspendisse lacus dolor, consectetur non augue vel, vehicula interdum libero. Morbi euismod sagittis libero sed lacinia.
		Sed tempus felis lobortis leo pulvinar rutrum. Nam mattis velit nisl, eu condimentum ligula luctus nec. Phasellus semper velit eget aliquet faucibus. In a mattis elit. Phasellus vel urna viverra, condimentum lorem id, rhoncus nibh. Ut pellentesque posuere elementum. Sed a varius odio. Morbi rhoncus ligula libero, vel eleifend nunc tristique vitae. Fusce et sem dui. Aenean nec scelerisque tortor. Fusce malesuada accumsan magna vel tempus. Quisque mollis felis eu dolor tristique, sit amet auctor felis gravida. Sed libero lorem, molestie sed nisl in, accumsan tempor nisi. Fusce sollicitudin massa ut lacinia mattis. Sed vel eleifend lorem. Pellentesque vitae felis pretium, pulvinar elit eu, euismod sapien.
		`

		define.Get(bc, "/pets", api.FindPets).Description(longDescription)
		define.Post(bc, "/pets", api.AddPet).
			AnotherError(bc, 400, Error{}, "validation error").
			Example(400, "bad request", Error{Code: 400, Message: "validation error"}).
			DefaultInput(func() AddPetInput { return AddPetInput{Name: "foo"} })
		define.Get(bc, "/pets/{id}", api.FindPetByID)
		define.Delete(bc, "/pets/{id}", api.DeletePet).Status(204)
	}

	define.GetHTML(bc, "/hello/{name}", Hello, nil)
}

type Pet struct { // allOf is not supported
	ID   string `json:"id"`                                       // unique id of the pet
	Name string `json:"name" openapi-override:"{'minLength': 1}"` // name of the pet
	Tag  string `json:"tag,omitempty"`                            // id of the pet
}

type Error struct {
	Code    int32  `json:"code"`              // Error code
	Message string `json:"message,omitempty"` // message
}

type PetAPI struct {
}

// FindPets returns all pets
func (api *PetAPI) FindPets(context.Context, struct {
	Tags  []string `query:"tags" in:"query"`  // tags to filter by. (style: form)
	Limit int32    `query:"limit" in:"query"` // maximum number of results to return (format: int32)
}) (
	output struct { // list of pets
		Items []Pet `json:"items"`
	},
	err error,
) {
	return
}

type AddPetInput struct {
	Name string `json:"name" openapi-override:"{'minLength': 1}"` // Name of the pet
	Tag  string `json:"tag,omitempty"`                            // Type of the pet

	Pretty bool `json:"pretty" in:"query" query:"pretty"` // pretty output (hmm)
}

// AddPet creates a new pet in the store. Duplicates are allowed
func (api *PetAPI) AddPet(ctx context.Context, input AddPetInput) (
	output Pet,
	err error,
) {
	return Pet{Name: input.Name, Tag: input.Tag}, nil
}

// FindPetByID returns a pet based on a single ID
func (api *PetAPI) FindPetByID(context.Context, struct {
	ID int64 `path:"id" in:"path"` // ID of pet to fetch
}) (
	output Pet,
	err error,
) {
	return
}

// DeletePet deletes a pet by ID
func (api *PetAPI) DeletePet(context.Context, struct {
	ID int64 `path:"id" in:"path"` // ID of pet to fetch
}) (
	output quickapi.Empty, // return 204
	err error,
) {
	return
}

func Hello(ctx context.Context, input struct {
	Name string `in:"path" path:"name"`
}) (string /* return greeting text */, error) {
	return fmt.Sprintf("hello %s", input.Name), nil
}
