package handler

type ErrorResponse struct {
	error string
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		error: err.Error(),
	}
}

func (er ErrorResponse) getJSON() []byte {
	return []byte("{\"error\": \"" + er.error + "\"}")
}
