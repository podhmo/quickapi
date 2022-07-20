package quickapi

import "github.com/podhmo/quickapi/shared"

var NewAPIError = shared.NewAPIError
var NoContent = shared.NoContent
var GetRequest = shared.GetRequest
var Redirect = shared.Redirect

// Empty is zero Input
type Empty = shared.Empty

// ErrorResponse represents a normal error response type
type ErrorResponse = shared.ErrorResponse
