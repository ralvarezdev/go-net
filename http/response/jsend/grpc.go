package jsend

import (
	"context"
	"errors"
	"net/http"

	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

// ParseGRPCError parses a gRPC error to a JSend error response
//
// Parameters:
//
//   - err: the gRPC error
//
// Returns:
//   - *ErrorResponse: the JSend error response
func ParseGRPCError(
	err error,
) error {
	// Check if the error is nil
	if err == nil {
		return nil
	}

	// Check if the error comes from a gRPC status
	st, ok := status.FromError(err)
	if ok {
		// Try to parse detailed error information
		for _, detail := range st.Details() {
			switch info := detail.(type) {
			case *errdetails.BadRequest:
				return NewFailResponseFromErrorDetailsBadRequest(
					info,
					http.StatusBadRequest,
				)
			case *errdetails.PreconditionFailure:
				// Precondition failure error
				return gonethttpresponse.NewDebugErrorWithCode(
					err,
					ErrGRPCPreconditionFailed,
					ErrCodeGRPCPreconditionFailed,
					http.StatusPreconditionFailed,
				)
			case *errdetails.QuotaFailure:
				// Quota failure error
				return gonethttpresponse.NewDebugErrorWithCode(
					err,
					gonethttp.ErrTooManyRequests,
					ErrCodeGRPCQuotaFailure,
					http.StatusTooManyRequests,
				)
			case *errdetails.RequestInfo:
				// Request info error
				return gonethttpresponse.NewDebugErrorWithCode(
					err,
					gonethttp.ErrBadRequest,
					ErrCodeGRPCRequestInfo,
					http.StatusBadRequest,
				)
			case *errdetails.ResourceInfo:
				// Resource info error
				return gonethttpresponse.NewDebugErrorWithCode(
					err,
					gonethttp.ErrNotFound,
					ErrCodeGRPCResourceInfo,
					http.StatusNotFound,
				)
			case *errdetails.Help:
				// Help error
				return gonethttpresponse.NewDebugErrorWithCode(
					err,
					gonethttp.ErrBadRequest,
					ErrCodeGRPCHelp,
					http.StatusBadRequest,
				)
			case *errdetails.LocalizedMessage:
				// Localized message error
				return gonethttpresponse.NewDebugErrorWithCode(
					err,
					gonethttp.ErrBadRequest,
					ErrCodeGRPCLocalizedMessage,
					http.StatusBadRequest,
				)
			default:
				_ = info
			}
		}
		// gRPC status error
		return gonethttpresponse.NewDebugErrorWithCode(
			err,
			gonethttp.ErrBadRequest,
			ErrCodeGRPCCodePrefix+st.Code().String(),
			http.StatusBadRequest,
		)
	}
	if errors.Is(err, context.Canceled) {
		return gonethttpresponse.NewDebugErrorWithCode(
			err,
			gonethttp.ErrRequestTimeout,
			ErrCodeGRPCCtxCanceled,
			http.StatusRequestTimeout,
		)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return gonethttpresponse.NewDebugErrorWithCode(
			err,
			gonethttp.ErrRequestTimeout,
			ErrCodeGRPCCtxDeadlineExceeded,
			http.StatusRequestTimeout,
		)
	}

	// Other generic error
	return gonethttpresponse.NewDebugErrorWithCode(
		err,
		gonethttp.ErrInternalServerError,
		ErrCodeGRPCUnknown,
		http.StatusInternalServerError,
	)
}
