package request

import (
	"io"
	"net/http"
)

type (
	// Decoder interface
	Decoder interface {
		Decode(
			body any,
			dest any,
		) error
		DecodeReader(
			reader io.Reader,
			dest any,
		) error
		DecodeRequest(
			r *http.Request,
			dest any,
		) error
	}
)
