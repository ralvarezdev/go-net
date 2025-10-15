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
	// StreamHandler is the handler implementation for handling JSON stream requests and responses in JSend format
	StreamHandler struct {
		gonethttphandler.RequestsHandler
		gonethttphandler.ResponsesHandler
	}
)

// NewStreamHandler creates a new JSON handler in JSend format
//
// Parameters:
//
//   - mode: the flag mode
//   - logger: the logger
//
// Returns:
//
//   - *StreamHandler: the created JSON handler
//   - error: the error if any
func NewStreamHandler(mode *goflagsmode.Flag, logger *slog.Logger) (
	*Handler,
	error,
) {
	// Create the JSON stream encoder
	streamEncoder := gonethttpresponsejson.NewStreamEncoder(mode, logger)

	// Create the responses handler
	responsesHandler, err := gonethttpresponsehandlerjsend.NewResponsesHandler(
		mode,
		streamEncoder,
	)
	if err != nil {
		return nil, err
	}

	// Create the JSON stream decoder
	streamDecoder, err := gonethttprequestjson.NewDecoder(
		mode,
		responsesHandler,
	)
	if err != nil {
		return nil, err
	}

	// Create the requests handler
	requestsHandler, err := gonethttprequesthandler.NewDefaultRequestsHandler(
		mode,
		streamDecoder,
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
