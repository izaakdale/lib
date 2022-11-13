package response

type Response struct {
	Message string `json:"error_message"`
}

func (e Response) AsMessage() string {
	return e.Message
}

func NewErrorResponse(message string) *Response {
	return &Response{
		Message: message,
	}
}
