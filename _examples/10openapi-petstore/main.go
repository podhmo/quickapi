package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/experimental/define"
	rohandler "github.com/podhmo/reflect-openapi/handler"
)

var (
	gendoc bool
	port   int
)

func main() {
	flag.BoolVar(&gendoc, "gendoc", false, "generate openapi doc")
	flag.IntVar(&port, "port", 8080, "port")
	flag.Parse()
	if err := run(); err != nil {
		log.Fatalf("!! %+v", err)
	}

}

func run() error {
	ctx := context.Background()
	bc := build(port)

	if gendoc {
		if err := bc.EmitDoc(ctx, os.Stdout); err != nil {
			return err
		}
		return nil
	}

	handler, err := bc.BuildHandler(ctx)
	if err != nil {
		return err
	}
	bc.Router().Mount("/openapi", rohandler.NewHandler(bc.Doc(), "/openapi"))

	if err := quickapi.NewServer(fmt.Sprintf(":%d", port), handler, 5*time.Second).ListenAndServe(ctx); err != nil {
		log.Printf("[Error] !! %+v", err)
	}
	return nil
}

// see: https://github.com/deepmap/oapi-codegen/blob/master/examples/petstore-expanded/petstore-expanded.yaml

func build(port int) *define.BuildContext {
	doc := define.Doc().
		Title("Swagger Petstore").
		Version("1.0.0").
		Server("http://petstore.swagger.io/api", "").
		Server(fmt.Sprintf("http://localhost:%d", port), "local development")

	router := quickapi.DefaultRouter()
	bc := define.MustBuildContext(doc, router)

	define.DefaultError(bc, Error{})
	define.Type(bc, Pet{ID: "1", Name: "foo", Tag: "Cat"})

	longDescription := `Returns all pets from the system that the user has access to
	Nam sed condimentum est. Maecenas tempor sagittis sapien, nec rhoncus sem sagittis sit amet. Aenean at gravida augue, ac iaculis sem. Curabitur odio lorem, ornare eget elementum nec, cursus id lectus. Duis mi turpis, pulvinar ac eros ac, tincidunt varius justo. In hac habitasse platea dictumst. Integer at adipiscing ante, a sagittis ligula. Aenean pharetra tempor ante molestie imperdiet. Vivamus id aliquam diam. Cras quis velit non tortor eleifend sagittis. Praesent at enim pharetra urna volutpat venenatis eget eget mauris. In eleifend fermentum facilisis. Praesent enim enim, gravida ac sodales sed, placerat id erat. Suspendisse lacus dolor, consectetur non augue vel, vehicula interdum libero. Morbi euismod sagittis libero sed lacinia.
	Sed tempus felis lobortis leo pulvinar rutrum. Nam mattis velit nisl, eu condimentum ligula luctus nec. Phasellus semper velit eget aliquet faucibus. In a mattis elit. Phasellus vel urna viverra, condimentum lorem id, rhoncus nibh. Ut pellentesque posuere elementum. Sed a varius odio. Morbi rhoncus ligula libero, vel eleifend nunc tristique vitae. Fusce et sem dui. Aenean nec scelerisque tortor. Fusce malesuada accumsan magna vel tempus. Quisque mollis felis eu dolor tristique, sit amet auctor felis gravida. Sed libero lorem, molestie sed nisl in, accumsan tempor nisi. Fusce sollicitudin massa ut lacinia mattis. Sed vel eleifend lorem. Pellentesque vitae felis pretium, pulvinar elit eu, euismod sapien.
	`
	define.Get(bc, "/pets", FindPets).Description(longDescription)
	define.Post(bc, "/pets", AddPet).Description(`Creates a new pet in the store. Duplicates are allowed`).
		AnotherError(bc, 400, Error{}, "validation error").
		Example(400, "validation error", quickapi.ErrorResponse{Code: 400, Error: "validation error"})
	define.Get(bc, "/pets/{id}", FindPetByID).Description(`Returns a pet based on a single ID`)
	define.Delete(bc, "/pets/{id}", DeletePet).Description(`delete a single pet based on the ID supplied`).
		Status(204)
	return bc
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

// FindPets returns all pets
func FindPets(context.Context, struct {
	Tags  []string `query:"tags" openapi:"query"`  // tags to filter by. (style: form)
	Limit int32    `query:"limit" openapi:"query"` // maximum number of results to return (format: int32)
}) (
	output struct {
		Items []Pet `json:"items"`
	},
	err error,
) {
	return
}

// AddPet creates a new pet in the store. Duplicates are allowed
func AddPet(context.Context, struct {
	Name string `json:"name"`          // Name of the pet
	Tag  string `json:"tag,omitempty"` // Type of the pet
}) (
	output Pet,
	err error,
) {
	return
}

// FindPetByID returns a pet based on a single ID
func FindPetByID(context.Context, struct {
	ID int64 `path:"id" openapi:"path"` // ID of pet to fetch
}) (
	output Pet,
	err error,
) {
	return
}

// DeletePet deletes a pet by ID
func DeletePet(context.Context, struct {
	ID int64 `path:"id" openapi:"path"` // ID of pet to fetch
}) (
	output quickapi.Empty, // return 204
	err error,
) {
	return
}
