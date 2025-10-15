package grpc

import (
	"context"
	"errors"
	"net/http"

	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

// ParseError parses a gRPC error to a JSend error response
//
// Parameters:
//
//   - err: the gRPC error
//   - parseAsValidations: whether to parse error details as validation errors
//
// Returns:
//   - *ErrorResponse: the JSend error response
func ParseError(
	err error,
	parseAsValidations bool,
) error {
	// Check if the error is nil
	if err == nil {
		return nil
	}

	// Check if the error comes from a gRPC status
	st, ok := status.FromError(err)
	if ok {
		// Try to parse detailed error information. The first matching detail type will be used.
		for _, detail := range st.Details() {
			switch info := detail.(type) {
			case *errdetails.BadRequest:
				return NewFailErrorFromErrorDetailsBadRequest(
					info,
					parseAsValidations,
				)
			case *errdetails.PreconditionFailure:
				return NewFailErrorFromErrorDetailsPreconditionFailure(
					info,
				)
			case *errdetails.QuotaFailure:
				return NewFailErrorFromErrorDetailsQuotaFailure(
					info,
				)
			case *errdetails.RequestInfo:
				return NewFailErrorFromErrorDetailsRequestInfo(
					info,
				)
			case *errdetails.ResourceInfo:
				return NewFailErrorFromErrorDetailsResourceInfo(
					info,
				)
			case *errdetails.Help:
				return NewFailErrorFromErrorDetailsHelp(
					info,
				)
			case *errdetails.LocalizedMessage:
				return NewFailErrorFromErrorDetailsLocalizedMessage(
					info,
				)
			default:
				_ = info
			}
		}
		// gRPC status error
		return gonethttpresponse.NewDebugErrorWithCode(
			err,
			gonethttp.ErrBadRequest,
			ErrCodeCCodePrefix+st.Code().String(),
			http.StatusBadRequest,
		)
	}
	if errors.Is(err, context.Canceled) {
		return gonethttpresponse.NewDebugErrorWithCode(
			err,
			gonethttp.ErrRequestTimeout,
			ErrCodeCtxCanceled,
			http.StatusRequestTimeout,
		)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return gonethttpresponse.NewDebugErrorWithCode(
			err,
			gonethttp.ErrRequestTimeout,
			ErrCodeCtxDeadlineExceeded,
			http.StatusRequestTimeout,
		)
	}

	// Other generic error
	return gonethttpresponse.NewDebugErrorWithCode(
		err,
		gonethttp.ErrInternalServerError,
		ErrCodeUnknown,
		http.StatusInternalServerError,
	)
}
