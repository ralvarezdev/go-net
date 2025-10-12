package response

import (
	"net/http"

	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpstatuserrors "github.com/ralvarezdev/go-net/http/status/errors"
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
//   - gonethttpresponse.Response: The response
func NewJSendDebugInternalServerError(
	debugErr error,
	errorCode *string,
) gonethttpresponse.Response {
	if debugErr == nil {
		return NewJSendInternalServerError(errorCode)
	}

	return gonethttpresponse.NewJSendErrorDebugResponse(
		gonethttpstatuserrors.InternalServerError.Error(),
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
//   - gonethttpresponse.Response: The response
func NewJSendInternalServerError(
	errorCode *string,
) gonethttpresponse.Response {
	return gonethttpresponse.NewJSendErrorResponse(
		gonethttpstatuserrors.InternalServerError.Error(),
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
//   - gonethttpresponse.Response: The response
func NewJSendDebugNotImplemented(
	debugErr error,
	errorCode *string,
) gonethttpresponse.Response {
	if debugErr == nil {
		return NewJSendNotImplemented(errorCode)
	}

	return gonethttpresponse.NewJSendErrorDebugResponse(
		gonethttpstatuserrors.NotImplemented.Error(),
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
//   - gonethttpresponse.Response: The response
func NewJSendNotImplemented(
	errorCode *string,
) gonethttpresponse.Response {
	return gonethttpresponse.NewJSendErrorResponse(
		gonethttpstatuserrors.NotImplemented.Error(),
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
//   - gonethttpresponse.Response: The response
func NewJSendDebugBadRequest(
	debugErr error,
	errorCode *string,
) gonethttpresponse.Response {
	if debugErr == nil {
		return NewJSendBadRequest(errorCode)
	}

	return gonethttpresponse.NewJSendErrorDebugResponse(
		gonethttpstatuserrors.BadRequest.Error(),
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
//   - gonethttpresponse.Response: The response
func NewJSendBadRequest(
	errorCode *string,
) gonethttpresponse.Response {
	return gonethttpresponse.NewJSendErrorResponse(
		gonethttpstatuserrors.BadRequest.Error(),
		errorCode,
		http.StatusBadRequest,
	)
}
