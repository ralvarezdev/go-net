package handler

type (
	// Response struct
	Response struct {
		Data      interface{}
		DebugData interface{}
		Code      *int
	}

	// JSONErrorResponse struct
	JSONErrorResponse struct {
		Error string `json:"error"`
	}
)

// NewJSONErrorResponse creates a new error response
func NewJSONErrorResponse(err error) JSONErrorResponse {
	return JSONErrorResponse{Error: err.Error()}
}

// newResponse creates a new response
func newResponse(
	data interface{},
	debugData interface{},
	code *int,
) *Response {
	return &Response{Data: data, DebugData: debugData, Code: code}
}

// NewDebugResponseWithCode creates a new response with a code
func NewDebugResponseWithCode(
	data interface{},
	debugData interface{},
	code int,
) *Response {
	return newResponse(data, debugData, &code)
}

// NewDebugErrorResponseWithCode creates a new error response with a code
func NewDebugErrorResponseWithCode(
	err error,
	debugErr error,
	code int,
) *Response {
	return newResponse(
		NewJSONErrorResponse(err),
		NewJSONErrorResponse(debugErr),
		&code,
	)
}

// NewResponseWithCode creates a new response with a code
func NewResponseWithCode(data interface{}, code int) *Response {
	return newResponse(data, nil, &code)
}

// NewErrorResponseWithCode creates a new error response with a code
func NewErrorResponseWithCode(err error, code int) *Response {
	return newResponse(NewJSONErrorResponse(err), nil, &code)
}

// NewDebugResponse creates a new response
func NewDebugResponse(data interface{}, debugData interface{}) *Response {
	return newResponse(data, debugData, nil)
}

// NewDebugErrorResponse creates a new error response
func NewDebugErrorResponse(err error, debugErr error) *Response {
	return newResponse(
		NewJSONErrorResponse(err),
		NewJSONErrorResponse(debugErr),
		nil,
	)
}

// NewResponse creates a new response
func NewResponse(data interface{}) *Response {
	return newResponse(data, nil, nil)
}

// NewErrorResponse creates a new error response
func NewErrorResponse(err error) *Response {
	return newResponse(NewJSONErrorResponse(err), nil, nil)
}
