package response

import (
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpstatuserrors "github.com/ralvarezdev/go-net/http/status/errors"
	"net/http"
)

// NewJSendDebugInternalServerError creates a new internal server error JSend response with debug information
func NewJSendDebugInternalServerError(
	debugErr error,
	errorCode *string,
) gonethttpresponse.Response {
	return gonethttpresponse.NewJSendErrorResponse(
		gonethttpstatuserrors.InternalServerError,
		debugErr,
		nil,
		errorCode,
		http.StatusInternalServerError,
	)
}

// NewJSendInternalServerError creates a new internal server error JSend response
func NewJSendInternalServerError(
	errorCode *string,
) gonethttpresponse.Response {
	return gonethttpresponse.NewJSendErrorResponse(
		gonethttpstatuserrors.InternalServerError,
		nil,
		nil,
		errorCode,
		http.StatusInternalServerError,
	)
}

// NewJSendDebugNotImplemented creates a new not implemented JSend response with debug information
func NewJSendDebugNotImplemented(
	debugErr error,
	errorCode *string,
) gonethttpresponse.Response {
	return gonethttpresponse.NewJSendErrorResponse(
		gonethttpstatuserrors.NotImplemented,
		debugErr,
		nil,
		errorCode,
		http.StatusNotImplemented,
	)
}

// NewJSendNotImplemented creates a new not implemented JSend response
func NewJSendNotImplemented(
	errorCode *string,
) gonethttpresponse.Response {
	return gonethttpresponse.NewJSendErrorResponse(
		gonethttpstatuserrors.NotImplemented,
		nil,
		nil,
		errorCode,
		http.StatusNotImplemented,
	)
}
