package types

// HTTPError example used for error responses.
// swagger:model HTTPError
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
