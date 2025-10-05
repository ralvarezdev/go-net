package json

import (
	"net/http"

	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	// Decoder interface
	Decoder interface {
		Decode(
			w http.ResponseWriter,
			r *http.Request,
			dest interface{},
		) (err error)
	}

	// Encoder interface
	Encoder interface {
		Encode(
			w http.ResponseWriter,
			response gonethttpresponse.Response,
		) error
	}
)
