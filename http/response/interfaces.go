package response

import (
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gojsonecoder "github.com/ralvarezdev/go-json/encoder"
)

type (
	// Response is the interface for the responses
	Response interface {
		Body(mode *goflagsmode.Flag) any
		HTTPStatus() int
	}

	// Encoder interface
	Encoder interface {
		gojsonecoder.Encoder
		EncodeResponse(
			response Response,
		) ([]byte, error)
		EncodeAndWriteResponse(
			w http.ResponseWriter,
			response Response,
		) error
	}

	// ProtoJSONEncoder interface
	ProtoJSONEncoder interface {
		PrecomputeMarshal(
			body any,
		) (map[string]any, error)
	}
)
