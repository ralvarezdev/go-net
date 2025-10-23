package jsend

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"

	gonethttphandlerjson "github.com/ralvarezdev/go-net/http/handler/json"
	gonethttpresponsehandlerjsend "github.com/ralvarezdev/go-net/http/response/handler/jsend"
)

// NewHandler creates a new JSON handler in JSend format
//
// Parameters:
//
//   - mode: the flag mode
//
// Returns:
//
//   - *Handler: the created JSON handler
//   - error: the error if any
func NewHandler(mode *goflagsmode.Flag) (
	*gonethttphandlerjson.Handler,
	error,
) {
	// Create the raw error handler for the responses handler
	rawErrorHandler := gonethttpresponsehandlerjsend.NewRawErrorHandler()

	return gonethttphandlerjson.NewHandler(mode, rawErrorHandler)
}
