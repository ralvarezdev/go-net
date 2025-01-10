package response

type (
	// Response struct
	Response struct {
		Response      interface{}
		DebugResponse interface{}
		HTTPStatus    int
	}

	// JSendResponse struct
	JSendResponse struct {
		Status  string      `json:"status"`
		Data    interface{} `json:"data"`
		Message *string     `json:"message,omitempty"`
		Code    *int        `json:"code,omitempty"`
	}
)

// NewJSendSuccessResponse creates a new success response
func NewJSendSuccessResponse(
	data interface{},
) *JSendResponse {
	return &JSendResponse{
		Status: "success",
		Data:   data,
	}
}

// NewJSendFailResponse creates a new fail response
func NewJSendFailResponse(
	data interface{},
	code *int,
) *JSendResponse {
	return &JSendResponse{
		Status: "fail",
		Data:   data,
		Code:   code,
	}
}

// NewJSendErrorResponse creates a new error response
func NewJSendErrorResponse(
	message string,
	data interface{},
	code *int,
) *JSendResponse {
	return &JSendResponse{
		Status:  "error",
		Data:    data,
		Message: &message,
		Code:    code,
	}
}

// newResponse creates a new response
func newResponse(
	response interface{},
	debugResponse interface{},
	httpStatus int,
) *Response {
	return &Response{
		Response:      response,
		DebugResponse: debugResponse,
		HTTPStatus:    httpStatus,
	}
}

// NewDebugSuccessResponse creates a new success response
func NewDebugSuccessResponse(
	data interface{},
	debugData interface{},
	httpStatus int,
) *Response {
	return newResponse(
		NewJSendSuccessResponse(data),
		NewJSendSuccessResponse(debugData),
		httpStatus,
	)
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(
	data interface{},
	httpStatus int,
) *Response {
	return NewDebugSuccessResponse(data, data, httpStatus)
}

// NewDebugFailResponse creates a new fail response
func NewDebugFailResponse(
	data interface{},
	debugData interface{},
	errorCode *int,
	httpStatus int,
) *Response {
	return newResponse(
		NewJSendFailResponse(data, errorCode),
		NewJSendFailResponse(debugData, errorCode),
		httpStatus,
	)
}

// NewFailResponse creates a new fail response
func NewFailResponse(
	data interface{},
	httpStatus int,
	errorCode *int,
) *Response {
	return NewDebugFailResponse(data, data, errorCode, httpStatus)
}

// NewDebugErrorResponse creates a new error response
func NewDebugErrorResponse(
	err error,
	debugErr error,
	data interface{},
	errorCode *int,
	httpStatus int,
) *Response {
	return newResponse(
		NewJSendErrorResponse(err.Error(), data, errorCode),
		NewJSendErrorResponse(debugErr.Error(), data, errorCode),
		httpStatus,
	)
}

// NewErrorResponse creates a new error response
func NewErrorResponse(
	err error,
	data interface{},
	errorCode *int,
	httpStatus int,
) *Response {
	return NewDebugErrorResponse(err, err, data, errorCode, httpStatus)
}
