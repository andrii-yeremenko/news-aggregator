package handler

// ErrorResponse is a structured error.
type ErrorResponse struct {
	error string
}

// NewErrorResponse creates a new ErrorResponse instance.
func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		error: err.Error(),
	}
}

// getJSON returns a JSON representation of the error.
func (er ErrorResponse) getJSON() []byte {
	return []byte("{\"error\": \"" + er.error + "\"}")
}
