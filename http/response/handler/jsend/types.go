package jsend

import (
	"errors"
	"log/slog"
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"

	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpresponsehandler "github.com/ralvarezdev/go-net/http/response/handler"
	gonethttpresponsejsend "github.com/ralvarezdev/go-net/http/response/jsend"
)

type (
	// RawErrorHandler struct
	RawErrorHandler struct {
		logger *slog.Logger
	}
)

// NewRawErrorHandler creates a new default response handler
//
// Parameters:
//
//   - logger: The logger instance
//
// Returns:
//
//   - *ResponsesHandler: The default handler
func NewRawErrorHandler(logger *slog.Logger) *RawErrorHandler {
	if logger != nil {
		logger = logger.With(
			slog.String("component", "http_response_jsend_raw_error_handler"),
		)
	}

	return &RawErrorHandler{logger}
}

// HandleRawError handles the raw error response
//
// Parameters:
//
//   - w: The HTTP response writer
//   - req: The HTTP request
//   - err: The error to handle
//   - stackTrace: The error debug stack trace
//   - stackTrace: The error debug stack trace
//   - handleResponseFn: The function to handle the response
func (r RawErrorHandler) HandleRawError(
	w http.ResponseWriter,
	req *http.Request,
	err error,
	stackTrace []byte,
	handleResponseFn func(
		w http.ResponseWriter,
		req *http.Request,
		response gonethttpresponse.Response,
	),
) {
	var parsedErr1 *gonethttpresponse.FailFieldError
	var parsedErr2 *gonethttpresponse.Error
	switch {
	case errors.As(err, &parsedErr1):
		handleResponseFn(w, req, gonethttpresponsejsend.NewResponseFromFailFieldError(parsedErr1))
	case errors.As(err, &parsedErr2):
		handleResponseFn(w, req, gonethttpresponsejsend.NewResponseFromError(parsedErr2))
	default:
		handleResponseFn(
			w,
			req,
			gonethttpresponsejsend.NewResponseFromError(
				gonethttpresponse.NewDebugErrorWithCode(
					err,
					gonethttp.ErrInternalServerError,
					ErrCodeRequestFatalError,
					http.StatusInternalServerError,
				),
			),
		)

		// Log the unhandled error
		if r.logger == nil {
			return
		}

		if stackTrace != nil {
			r.logger.Error(
				"An unhandled error caught in RawErrorHandler",
				slog.Any("error", err),
				slog.String("route", req.URL.Path),
				slog.String("stack_trace", string(stackTrace)),
			)
		} else {
			r.logger.Error(
				"An unhandled error caught in RawErrorHandler",
				slog.Any("error", err),
				slog.String("route", req.URL.Path),
			)
		}
	}
}

// NewResponsesHandler creates a new default response handler
//
// Parameters:
//
//   - mode: The flag mode
//   - encoder: The HTTP response encoder
//   - logger: The logger instance
//
// Returns:
//
//   - *ResponsesHandler: The default handler
//   - error: The error if any
func NewResponsesHandler(
	mode *goflagsmode.Flag,
	encoder gonethttpresponse.Encoder,
	logger *slog.Logger,
) (*gonethttpresponsehandler.ResponsesHandler, error) {
	// Create the raw error handler
	rawErrorHandler := NewRawErrorHandler(logger)

	// Create the responses handler
	return gonethttpresponsehandler.NewResponsesHandler(
		mode,
		encoder,
		rawErrorHandler,
	)
}
