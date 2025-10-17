package response

import (
	"io"
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
)

type (
	// Response is the interface for the responses
	Response interface {
		Body(mode *goflagsmode.Flag) interface{}
		HTTPStatus() int
	}

	// Encoder interface
	Encoder interface {
		Encode(
			body interface{},
		) ([]byte, error)
		EncodeAndWrite(
			writer io.Writer,
			beforeWriteFn func() error,
			body interface{},
		) error
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
			body interface{},
		) (map[string]interface{}, error)
	}
)
