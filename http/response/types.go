package response

type (
	// Response struct
	Response struct {
		Data interface{}
		Code *int
	}

	// JSONErrorResponse struct
	JSONErrorResponse struct {
		Error string `json:"error"`
	}
)

// NewResponseWithCode creates a new response with a code
func NewResponseWithCode(data interface{}, code int) *Response {
	return &Response{Data: data, Code: &code}
}

// NewErrorResponseWithCode creates a new error response with a code
func NewErrorResponseWithCode(err error, code int) *Response {
	return &Response{Data: NewJSONErrorResponse(err), Code: &code}
}

// NewResponse creates a new response
func NewResponse(data interface{}) *Response {
	return &Response{Data: data}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(err error) *Response {
	return &Response{Data: NewJSONErrorResponse(err)}
}

// NewJSONErrorResponse creates a new error response
func NewJSONErrorResponse(err error) JSONErrorResponse {
	return JSONErrorResponse{Error: err.Error()}
}
