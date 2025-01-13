package response

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gostringsconvert "github.com/ralvarezdev/go-strings/convert"
)

type (
	// Response struct
	Response struct {
		body       interface{}
		debugBody  interface{}
		httpStatus int
	}

	// JSendBody struct
	JSendBody struct {
		Status  string      `json:"status"`
		Data    interface{} `json:"data"`
		Message *string     `json:"message,omitempty"`
		Code    *int        `json:"code,omitempty"`
	}
)

// NewJSendSuccessBody creates a new success response body
func NewJSendSuccessBody(
	data interface{},
) *JSendBody {
	return &JSendBody{
		Status: "success",
		Data:   data,
	}
}

// NewJSendFailBody creates a new fail response body
func NewJSendFailBody(
	data interface{},
	code *int,
) *JSendBody {
	return &JSendBody{
		Status: "fail",
		Data:   data,
		Code:   code,
	}
}

// NewJSendErrorBody creates a new error response body
func NewJSendErrorBody(
	message string,
	data interface{},
	code *int,
) *JSendBody {
	return &JSendBody{
		Status:  "error",
		Data:    data,
		Message: &message,
		Code:    code,
	}
}

// newResponse creates a new response
func newResponse(
	body interface{},
	debugBody interface{},
	httpStatus int,
) *Response {
	return &Response{
		body:       body,
		debugBody:  debugBody,
		httpStatus: httpStatus,
	}
}

// GetBody returns the response body
func (r *Response) GetBody(mode *goflagsmode.Flag) interface{} {
	// Check if the response contains the debug response body
	if r.debugBody != nil && mode != nil && mode.IsDebug() {
		return r.debugBody
	}
	return r.body
}

// GetHTTPStatus returns the HTTP status
func (r *Response) GetHTTPStatus() int {
	return r.httpStatus
}

// NewDebugSuccessResponse creates a new success response
func NewDebugSuccessResponse(
	data interface{},
	debugData interface{},
	httpStatus int,
) *Response {
	return newResponse(
		NewJSendSuccessBody(data),
		NewJSendSuccessBody(debugData),
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
		NewJSendFailBody(data, errorCode),
		NewJSendFailBody(debugData, errorCode),
		httpStatus,
	)
}

// NewFailResponse creates a new fail response
func NewFailResponse(
	data interface{},
	errorCode *int,
	httpStatus int,
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
		NewJSendErrorBody(err.Error(), data, errorCode),
		NewJSendErrorBody(debugErr.Error(), data, errorCode),
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

// NewSingleFieldBodyData creates a new single field body data
func NewSingleFieldBodyData(
	fieldName string,
	fieldValue ...interface{},
) *map[string]interface{} {
	return &map[string]interface{}{
		fieldName: &[]interface{}{fieldValue},
	}
}

// NewSingleFieldErrorsBodyData creates a new single field errors body data
func NewSingleFieldErrorsBodyData(
	fieldName string,
	fieldValue ...error,
) *map[string]*[]string {
	return &map[string]*[]string{
		fieldName: gostringsconvert.ErrorArrayToStringArray(&fieldValue),
	}
}
