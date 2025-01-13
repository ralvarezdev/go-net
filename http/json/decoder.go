package json

import (
	"encoding/json"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpstatuserrors "github.com/ralvarezdev/go-net/http/status/errors"
	"io"
	"net/http"
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

// Decode decodes the JSON request body and stores it in the destination
func (d *DefaultDecoder) Decode(
	w http.ResponseWriter,
	r *http.Request,
	dest interface{},
) (err error) {
	// Check the decoder destination
	if dest == nil {
		_ = d.encoder.Encode(
			w, gonethttpresponse.NewDebugErrorResponse(
				gonethttpstatuserrors.InternalServerError,
				err,
				nil,
				nil,
				http.StatusInternalServerError,
			),
		)
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

	// Decode JSON body into destination
	if err = json.Unmarshal(body, dest); err != nil {
		_ = bodyDecodeErrorHandler(w, err, d.encoder)
	}
	return err
}
