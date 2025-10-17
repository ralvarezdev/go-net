package request

import (
	"io"
	"net/http"
)

type (
	// Decoder interface
	Decoder interface {
		Decode(
			body interface{},
			dest interface{},
		) error
		DecodeReader(
			reader io.Reader,
			dest interface{},
		) error
		DecodeRequest(
			r *http.Request,
			dest interface{},
		) error
	}
)
