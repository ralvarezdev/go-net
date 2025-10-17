package request

import (
	"io"
	"net/http"
)

type (
	// Decoder interface
	Decoder interface {
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
