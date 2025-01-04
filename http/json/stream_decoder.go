package json

import (
	"encoding/json"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"net/http"
)

type (
	// DefaultStreamDecoder is the JSON decoder struct
	DefaultStreamDecoder struct {
		mode    *goflagsmode.Flag
		encoder Encoder
	}
)

// NewDefaultStreamDecoder creates a new JSON decoder
func NewDefaultStreamDecoder(
	mode *goflagsmode.Flag,
	encoder Encoder,
) (*DefaultStreamDecoder, error) {
	// Check if the stream encoder is nil
	if encoder == nil {
		return nil, ErrNilEncoder
	}

	return &DefaultStreamDecoder{
		mode:    mode,
		encoder: encoder,
	}, nil
}

// Decode decodes the JSON data
func (d *DefaultStreamDecoder) Decode(
	w http.ResponseWriter,
	r *http.Request,
	data interface{},
) (err error) {
	// Check the data type
	if err = checkJSONData(w, data, d.mode, d.encoder); err != nil {
		return err
	}

	// Decode JSON data
	if err = json.NewDecoder(r.Body).Decode(data); err != nil {
		_ = d.encoder.Encode(
			w,
			gonethttpresponse.NewJSONErrorResponse(err),
			http.StatusInternalServerError,
		)
	}
	return err
}
