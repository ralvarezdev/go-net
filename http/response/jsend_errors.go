package response

import (
	"net/http"

	gonethttp "github.com/ralvarezdev/go-net/http"
)

// NewJSendDebugInternalServerError creates a new internal server error JSend response with debug information
//
// Parameters:
//
//   - debugErr: The debug error
//   - errorCode: The error code
//
// Returns:
//
//   - Response: The response
func NewJSendDebugInternalServerError(
	debugErr error,
	errorCode *string,
) Response {
	if debugErr == nil {
		return NewJSendInternalServerError(errorCode)
	}

	return NewJSendErrorDebugResponse(
		gonethttp.ErrInternalServerError.Error(),
		debugErr.Error(),
		errorCode,
		http.StatusInternalServerError,
	)
}

// NewJSendInternalServerError creates a new internal server error JSend response
//
// Parameters:
//
//   - errorCode: The error code
//
// Returns:
//
//   - Response: The response
func NewJSendInternalServerError(
	errorCode *string,
) Response {
	return NewJSendErrorResponse(
		gonethttp.ErrInternalServerError.Error(),
		errorCode,
		http.StatusInternalServerError,
	)
}

// NewJSendDebugNotImplemented creates a new not implemented JSend response with debug information
//
// Parameters:
//
//   - debugErr: The debug error
//   - errorCode: The error code
//
// Returns:
//
//   - Response: The response
func NewJSendDebugNotImplemented(
	debugErr error,
	errorCode *string,
) Response {
	if debugErr == nil {
		return NewJSendNotImplemented(errorCode)
	}

	return NewJSendErrorDebugResponse(
		gonethttp.ErrNotImplemented.Error(),
		debugErr.Error(),
		errorCode,
		http.StatusNotImplemented,
	)
}

// NewJSendNotImplemented creates a new not implemented JSend response
//
// Parameters:
//
//   - errorCode: The error code
//
// Returns:
//
//   - Response: The response
func NewJSendNotImplemented(
	errorCode *string,
) Response {
	return NewJSendErrorResponse(
		gonethttp.ErrNotImplemented.Error(),
		errorCode,
		http.StatusNotImplemented,
	)
}

// NewJSendDebugBadRequest creates a new bad request JSend response with debug information
//
// Parameters:
//
//   - debugErr: The debug error
//   - errorCode: The error code
//
// Returns:
//
//   - Response: The response
func NewJSendDebugBadRequest(
	debugErr error,
	errorCode *string,
) Response {
	if debugErr == nil {
		return NewJSendBadRequest(errorCode)
	}

	return NewJSendErrorDebugResponse(
		gonethttp.ErrBadRequest.Error(),
		debugErr.Error(),
		errorCode,
		http.StatusBadRequest,
	)
}

// NewJSendBadRequest creates a new bad request JSend response
//
// Parameters:
//
//   - errorCode: The error code
//
// Returns:
//
//   - Response: The response
func NewJSendBadRequest(
	errorCode *string,
) Response {
	return NewJSendErrorResponse(
		gonethttp.ErrBadRequest.Error(),
		errorCode,
		http.StatusBadRequest,
	)
}
