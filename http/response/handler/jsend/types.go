package jsend

import (
	"errors"
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"

	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpresponsehandler "github.com/ralvarezdev/go-net/http/response/handler"
	gonethttpresponsejsend "github.com/ralvarezdev/go-net/http/response/jsend"
)

type (
	// RawErrorHandler struct
	RawErrorHandler struct{}
)

// NewRawErrorHandler creates a new default response handler
//
// Returns:
//
//   - *ResponsesHandler: The default handler
func NewRawErrorHandler() *RawErrorHandler {
	return &RawErrorHandler{}
}

// HandleRawError handles the raw error response
//
// Parameters:
//
//   - w: The HTTP response writer
//   - err: The error to handle
//   - handleResponseFn: The function to handle the response
func (r RawErrorHandler) HandleRawError(
	w http.ResponseWriter,
	err error,
	handleResponseFn func(
		w http.ResponseWriter,
		response gonethttpresponse.Response,
	),
) {
	{
		var parsedErr1 *gonethttpresponse.FailFieldError
		var parsedErr2 *gonethttpresponse.Error
		switch {
		case errors.As(err, &parsedErr1):
			handleResponseFn(w, gonethttpresponsejsend.NewResponseFromFailFieldError(parsedErr1))
		case errors.As(err, &parsedErr2):
			handleResponseFn(w, gonethttpresponsejsend.NewResponseFromError(parsedErr2))
		default:
			handleResponseFn(
				w,
				gonethttpresponsejsend.NewResponseFromError(
					gonethttpresponse.NewDebugErrorWithCode(
						err,
						gonethttp.ErrInternalServerError,
						ErrCodeRequestFatalError,
						http.StatusInternalServerError,
					),
				),
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
//
// Returns:
//
//   - *ResponsesHandler: The default handler
//   - error: The error if any
func NewResponsesHandler(
	mode *goflagsmode.Flag,
	encoder gonethttpresponse.Encoder,
) (*gonethttpresponsehandler.ResponsesHandler, error) {
	// Create the raw error handler
	rawErrorHandler := NewRawErrorHandler()

	// Create the responses handler
	return gonethttpresponsehandler.NewResponsesHandler(
		mode,
		encoder,
		rawErrorHandler,
	)
}
