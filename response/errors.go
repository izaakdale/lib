package response

type errorResponse struct {
	Message string `json:"error_message"`
}

// AsMessage returns the message from an errorResponse.
func (e errorResponse) AsMessage() string {
	return e.Message
}

// NewError takes a message and composes an errorResponse, to use when json tags are required.
func NewError(message string) *errorResponse {
	return &errorResponse{
		Message: message,
	}
}
