package common

import "github.com/lib/pq"

// APIResponse is used by insert/update/delete methods to indicate to the frontend whether the operation was successful.
// When marshalled to JSON it will give the front-end an array of errors, length 0 indicating all went fine.
// And optionally the ID of a created object.
type APIResponse struct {
	ID     uint64   `json:"id,omitempty"`
	Errors []string `json:"errors"`
}

// NewAPIResponse creates a new emptry API response.
func NewAPIResponse() APIResponse {
	return APIResponse{
		Errors: make([]string, 0),
	}
}

// CheckPgErr checks an error. If it's a postgres error of a kind that is handled in the front-end,
// an error code is added to the API response object and the function returns true.
// TODO: accept a whitelist of database constraints and violations that can go through to the frontend.
func (cr *APIResponse) CheckPgErr(err error) bool {
	if perr, ok := err.(*pq.Error); ok {
		switch perr.Code.Name() {
		case `unique_violation`:
			cr.Errors = append(cr.Errors, `pgerr:`+perr.Constraint)
			return true
		case `check_violation`:
			cr.Errors = append(cr.Errors, `pgerr:check_violation:`+perr.Constraint)
			return true
		case `datetime_field_overflow`:
			cr.Errors = append(cr.Errors, `pgerr:datetime_field_overflow`)
			return true
		}
	}
	return false
}

// AddErrorCode adds an error code to the APIResponse
func (cr *APIResponse) AddErrorCode(err string) {
	cr.Errors = append(cr.Errors, err)
}
