package validator

import (
	gonethttpjson "github.com/ralvarezdev/go-net/http/json"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
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
		return nil, gonethttpjson.ErrNilEncoder
	}

	return func(w http.ResponseWriter, err error) {
		_ = jsonEncoder.Encode(
			w, gonethttpresponse.NewErrorResponse(
				err,
				nil,
				nil,
				http.StatusInternalServerError,
			),
		)
	}, nil
}
