package json

import (
	"encoding/json"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"io"
	"net/http"
)

type (
	// Decoder interface
	Decoder interface {
		Decode(
			w http.ResponseWriter,
			r *http.Request,
			data interface{},
		) (err error)
	}

	// DefaultDecoder struct
	DefaultDecoder struct {
		mode    *goflagsmode.Flag
		encoder Encoder
	}
)

// NewDefaultDecoder creates a new JSON decoder
func NewDefaultDecoder(
	mode *goflagsmode.Flag,
	encoder Encoder,
) (*DefaultDecoder, error) {
	// Check if the encoder is nil
	if encoder == nil {
		return nil, ErrNilEncoder
	}

	return &DefaultDecoder{
		mode:    mode,
		encoder: encoder,
	}, nil
}

// Decode decodes the JSON data from the request
func (d *DefaultDecoder) Decode(
	w http.ResponseWriter,
	r *http.Request,
	data interface{},
) (err error) {
	// Check the data type
	if err = checkJSONData(w, data, d.mode, d.encoder); err != nil {
		return err
	}

	// Get the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		_ = d.encoder.Encode(
			w,
			gonethttpresponse.NewErrorResponse(
				err,
				nil,
				nil,
				http.StatusBadRequest,
			),
		)
		return err
	}

	// Decode JSON data
	if err = json.Unmarshal(body, data); err != nil {
		_ = d.encoder.Encode(
			w,
			gonethttpresponse.NewErrorResponse(
				err,
				nil,
				nil,
				http.StatusBadRequest,
			),
		)
	}
	return err
}
