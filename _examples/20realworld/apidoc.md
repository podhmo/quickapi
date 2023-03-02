---
title: Conduit API
version: 1.0.0
---

# Conduit API

Conduit API documentation

- [paths](#paths)
- [schemas](#schemas)

## paths

| endpoint | operationId | tags | summary |
| --- | --- | --- | --- |







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
