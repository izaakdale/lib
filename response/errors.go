package response

type Error struct {
	Message string `json:"error_message"`
}

func (e Error) AsMessage() string {
	return e.Message
}

func NewErrorResponse(message string) *Error {
	return &Error{
		Message: message,
	}
}
