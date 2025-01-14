package json

import (
	"encoding/json"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttpstatusresponse "github.com/ralvarezdev/go-net/http/status/response"
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

// Decode decodes the JSON request body and stores it in the destination
func (d *DefaultStreamDecoder) Decode(
	w http.ResponseWriter,
	r *http.Request,
	dest interface{},
) (err error) {
	// Check the decoder destination
	if dest == nil {
		_ = d.encoder.Encode(
			w,
			gonethttpstatusresponse.NewJSendDebugInternalServerError(
				err,
				ErrCodeNilDestination,
			),
		)
	}

	// Decode JSON body into destination
	if err = json.NewDecoder(r.Body).Decode(dest); err != nil {
		_ = bodyDecodeErrorHandler(w, err, d.encoder)
	}
	return err
}
