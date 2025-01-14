package validator

import (
	gojwtnethttp "github.com/ralvarezdev/go-jwt/net/http"
	gonethttpjson "github.com/ralvarezdev/go-net/http/json"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"net/http"
)

// FailHandler handles the possible JWT validation errors
type FailHandler func(
	w http.ResponseWriter,
	error error,
)

// NewDefaultFailHandler function
func NewDefaultFailHandler(
	jsonEncoder gonethttpjson.Encoder,
) (FailHandler, error) {
	// Check if the JSON encoder is nil
	if jsonEncoder == nil {
		return nil, gonethttpjson.ErrNilEncoder
	}

	return func(
		w http.ResponseWriter,
		error error,
	) {
		// Encode the response
		_ = jsonEncoder.Encode(
			w, gonethttpresponse.NewFailResponseFromRequestError(
				gonethttpresponse.NewHeaderError(
					gojwtnethttp.AuthorizationHeaderKey,
					error.Error(),
					http.StatusUnauthorized,
				),
			),
		)
	}, nil
}
