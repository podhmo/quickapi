package main

import (
	"context"
	"log"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/experimental/define"
)

// see: https://github.com/deepmap/oapi-codegen/blob/master/examples/petstore-expanded/petstore-expanded.yaml

type Pet struct { // allOf is not supported
	ID   string `json:"id"`            // unique id of the pet
	Name string `json:"name"`          // name of the pet
	Tag  string `json:"tag,omitempty"` // id of the pet
}

type Error struct {
	Code    int32  `json:"code"`              // Error code
	Message string `json:"message,omitempty"` // message
}

////////////////////////////////////////

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
	ID int64 `query:"id" openapi:"path"` // ID of pet to fetch
}) (
	output Pet,
	err error,
) {
	return
}

// DeletePet deletes a pet by ID
func DeletePet(context.Context, struct {
	ID int64 `query:"id" openapi:"path"` // ID of pet to fetch
}) (
	output quickapi.Empty, // return 204
	err error,
) {
	return
}

func main() {
	bc, err := define.NewBuildContext(chi.NewRouter())
	if err != nil {
		log.Fatalf("!! %+v", err)
	}

	define.DefaultError(bc, Error{})

	longDescription := `Returns all pets from the system that the user has access to
	Nam sed condimentum est. Maecenas tempor sagittis sapien, nec rhoncus sem sagittis sit amet. Aenean at gravida augue, ac iaculis sem. Curabitur odio lorem, ornare eget elementum nec, cursus id lectus. Duis mi turpis, pulvinar ac eros ac, tincidunt varius justo. In hac habitasse platea dictumst. Integer at adipiscing ante, a sagittis ligula. Aenean pharetra tempor ante molestie imperdiet. Vivamus id aliquam diam. Cras quis velit non tortor eleifend sagittis. Praesent at enim pharetra urna volutpat venenatis eget eget mauris. In eleifend fermentum facilisis. Praesent enim enim, gravida ac sodales sed, placerat id erat. Suspendisse lacus dolor, consectetur non augue vel, vehicula interdum libero. Morbi euismod sagittis libero sed lacinia.
	Sed tempus felis lobortis leo pulvinar rutrum. Nam mattis velit nisl, eu condimentum ligula luctus nec. Phasellus semper velit eget aliquet faucibus. In a mattis elit. Phasellus vel urna viverra, condimentum lorem id, rhoncus nibh. Ut pellentesque posuere elementum. Sed a varius odio. Morbi rhoncus ligula libero, vel eleifend nunc tristique vitae. Fusce et sem dui. Aenean nec scelerisque tortor. Fusce malesuada accumsan magna vel tempus. Quisque mollis felis eu dolor tristique, sit amet auctor felis gravida. Sed libero lorem, molestie sed nisl in, accumsan tempor nisi. Fusce sollicitudin massa ut lacinia mattis. Sed vel eleifend lorem. Pellentesque vitae felis pretium, pulvinar elit eu, euismod sapien.
	`
	define.Get(bc, "/pets", FindPets).Description(strings.TrimSpace(longDescription))
	define.Post(bc, "/pets", AddPet).Description(strings.TrimSpace(`Creates a new pet in the store. Duplicates are allowed`))
	define.Get(bc, "/pets/{id}", FindPetByID).Description(strings.TrimSpace(`Returns a pet based on a single ID`))
	define.Delete(bc, "/pets/{id}", DeletePet).Description(strings.TrimSpace(`delete a single pet based on the ID supplied`))

	ctx := context.Background()
	if err := define.EmitDoc(ctx, bc); err != nil {
		log.Fatalf("!! %+v", err)
	}
}
