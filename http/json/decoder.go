package json

import (
	"encoding/json"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttpstatusresponse "github.com/ralvarezdev/go-net/http/status/response"
	"io"
	"net/http"
)

var (
	ErrCodeFailedToReadBody *string
	ErrCodeNilDestination   *string
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
		mode,
		encoder,
	}, nil
}

// Decode decodes the JSON request body and stores it in the destination
func (d *DefaultDecoder) Decode(
	w http.ResponseWriter,
	r *http.Request,
	dest interface{},
) (err error) {
	// Check the content type
	if !CheckContentType(r) {
		_ = d.encoder.Encode(
			w,
			ErrInvalidContentType.Response(),
		)
	}

	// Check the decoder destination
	if dest == nil {
		_ = d.encoder.Encode(
			w,
			gonethttpstatusresponse.NewJSendInternalServerError(ErrCodeNilDestination),
		)
	}

	// Get the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		_ = d.encoder.Encode(
			w,
			gonethttpstatusresponse.NewJSendDebugInternalServerError(
				err,
				ErrCodeFailedToReadBody,
			),
		)
		return err
	}

	// Decode JSON body into destination
	if err = json.Unmarshal(body, dest); err != nil {
		_ = BodyDecodeErrorHandler(w, err, d.encoder)
	}
	return err
}
