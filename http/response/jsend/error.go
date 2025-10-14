package jsend

import (
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	// ErrorBody struct
	ErrorBody struct {
		Status  Status `json:"status"`
		Message string `json:"message,omitempty"`
		Code    string `json:"code,omitempty"`
	}
)

// NewErrorBodyWithCode creates a new JSend error response body with error code
//
// Parameters:
//
//   - message: The error message
//   - code: The error code
//
// Returns:
//
//   - *ErrorBody: The JSend error body
func NewErrorBodyWithCode(
	message string,
	code string,
) *ErrorBody {
	return &ErrorBody{
		Status:  StatusError,
		Message: message,
		Code:    code,
	}
}

// NewErrorBody creates a new JSend error response body
//
// Parameters:
//
//   - message: The error message
//
// Returns:
//
//   - *ErrorBody: The JSend error body
func NewErrorBody(
	message string,
) *ErrorBody {
	return NewErrorBodyWithCode(message, "")
}

// NewErrorResponseWithCode creates a new JSend error response with error code
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
func NewErrorResponseWithCode(
	message string,
	code string,
	httpStatus int,
) gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(
		NewErrorBodyWithCode(message, code),
		httpStatus,
	)
}

// NewErrorResponse creates a new JSend error response
//
// Parameters:
//
//   - message: The error message
//   - httpStatus: The HTTP status code
//
// Returns:
//
//   - Response: The response
func NewErrorResponse(
	message string,
	httpStatus int,
) gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(
		NewErrorBody(message),
		httpStatus,
	)
}

// NewDebugErrorResponseWithCode creates a new JSend error response with debug information and error code
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
func NewDebugErrorResponseWithCode(
	message string,
	debugMessage string,
	code string,
	httpStatus int,
) gonethttpresponse.Response {
	if debugMessage == "" {
		return NewErrorResponseWithCode(message, code, httpStatus)
	}
	return gonethttpresponse.NewDebugResponse(
		NewErrorBodyWithCode(message, code),
		NewErrorBodyWithCode(debugMessage, code),
		httpStatus,
	)
}

// NewDebugErrorResponse creates a new JSend error response with debug information
//
// Parameters:
//
//   - message: The error message
//   - debugMessage: The debug message
//   - httpStatus: The HTTP status code
//
// Returns:
//
//   - Response: The response
func NewDebugErrorResponse(
	message string,
	debugMessage string,
	httpStatus int,
) gonethttpresponse.Response {
	return NewDebugErrorResponseWithCode(message, debugMessage, "", httpStatus)
}
