package shared

// Empty is zero Input
type Empty struct{}

// ErrorResponse represents a normal error response type
type ErrorResponse struct {
	Code   int      `json:"code"`
	Error  string   `json:"error"`
	Detail []string `json:"detail,omitempty"`
}

// NoContent is similar to Empty but has status code is changed
type noContent struct {
	code int
}

func (ob noContent) StatusCode() int {
	return ob.code
}

func NoContent(code int) StatusCoder {
	return noContent{code: code}
}
