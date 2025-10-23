package jsend

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"

	gojsondecoderjson "github.com/ralvarezdev/go-json/decoder/json"

	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
	gonethttprequesthandler "github.com/ralvarezdev/go-net/http/request/handler"
	gonethttprequestjson "github.com/ralvarezdev/go-net/http/request/json"
	gonethttpresponsehandler "github.com/ralvarezdev/go-net/http/response/handler"
	gonethttpresponsejson "github.com/ralvarezdev/go-net/http/response/json"
)

type (
	// Handler is the handler implementation for handling JSON requests and responses
	Handler struct {
		gonethttphandler.RequestsHandler
		gonethttphandler.ResponsesHandler
	}
)

// NewHandler creates a new JSON handler
//
// Parameters:
//
//   - mode: the flag mode
//   - rawErrorHandler: the raw error handler
//
// Returns:
//
//   - *Handler: the created JSON handler
//   - error: the error if any
func NewHandler(
	mode *goflagsmode.Flag,
	rawErrorHandler gonethttphandler.RawErrorHandler,
) (*Handler, error) {
	// Create the JSON encoder
	encoder := gonethttpresponsejson.NewEncoder(mode)

	// Create the responses handler
	responsesHandler, err := gonethttpresponsehandler.NewResponsesHandler(
		mode,
		encoder,
		rawErrorHandler,
	)
	if err != nil {
		return nil, err
	}

	// Create the JSON decoder
	decoder, err := gonethttprequestjson.NewDecoder(
		mode,
		gojsondecoderjson.NewDecoder(),
	)
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
		requestsHandler,
		responsesHandler,
	}, nil
}
