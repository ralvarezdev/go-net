package jsend

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
	gonethttprequesthandler "github.com/ralvarezdev/go-net/http/request/handler"
	gonethttprequestjson "github.com/ralvarezdev/go-net/http/request/json"
	gonethttpresponsehandler "github.com/ralvarezdev/go-net/http/response/handler"
	gonethttpresponsejson "github.com/ralvarezdev/go-net/http/response/json"
)

type (
	// StreamHandler is the handler implementation for handling JSON stream requests and responses
	StreamHandler struct {
		gonethttphandler.RequestsHandler
		gonethttphandler.ResponsesHandler
	}
)

// NewStreamHandler creates a new JSON handler
//
// Parameters:
//
//   - mode: the flag mode
//   - rawErrorhandler: the raw error handler
//
// Returns:
//
//   - *StreamHandler: the created JSON handler
//   - error: the error if any
func NewStreamHandler(
	mode *goflagsmode.Flag,
	rawErrorhandler gonethttphandler.RawErrorHandler,
) (
	*StreamHandler,
	error,
) {
	// Create the JSON stream encoder
	streamEncoder := gonethttpresponsejson.NewStreamEncoder(mode)

	// Create the responses handler
	responsesHandler, err := gonethttpresponsehandler.NewResponsesHandler(
		mode,
		streamEncoder,
		rawErrorhandler,
	)
	if err != nil {
		return nil, err
	}

	// Create the JSON stream decoder
	streamDecoder := gonethttprequestjson.NewStreamDecoder(
		mode,
	)

	// Create the requests handler
	requestsHandler, err := gonethttprequesthandler.NewDefaultRequestsHandler(
		mode,
		streamDecoder,
		responsesHandler,
	)
	if err != nil {
		return nil, err
	}

	return &StreamHandler{
		requestsHandler,
		responsesHandler,
	}, nil
}
