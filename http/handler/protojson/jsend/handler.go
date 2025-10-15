package jsend

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
	gonethttprequesthandler "github.com/ralvarezdev/go-net/http/request/handler"
	gonethttprequestprotojson "github.com/ralvarezdev/go-net/http/request/protojson"
	gonethttpresponsehandlerjsend "github.com/ralvarezdev/go-net/http/response/handler/jsend"
	gonethttpresponseprotojson "github.com/ralvarezdev/go-net/http/response/protojson"
)

type (
	// Handler is the handler implementation for handling protoJSON requests and responses in JSend format
	Handler struct {
		gonethttphandler.RequestsHandler
		gonethttphandler.ResponsesHandler
	}
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
func NewHandler(mode *goflagsmode.Flag) (*Handler, error) {
	// Create the protoJSON encoder
	encoder := gonethttpresponseprotojson.NewEncoder(mode)

	// Create the responses handler
	responsesHandler, err := gonethttpresponsehandlerjsend.NewResponsesHandler(
		mode,
		encoder,
	)
	if err != nil {
		return nil, err
	}

	// Create the protoJSON decoder
	decoder := gonethttprequestprotojson.NewDecoder()

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
