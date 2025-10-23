package jsend

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"

	gonethttphandlerprotojson "github.com/ralvarezdev/go-net/http/handler/protojson"
	gonethttpresponsehandlerjsend "github.com/ralvarezdev/go-net/http/response/handler/jsend"
)

// NewHandler creates a new protoJSON handler in JSend format
//
// Parameters:
//
//   - mode: the flag mode
//
// Returns:
//
//   - *Handler: the created protoJSON handler
//   - error: the error if any
func NewHandler(mode *goflagsmode.Flag) (
	*gonethttphandlerprotojson.Handler,
	error,
) {
	// Create the error raw handler for the responses handler
	rawErrorHandler := gonethttpresponsehandlerjsend.NewRawErrorHandler()

	return gonethttphandlerprotojson.NewHandler(mode, rawErrorHandler)
}
