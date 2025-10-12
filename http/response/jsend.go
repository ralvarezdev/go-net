package response

type (
	// JSendSuccessBody struct
	JSendSuccessBody[T interface{}] struct {
		Status Status `json:"status"`
		Data   T      `json:"data,omitempty"`
	}

	// JSendFailBody struct
	JSendFailBody struct {
		Status Status      `json:"status"`
		Code   *string     `json:"code,omitempty"`
		Data   interface{} `json:"data,omitempty"`
	}

	// JSendErrorBody struct
	JSendErrorBody struct {
		Status  Status  `json:"status"`
		Message *string `json:"message,omitempty"`
		Code    *string `json:"code,omitempty"`
	}
)

// NewJSendSuccessBody creates a new JSend success response body
//
// Parameters:
//
//   - data: The data
//
// Returns:
//
//   - *JSendSuccessBody: The JSend success body
func NewJSendSuccessBody[T interface{}](
	data T,
) *JSendSuccessBody[T] {
	return &JSendSuccessBody[T]{
		Status: StatusSuccess,
		Data:   data,
	}
}

// NewJSendSuccessResponse creates a new JSend success response
//
// Parameters:
//
//   - data: The data
//   - httpStatus: The HTTP status code
//
// Returns:
//
//   - Response: The response
func NewJSendSuccessResponse(
	data interface{},
	httpStatus int,
) Response {
	return NewResponse(NewJSendSuccessBody(data), httpStatus)
}

// NewJSendFailBody creates a new JSend fail response body
//
// Parameters:
//
//   - data: The data
//   - code: The error code
//
// Returns:
//
//   - *JSendFailBody: The JSend fail body
func NewJSendFailBody(
	data interface{},
	code *string,
) *JSendFailBody {
	return &JSendFailBody{
		Status: StatusFail,
		Code:   code,
		Data:   data,
	}
}

// NewJSendFailResponse creates a new JSend fail response
//
// Parameters:
//
//   - data: The data
//   - code: The error code
//   - httpStatus: The HTTP status code
//
// Returns:
//
//   - Response: The response
func NewJSendFailResponse(
	data interface{},
	code *string,
	httpStatus int,
) Response {
	return NewResponse(NewJSendFailBody(data, code), httpStatus)
}

// NewJSendErrorBody creates a new JSend error response body
//
// Parameters:
//
//   - message: The error message
//   - code: The error code
//
// Returns:
//
//   - *JSendErrorBody: The JSend error body
func NewJSendErrorBody(
	message string,
	code *string,
) *JSendErrorBody {
	return &JSendErrorBody{
		Status:  StatusError,
		Message: &message,
		Code:    code,
	}
}

// NewJSendErrorResponse creates a new JSend error response
//
// Parameters:
//
//   - message: The error message
//   - code: The error code
//   - httpStatus: The HTTP status code
//
// Returns:
//
//   - Response: The response
func NewJSendErrorResponse(
	message string,
	code *string,
	httpStatus int,
) Response {
	return NewResponse(NewJSendErrorBody(message, code), httpStatus)
}

// NewJSendErrorDebugResponse creates a new JSend error response with debug information
//
// Parameters:
//
//   - message: The error message
//   - debugMessage: The debug message
//   - code: The error code
//   - httpStatus: The HTTP status code
//
// Returns:
//
//   - Response: The response
func NewJSendErrorDebugResponse(
	message string,
	debugMessage string,
	code *string,
	httpStatus int,
) Response {
	if debugMessage == "" {
		return NewJSendErrorResponse(message, code, httpStatus)
	}
	return NewDebugResponse(
		NewJSendErrorBody(message, code),
		NewJSendErrorBody(debugMessage, code),
		httpStatus,
	)
}
