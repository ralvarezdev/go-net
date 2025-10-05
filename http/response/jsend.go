package response

type (
	// BaseJSendSuccessBody struct
	BaseJSendSuccessBody struct {
		Status Status `json:"status"`
	}

	// JSendSuccessBody struct
	JSendSuccessBody struct {
		BaseJSendSuccessBody
		Data interface{} `json:"data,omitempty"`
	}

	// BaseJSendFailBody struct
	BaseJSendFailBody struct {
		Status Status  `json:"status"`
		Code   *string `json:"code,omitempty"`
	}

	// JSendFailBody struct
	JSendFailBody struct {
		BaseJSendFailBody
		Data interface{} `json:"data,omitempty"`
	}

	// BaseJSendErrorBody struct
	BaseJSendErrorBody struct {
		Status  Status  `json:"status"`
		Message *string `json:"message,omitempty"`
		Code    *string `json:"code,omitempty"`
	}

	// JSendErrorBody struct
	JSendErrorBody struct {
		BaseJSendErrorBody
		Data interface{} `json:"data,omitempty"`
	}
)

// NewBaseJSendSuccessBody creates a new base JSend success response body
//
// Returns:
//
//   - *BaseJSendSuccessBody: The base JSend success body
func NewBaseJSendSuccessBody() *BaseJSendSuccessBody {
	return &BaseJSendSuccessBody{
		Status: StatusSuccess,
	}
}

// NewJSendSuccessBody creates a new JSend success response body
//
// Parameters:
//
//   - data: The data
//
// Returns:
//
//   - *JSendSuccessBody: The JSend success body
func NewJSendSuccessBody(
	data interface{},
) *JSendSuccessBody {
	return &JSendSuccessBody{
		BaseJSendSuccessBody: *NewBaseJSendSuccessBody(),
		Data:                 data,
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

// NewBaseJSendFailBody creates a new base JSend fail response body
//
// Parameters:
//
//   - code: The error code
//
// Returns:
//
//   - *BaseJSendFailBody: The base JSend fail body
func NewBaseJSendFailBody(
	code *string,
) *BaseJSendFailBody {
	return &BaseJSendFailBody{
		Status: StatusFail,
		Code:   code,
	}
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
		BaseJSendFailBody: *NewBaseJSendFailBody(code),
		Data:              data,
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

// NewBaseJSendErrorBody creates a new base JSend error response body
//
// Parameters:
//
//   - message: The error message
//   - code: The error code
//
// Returns:
//
//   - *BaseJSendErrorBody: The base JSend error body
func NewBaseJSendErrorBody(
	message string,
	code *string,
) *BaseJSendErrorBody {
	return &BaseJSendErrorBody{
		Status:  StatusError,
		Message: &message,
		Code:    code,
	}
}

// NewJSendErrorBody creates a new JSend error response body
//
// Parameters:
//
//   - data: The data
//   - message: The error message
//   - code: The error code
//
// Returns:
//
//   - *JSendErrorBody: The JSend error body
func NewJSendErrorBody(
	data interface{},
	message string,
	code *string,
) *JSendErrorBody {
	return &JSendErrorBody{
		BaseJSendErrorBody: *NewBaseJSendErrorBody(message, code),
		Data:               data,
	}
}

// NewJSendErrorResponse creates a new JSend error response
//
// Parameters:
//
//   - data: The data
//   - message: The error message
//   - code: The error code
//   - httpStatus: The HTTP status code
//
// Returns:
//
//   - Response: The response
func NewJSendErrorResponse(
	data interface{},
	message string,
	code *string,
	httpStatus int,
) Response {
	return NewResponse(NewJSendErrorBody(data, message, code), httpStatus)
}

// NewJSendErrorDebugResponse creates a new JSend error response with debug information
//
// Parameters:
//
//   - data: The data
//   - message: The error message
//   - debugMessage: The debug message
//   - code: The error code
//   - httpStatus: The HTTP status code
//
// Returns:
//
//   - Response: The response
func NewJSendErrorDebugResponse(
	data interface{},
	message string,
	debugMessage string,
	code *string,
	httpStatus int,
) Response {
	if debugMessage == "" {
		return NewJSendErrorResponse(data, message, code, httpStatus)
	}
	return NewDebugResponse(
		NewJSendErrorBody(data, message, code),
		NewJSendErrorBody(data, debugMessage, code),
		httpStatus,
	)
}
