---
title: Swagger Petstore
version: 1.0.0
---

# Swagger Petstore

This is a sample server Petstore server.  You can find out more about Swagger at [http://swagger.io](http://swagger.io) or on [irc.freenode.net, #swagger](http://swagger.io/irc/).  For this sample, you can use the api key `special-key` to test the authorization filters.

- [paths](#paths)
- [schemas](#schemas)

## paths

| endpoint | operationId | tags | summary |
| --- | --- | --- | --- |
| `GET /pets` | [main.PetAPI.FindPets](#mainpetapifindpets-get-pets)  |  | FindPets returns all pets |
| `POST /pets` | [main.PetAPI.AddPet](#mainpetapiaddpet-post-pets)  |  | AddPet creates a new pet in the store. Duplicates are allowed |
| `DELETE /pets/{id}` | [main.PetAPI.DeletePet](#mainpetapideletepet-delete-petsid)  |  | DeletePet deletes a pet by ID |
| `GET /pets/{id}` | [main.PetAPI.FindPetByID](#mainpetapifindpetbyid-get-petsid)  |  | FindPetByID returns a pet based on a single ID |


### main.PetAPI.FindPets `GET /pets`

FindPets returns all pets

| name | value | 
| --- | --- |
| operationId | main.PetAPI.FindPets |
| endpoint | `GET /pets` |
| tags |  |


#### input (application/json)

```go
// GET /pets
type Input struct {
	tags? []string `in:"query"`
	limit? integer `in:"query"`
}
```

#### output (application/json)

```go

// GET /pets (200)
// list of pets
type Output200 struct {	// 
	items []struct {	// Pet
		// unique id of the pet
		id string

		// name of the pet
		name string `minLength:"1"`

		// id of the pet
		tag? string
	}
}

// GET /pets (default)
// default error
type OutputDefault struct {	// Error
	// Error code
	code integer `format:"int32"`

	// message
	message? string
}
```

#### description

Returns all pets from the system that the user has access to
		Nam sed condimentum est. Maecenas tempor sagittis sapien, nec rhoncus sem sagittis sit amet. Aenean at gravida augue, ac iaculis sem. Curabitur odio lorem, ornare eget elementum nec, cursus id lectus. Duis mi turpis, pulvinar ac eros ac, tincidunt varius justo. In hac habitasse platea dictumst. Integer at adipiscing ante, a sagittis ligula. Aenean pharetra tempor ante molestie imperdiet. Vivamus id aliquam diam. Cras quis velit non tortor eleifend sagittis. Praesent at enim pharetra urna volutpat venenatis eget eget mauris. In eleifend fermentum facilisis. Praesent enim enim, gravida ac sodales sed, placerat id erat. Suspendisse lacus dolor, consectetur non augue vel, vehicula interdum libero. Morbi euismod sagittis libero sed lacinia.
		Sed tempus felis lobortis leo pulvinar rutrum. Nam mattis velit nisl, eu condimentum ligula luctus nec. Phasellus semper velit eget aliquet faucibus. In a mattis elit. Phasellus vel urna viverra, condimentum lorem id, rhoncus nibh. Ut pellentesque posuere elementum. Sed a varius odio. Morbi rhoncus ligula libero, vel eleifend nunc tristique vitae. Fusce et sem dui. Aenean nec scelerisque tortor. Fusce malesuada accumsan magna vel tempus. Quisque mollis felis eu dolor tristique, sit amet auctor felis gravida. Sed libero lorem, molestie sed nisl in, accumsan tempor nisi. Fusce sollicitudin massa ut lacinia mattis. Sed vel eleifend lorem. Pellentesque vitae felis pretium, pulvinar elit eu, euismod sapien.
### main.PetAPI.AddPet `POST /pets`

AddPet creates a new pet in the store. Duplicates are allowed

| name | value | 
| --- | --- |
| operationId | main.PetAPI.AddPet |
| endpoint | `POST /pets` |
| tags |  |


#### input (application/json)

```go
// POST /pets
type Input struct {
	// pretty output (hmm)
	pretty? boolean `in:"query"`

	JSONBody struct {	// AddPetInput
		// Name of the pet
		name string

		// Type of the pet
		tag? string
	}
}
```

#### output (application/json)

```go

// POST /pets (200)
type Output200 struct {	// Pet
	// unique id of the pet
	id string

	// name of the pet
	name string `minLength:"1"`

	// id of the pet
	tag? string
}

// POST /pets (400)
// -
type Output400 struct {	// Error
	// Error code
	code integer `format:"int32"`

	// message
	message? string
}

// POST /pets (default)
// default error
type OutputDefault struct {	// Error
	// Error code
	code integer `format:"int32"`

	// message
	message? string
}
```

examples

```json

// POST /pets (400)
// bad request
{
  "code": 400,
  "message": "validation error"
}
```

#### description

AddPet creates a new pet in the store. Duplicates are allowed
### main.PetAPI.DeletePet `DELETE /pets/{id}`

DeletePet deletes a pet by ID

| name | value | 
| --- | --- |
| operationId | main.PetAPI.DeletePet |
| endpoint | `DELETE /pets/{id}` |
| tags |  |


#### input (application/json)

```go
// DELETE /pets/{id}
type Input struct {
	id integer `in:"path"`
}
```

#### output (application/json)

```go

// DELETE /pets/{id} (204)
// return 204
type Output204 struct {	// Empty
}

// DELETE /pets/{id} (default)
// default error
type OutputDefault struct {	// Error
	// Error code
	code integer `format:"int32"`

	// message
	message? string
}
```

#### description

DeletePet deletes a pet by ID
### main.PetAPI.FindPetByID `GET /pets/{id}`

FindPetByID returns a pet based on a single ID

| name | value | 
| --- | --- |
| operationId | main.PetAPI.FindPetByID |
| endpoint | `GET /pets/{id}` |
| tags |  |


#### input (application/json)

```go
// GET /pets/{id}
type Input struct {
	id integer `in:"path"`
}
```

#### output (application/json)

```go

// GET /pets/{id} (200)
type Output200 struct {	// Pet
	// unique id of the pet
	id string

	// name of the pet
	name string `minLength:"1"`

	// id of the pet
	tag? string
}

// GET /pets/{id} (default)
// default error
type OutputDefault struct {	// Error
	// Error code
	code integer `format:"int32"`

	// message
	message? string
}
```

#### description

FindPetByID returns a pet based on a single ID



----------------------------------------

## schemas

| name | summary |
| --- | --- |
| [AddPetInput](#addpetinput) |  |
| [Empty](#empty) | is zero Input |
| [Error](#error) |  |
| [Pet](#pet) |  |



### AddPetInput

```go
type AddPetInput struct {
	// Name of the pet
	name string

	// Type of the pet
	tag? string
}
```

- [input of main.PetAPI.AddPet](#mainpetapiaddpet-post-pets)

### Empty

```go
// Empty is zero Input
type Empty struct {
}
```

- [output of main.PetAPI.DeletePet (204)](#mainpetapideletepet-delete-petsid)

### Error

```go
type Error struct {
	// Error code
	code integer `format:"int32"`

	// message
	message? string
}
```

- [output of main.PetAPI.FindPets (default)](#mainpetapifindpets-get-pets)
- [output of main.PetAPI.AddPet (400)](#mainpetapiaddpet-post-pets)
- [output of main.PetAPI.AddPet (default)](#mainpetapiaddpet-post-pets)
- [output of main.PetAPI.DeletePet (default)](#mainpetapideletepet-delete-petsid)
- [output of main.PetAPI.FindPetByID (default)](#mainpetapifindpetbyid-get-petsid)

### Pet

```go
type Pet struct {
	// unique id of the pet
	id string

	// name of the pet
	name string `minLength:"1"`

	// id of the pet
	tag? string
}
```

exmaples

```js
// 
{
  "id": "1",
  "name": "foo",
  "tag": "Cat"
}
```

- [output of main.PetAPI.AddPet (200)](#mainpetapiaddpet-post-pets)
- [output of main.PetAPI.FindPetByID (200)](#mainpetapifindpetbyid-get-petsid)