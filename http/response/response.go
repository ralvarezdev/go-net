package response

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
)

type (
	// Response interface
	Response interface {
		Body(mode *goflagsmode.Flag) interface{}
		HTTPStatus() int
	}

	// SuccessResponse struct
	SuccessResponse struct {
		body       interface{}
		httpStatus int
	}

	// FailResponse struct
	FailResponse struct {
		body       interface{}
		httpStatus int
	}

	// ErrorResponse struct
	ErrorResponse struct {
		body       interface{}
		debugBody  interface{}
		httpStatus int
	}
)

// NewSuccessResponse creates a new success response
func NewSuccessResponse(
	body interface{},
	httpStatus int,
) *SuccessResponse {
	return &SuccessResponse{
		body,
		httpStatus,
	}
}

// NewJSendSuccessResponse creates a new JSend success response
func NewJSendSuccessResponse(
	data interface{},
	httpStatus int,
) *SuccessResponse {
	return NewSuccessResponse(NewJSendSuccessBody(data), httpStatus)
}

// Body returns the response body
func (r *SuccessResponse) Body(mode *goflagsmode.Flag) interface{} {
	return r.body
}

// HTTPStatus returns the HTTP status
func (r *SuccessResponse) HTTPStatus() int {
	return r.httpStatus
}

// NewFailResponse creates a new fail response
func NewFailResponse(
	body interface{},
	httpStatus int,
) *FailResponse {
	return &FailResponse{
		body,
		httpStatus,
	}
}

// NewJSendFailResponse creates a new JSend fail response
func NewJSendFailResponse(
	data interface{},
	errorCode *string,
	httpStatus int,
) *FailResponse {
	return NewFailResponse(NewJSendFailBody(data, errorCode), httpStatus)
}

// Body returns the response body
func (r *FailResponse) Body(mode *goflagsmode.Flag) interface{} {
	return r.body
}

// HTTPStatus returns the HTTP status
func (r *FailResponse) HTTPStatus() int {
	return r.httpStatus
}

// NewErrorResponse creates a new error response
func NewErrorResponse(
	body interface{},
	debugBody interface{},
	httpStatus int,
) *ErrorResponse {
	// Check if the debug body is nil
	if debugBody == nil {
		debugBody = body
	}

	return &ErrorResponse{
		body,
		debugBody,
		httpStatus,
	}
}

// NewJSendErrorResponse creates a new JSend error response
func NewJSendErrorResponse(
	err error,
	debugError error,
	data interface{},
	errorCode *string,
	httpStatus int,
) *ErrorResponse {
	// Create the response body
	body := NewJSendErrorBody(err.Error(), data, errorCode)

	// Create the response debug body
	var debugBody interface{}
	if debugError == nil {
		debugBody = body
	} else {
		debugBody = NewJSendErrorBody(debugError.Error(), data, errorCode)
	}

	return NewErrorResponse(body, debugBody, httpStatus)
}

// Body returns the response body
func (r *ErrorResponse) Body(mode *goflagsmode.Flag) interface{} {
	// Check if the mode is debug
	if mode != nil && mode.IsDebug() {
		return r.debugBody
	}
	return r.body
}

// HTTPStatus returns the HTTP status
func (r *ErrorResponse) HTTPStatus() int {
	return r.httpStatus
}
