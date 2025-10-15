package jsend

import (
	"log/slog"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
	gonethttprequesthandler "github.com/ralvarezdev/go-net/http/request/handler"
	gonethttprequestjson "github.com/ralvarezdev/go-net/http/request/json"
	gonethttpresponsehandlerjsend "github.com/ralvarezdev/go-net/http/response/handler/jsend"
	gonethttpresponsejson "github.com/ralvarezdev/go-net/http/response/json"
)

type (
	// Handler is the handler implementation for handling JSON requests and responses in JSend format
	Handler struct {
		gonethttphandler.RequestsHandler
		gonethttphandler.ResponsesHandler
	}
)

// NewHandler creates a new JSON handler in JSend format
//
// Parameters:
//
//   - mode: the flag mode
//   - logger: the logger
//
// Returns:
//
//   - *Handler: the created JSON handler
//   - error: the error if any
func NewHandler(mode *goflagsmode.Flag, logger *slog.Logger) (*Handler, error) {
	// Create the JSON encoder
	encoder := gonethttpresponsejson.NewEncoder(mode, logger)

	// Create the responses handler
	responsesHandler, err := gonethttpresponsehandlerjsend.NewResponsesHandler(
		mode,
		encoder,
	)
	if err != nil {
		return nil, err
	}

	// Create the JSON decoder
	decoder, err := gonethttprequestjson.NewDecoder(mode, responsesHandler)
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
