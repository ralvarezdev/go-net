package validator

import (
	gonethttpresponse "github.com/ralvarezdev/go-net/http/handler"
	gonethttpjson "github.com/ralvarezdev/go-net/http/json"
	"net/http"
)

// ErrorHandler handles the possible JWT validation error
type ErrorHandler func(w http.ResponseWriter, err error)

// NewDefaultErrorHandler function
func NewDefaultErrorHandler(
	jsonEncoder gonethttpjson.Encoder,
) (ErrorHandler, error) {
	// Check if the JSON encoder is nil
	if jsonEncoder == nil {
		return nil, gonethttpjson.ErrNilJSONEncoder
	}

	return func(w http.ResponseWriter, err error) {
		response := gonethttpresponse.NewJSONErrorResponse(err)

		// Encode the response
		_ = jsonEncoder.Encode(w, &response, http.StatusInternalServerError)
	}, nil
}
