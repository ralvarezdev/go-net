package json

import (
	"encoding/json"
	"io"
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttpstatusresponse "github.com/ralvarezdev/go-net/http/status/response"
)

type (
	// DefaultDecoder struct
	DefaultDecoder struct {
		mode    *goflagsmode.Flag
		encoder Encoder
	}
)

// NewDefaultDecoder creates a new JSON decoder
//
// Parameters:
//
//   - mode: The flag mode
//   - encoder: The JSON encoder
//
// Returns:
//
//   - *DefaultDecoder: The default decoder
//   - error: The error if any
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
//
// Parameters:
//
//   - w: The HTTP response writer
//   - r: The HTTP request
//   - dest: The destination to store the decoded body
//
// Returns:
//
//   - error: The error if any
func (d DefaultDecoder) Decode(
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
