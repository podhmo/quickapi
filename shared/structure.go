package shared

type ErrorResponse struct {
	Code   int      `json:"code"`
	Error  string   `json:"error"`
	Detail []string `json:"detail,omitempty"`
}
