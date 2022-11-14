package response

type errorResponse struct {
	Message string `json:"error_message"`
}

func (e errorResponse) AsMessage() string {
	return e.Message
}

func NewError(message string) *errorResponse {
	return &errorResponse{
		Message: message,
	}
}
