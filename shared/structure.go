package shared

// Empty is zero Input
type Empty struct{}

// ErrorResponse represents a normal error response type
type ErrorResponse struct {
	Code   int      `json:"code"`
	Error  string   `json:"error"`
	Detail []string `json:"detail,omitempty"`
}
