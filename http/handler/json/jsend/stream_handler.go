package jsend

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttphandlerjson "github.com/ralvarezdev/go-net/http/handler/json"
	gonethttpresponsehandlerjsend "github.com/ralvarezdev/go-net/http/response/handler/jsend"
)

// NewStreamHandler creates a new JSON handler in JSend format
//
// Parameters:
//
//   - mode: the flag mode
//
// Returns:
//
//   - *StreamHandler: the created JSON handler
//   - error: the error if any
func NewStreamHandler(mode *goflagsmode.Flag) (
	*gonethttphandlerjson.StreamHandler,
	error,
) {
	// Create the raw error handler for the responses handler
	rawErrorHandler := gonethttpresponsehandlerjsend.NewRawErrorHandler()

	return gonethttphandlerjson.NewStreamHandler(mode, rawErrorHandler)
}
