package response

import (
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpstatuserrors "github.com/ralvarezdev/go-net/http/status/errors"
	"net/http"
)

// NewDebugInternalServerError creates a new internal server error debug response
func NewDebugInternalServerError(
	debugErr error,
	errorCode *string,
) *gonethttpresponse.Response {
	return gonethttpresponse.NewDebugErrorResponse(
		gonethttpstatuserrors.InternalServerError,
		debugErr,
		nil,
		errorCode,
		http.StatusInternalServerError,
	)
}

// NewInternalServerError creates a new internal server error response
func NewInternalServerError(
	errorCode *string,
) *gonethttpresponse.Response {
	return gonethttpresponse.NewErrorResponse(
		gonethttpstatuserrors.InternalServerError,
		nil,
		errorCode,
		http.StatusInternalServerError,
	)
}
