package api

// APIResponse is used by insert/update/delete methods to indicate to the frontend whether the operation was successful.
// When marshalled to JSON it will give the front-end an array of errors, length 0 indicating all went fine.
// And optionally the ID of a created object.
type APIResponse struct {
	ID     uint64   `json:"id,omitempty"`
	Errors []string `json:"errors"`
}

// NewAPIResponse creates a new emptry API response.
func NewAPIResponse() *APIResponse {
	return &APIResponse{
		Errors: make([]string, 0),
	}
}

// AddErrorCode adds an error code to the APIResponse
func (cr *APIResponse) AddErrorCode(err string) {
	cr.Errors = append(cr.Errors, err)
}
