package json

import (
	"encoding/json"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"net/http"
)

type (
	// StreamDecoder is the JSON decoder interface
	StreamDecoder interface {
		Decode(
			w http.ResponseWriter,
			r *http.Request,
			data interface{},
		) (err error)
	}

	// DefaultStreamDecoder is the JSON decoder struct
	DefaultStreamDecoder struct {
		mode          *goflagsmode.Flag
		streamEncoder StreamEncoder
	}
)

// NewDefaultStreamDecoder creates a new JSON decoder
func NewDefaultStreamDecoder(
	mode *goflagsmode.Flag,
	streamEncoder StreamEncoder,
) *DefaultStreamDecoder {
	return &DefaultStreamDecoder{
		mode:          mode,
		streamEncoder: streamEncoder,
	}
}

// Decode decodes the JSON data
func (d *DefaultStreamDecoder) Decode(
	w http.ResponseWriter,
	r *http.Request,
	data interface{},
) (err error) {
	// Check the data type
	if err = checkJSONData(w, data, d.mode); err != nil {
		return err
	}

	// Decode JSON data
	if err = json.NewDecoder(r.Body).Decode(data); err != nil {
		if d.streamEncoder != nil {
			_ = d.streamEncoder.Encode(
				w,
				gonethttpresponse.NewJSONErrorResponse(err),
				http.StatusInternalServerError,
			)
		} else {
			http.Error(
				w,
				err.Error(),
				http.StatusBadRequest,
			)
		}
	}
	return err
}
