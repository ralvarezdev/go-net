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
	return gonethttpresponse.NewJSendErrorDebugResponse(
		nil,
		gonethttpstatuserrors.InternalServerError.Error(),
		debugErr.Error(),
		errorCode,
		http.StatusInternalServerError,
	)
}

// NewJSendInternalServerError creates a new internal server error JSend response
func NewJSendInternalServerError(
	errorCode *string,
) gonethttpresponse.Response {
	return gonethttpresponse.NewJSendErrorResponse(
		nil,
		gonethttpstatuserrors.InternalServerError.Error(),
		errorCode,
		http.StatusInternalServerError,
	)
}

// NewJSendDebugNotImplemented creates a new not implemented JSend response with debug information
func NewJSendDebugNotImplemented(
	debugErr error,
	errorCode *string,
) gonethttpresponse.Response {
	return gonethttpresponse.NewJSendErrorDebugResponse(
		nil,
		gonethttpstatuserrors.NotImplemented.Error(),
		debugErr.Error(),
		errorCode,
		http.StatusNotImplemented,
	)
}

// NewJSendNotImplemented creates a new not implemented JSend response
func NewJSendNotImplemented(
	errorCode *string,
) gonethttpresponse.Response {
	return gonethttpresponse.NewJSendErrorResponse(
		nil,
		gonethttpstatuserrors.NotImplemented.Error(),
		errorCode,
		http.StatusNotImplemented,
	)
}

// NewJSendDebugBadRequest creates a new bad request JSend response with debug information
func NewJSendDebugBadRequest(
	debugErr error,
	errorCode *string,
) gonethttpresponse.Response {
	return gonethttpresponse.NewJSendErrorDebugResponse(
		nil,
		gonethttpstatuserrors.BadRequest.Error(),
		debugErr.Error(),
		errorCode,
		http.StatusBadRequest,
	)
}

// NewJSendBadRequest creates a new bad request JSend response
func NewJSendBadRequest(
	errorCode *string,
) gonethttpresponse.Response {
	return gonethttpresponse.NewJSendErrorResponse(
		nil,
		gonethttpstatuserrors.BadRequest.Error(),
		errorCode,
		http.StatusBadRequest,
	)
}
