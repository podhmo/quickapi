---
title: hello example
version: 0.0.0
---

# hello example

hello example

- [paths](#paths)
- [schemas](#schemas)

## paths

| endpoint | operationId | tags | summary |
| --- | --- | --- | --- |
| `POST /main.Hello` | [main.Hello](#mainhello-post-mainhello)  | `main` | hello, greeting message |


### main.Hello `POST /main.Hello`

hello, greeting message

| name | value |
| --- | --- |
| operationId | main.Hello[  <sub>(source)</sub>](https://github.com/podhmo/quickapi/blob/main/_examples/11tiny-rpc/main.go#L106) |
| endpoint | `POST /main.Hello` |
| input | Input[ [`HelloInput`](#helloinput) ] |
| output | [`HelloOutput`](#hellooutput) ｜ [`ErrorResponse`](#errorresponse) |
| tags | `main` |


#### input (application/json)

```go
// POST /main.Hello
type Input struct {
	JSONBody struct {	// HelloInput
		name string
	}
}
```

#### output (application/json)

```go
// POST /main.Hello (200)
type Output200 struct {	// HelloOutput
	message string
}

// POST /main.Hello (default)
// default error
type OutputDefault struct {	// ErrorResponse
	code integer

	error string

	detail? []string
}
```

#### description

hello, greeting message





----------------------------------------

## schemas

| name | summary |
| --- | --- |
| [ErrorResponse](#errorresponse) | represents a normal error response type |



### ErrorResponse

represents a normal error response type

```go
// ErrorResponse represents a normal error response type
type ErrorResponse struct {
	code integer

	error string

	detail? []string
}
```

- [output of main.Hello (default) as `ErrorResponse`](#mainhello-post-mainhello)