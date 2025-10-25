package jsend

import (
	"log/slog"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"

	gonethttphandlerprotojson "github.com/ralvarezdev/go-net/http/handler/protojson"
	gonethttpresponsehandlerjsend "github.com/ralvarezdev/go-net/http/response/handler/jsend"
)

// NewHandler creates a new protoJSON handler in JSend format
//
// Parameters:
//
//   - mode: the flag mode
//   - logger: the logger instance
//
// Returns:
//
//   - *Handler: the created protoJSON handler
//   - error: the error if any
func NewHandler(mode *goflagsmode.Flag, logger *slog.Logger) (
	*gonethttphandlerprotojson.Handler,
	error,
) {
	// Create the error raw handler for the responses handler
	rawErrorHandler := gonethttpresponsehandlerjsend.NewRawErrorHandler(logger)

	return gonethttphandlerprotojson.NewHandler(mode, rawErrorHandler)
}
