package validator

import (
	gonethttpjson "github.com/ralvarezdev/go-net/http/json"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"net/http"
)

type (
	// Handler interface
	Handler interface {
		HandleError(
			w http.ResponseWriter,
			err error,
		)
	}

	// DefaultHandler struct
	DefaultHandler struct {
		jsonEncoder gonethttpjson.Encoder
	}
)

// NewDefaultHandler function
func NewDefaultHandler(
	jsonEncoder gonethttpjson.Encoder,
) (*DefaultHandler, error) {
	// Check if the JSON encoder is nil
	if jsonEncoder == nil {
		return nil, gonethttpjson.ErrNilJSONEncoder
	}

	return &DefaultHandler{
		jsonEncoder: jsonEncoder,
	}, nil
}

// HandlerError handles the error
func (d *DefaultHandler) HandlerError(
	w http.ResponseWriter,
	err error,
) {
	response := gonethttpresponse.ErrorResponse{Error: err.Error()}

	// Encode the response
	_ = d.jsonEncoder.Encode(w, &response, http.StatusInternalServerError)
}
