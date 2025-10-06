package response

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(message string, data interface{}) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string) Response {
	return Response{
		Success: false,
		Message: message,
		Data:    nil,
	}
}