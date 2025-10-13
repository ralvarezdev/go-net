package request

import (
	"net/http"
)

type (
	// Decoder interface
	Decoder interface {
		Decode(
			w http.ResponseWriter,
			r *http.Request,
			dest interface{},
		) error
	}
)
