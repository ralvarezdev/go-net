package validator

import (
	gonethttpjson "github.com/ralvarezdev/go-net/http/json"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gostringsconvert "github.com/ralvarezdev/go-strings/convert"
	"net/http"
)

// FailHandler handles the possible JWT validation errors
type FailHandler func(w http.ResponseWriter, err ...error)

// NewDefaultFailHandler function
func NewDefaultFailHandler(
	jsonEncoder gonethttpjson.Encoder,
) (FailHandler, error) {
	// Check if the JSON encoder is nil
	if jsonEncoder == nil {
		return nil, gonethttpjson.ErrNilEncoder
	}

	return func(w http.ResponseWriter, err ...error) {
		// Create the body map
		var body = make(map[string]*[]string)
		body["authorization"] = gostringsconvert.ErrorArrayToStringArray(&err)

		// Encode the response
		_ = jsonEncoder.Encode(
			w, gonethttpresponse.NewFailResponse(
				&body,
				nil,
				http.StatusUnauthorized,
			),
		)
	}, nil
}
