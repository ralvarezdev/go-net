package jsend

import (
	"log/slog"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttprequesthandler "github.com/ralvarezdev/go-net/http/request/handler"
	gonethttprequestprotojson "github.com/ralvarezdev/go-net/http/request/protojson"
	gonethttpresponsehandler "github.com/ralvarezdev/go-net/http/response/handler"
	gonethttpresponsehandlerjsend "github.com/ralvarezdev/go-net/http/response/handler/jsend"
	gonethttpresponseprotojson "github.com/ralvarezdev/go-net/http/response/protojson"
)

type (
	// Handler is the handler implementation for handling protoJSON requests and responses in JSend format
	Handler struct {
		gonethttprequesthandler.RequestsHandler
		gonethttpresponsehandler.ResponsesHandler
	}
)

// NewHandler creates a new protoJSON handler in JSend format
//
// Parameters:
//
//   - mode: the flag mode
//   - logger: the logger
//
// Returns:
//
//   - *Handler: the created protoJSON handler
//   - error: the error if any
func NewHandler(mode *goflagsmode.Flag, logger *slog.Logger) (*Handler, error) {
	// Create the protoJSON encoder
	encoder := gonethttpresponseprotojson.NewEncoder(mode, logger)

	// Create the responses handler
	responsesHandler, err := gonethttpresponsehandlerjsend.NewResponsesHandler(
		mode,
		encoder,
	)
	if err != nil {
		return nil, err
	}

	// Create the protoJSON decoder
	decoder, err := gonethttprequestprotojson.NewDecoder(responsesHandler)
	if err != nil {
		return nil, err
	}

	// Create the requests handler
	requestsHandler, err := gonethttprequesthandler.NewDefaultRequestsHandler(
		mode,
		decoder,
		responsesHandler,
	)
	if err != nil {
		return nil, err
	}

	return &Handler{
		*requestsHandler,
		*responsesHandler,
	}, nil
}
