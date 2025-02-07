package response

type (
	// BaseJSendSuccessBody struct
	BaseJSendSuccessBody struct {
		Status string `json:"status"`
	}

	// JSendSuccessBody struct
	JSendSuccessBody struct {
		BaseJSendSuccessBody
		Data interface{} `json:"data,omitempty"`
	}

	// BaseJSendFailBody struct
	BaseJSendFailBody struct {
		Status string  `json:"status"`
		Code   *string `json:"code,omitempty"`
	}

	// JSendFailBody struct
	JSendFailBody struct {
		BaseJSendFailBody
		Data interface{} `json:"data,omitempty"`
	}

	// BaseJSendErrorBody struct
	BaseJSendErrorBody struct {
		Status  string  `json:"status"`
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
func NewBaseJSendSuccessBody() *BaseJSendSuccessBody {
	return &BaseJSendSuccessBody{
		Status: "success",
	}
}

// NewJSendSuccessBody creates a new JSend success response body
func NewJSendSuccessBody(
	data interface{},
) *JSendSuccessBody {
	return &JSendSuccessBody{
		BaseJSendSuccessBody: *NewBaseJSendSuccessBody(),
		Data:                 data,
	}
}

// NewJSendSuccessResponse creates a new JSend success response
func NewJSendSuccessResponse(
	data interface{},
	httpStatus int,
) Response {
	return NewResponse(NewJSendSuccessBody(data), httpStatus)
}

// NewBaseJSendFailBody creates a new base JSend fail response body
func NewBaseJSendFailBody(
	code *string,
) *BaseJSendFailBody {
	return &BaseJSendFailBody{
		Status: "fail",
		Code:   code,
	}
}

// NewJSendFailBody creates a new JSend fail response body
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
func NewJSendFailResponse(
	data interface{},
	code *string,
	httpStatus int,
) Response {
	return NewResponse(NewJSendFailBody(data, code), httpStatus)
}

// NewBaseJSendErrorBody creates a new base JSend error response body
func NewBaseJSendErrorBody(
	message string,
	code *string,
) *BaseJSendErrorBody {
	return &BaseJSendErrorBody{
		Status:  "error",
		Message: &message,
		Code:    code,
	}
}

// NewJSendErrorBody creates a new JSend error response body
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
func NewJSendErrorResponse(
	data interface{},
	message string,
	code *string,
	httpStatus int,
) Response {
	return NewResponse(NewJSendErrorBody(data, message, code), httpStatus)
}

// NewJSendErrorDebugResponse creates a new JSend error response with debug information
func NewJSendErrorDebugResponse(
	data interface{},
	message string,
	debugMessage string,
	code *string,
	httpStatus int,
) Response {
	return NewDebugResponse(
		NewJSendErrorBody(data, message, code),
		NewJSendErrorBody(data, debugMessage, code),
		httpStatus,
	)
}
