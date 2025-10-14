package response

import (
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
)

type (
	// Response is the interface for the responses
	Response interface {
		Body(mode *goflagsmode.Flag) interface{}
		HTTPStatus() int
	}

	// FailResponseError is the error for fail responses
	FailResponseError interface {
		Error() string
		Response() Response
	}

	// Encoder interface
	Encoder interface {
		Encode(
			w http.ResponseWriter,
			response Response,
		) error
	}
)
