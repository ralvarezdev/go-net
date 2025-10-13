package json

import (
	"encoding/json"
	"io"
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	// Decoder struct
	Decoder struct {
		mode    *goflagsmode.Flag
		encoder gonethttpresponse.Encoder
	}
)

// NewDecoder creates a new JSON decoder
//
// Parameters:
//
//   - mode: The flag mode
//   - encoder: The JSON encoder
//
// Returns:
//
//   - *Decoder: The default decoder
//   - error: The error if any
func NewDecoder(
	mode *goflagsmode.Flag,
	encoder gonethttpresponse.Encoder,
) (*Decoder, error) {
	// Check if the encoder is nil
	if encoder == nil {
		return nil, gonethttpresponse.ErrNilEncoder
	}

	return &Decoder{
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
func (d Decoder) Decode(
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
			gonethttpresponse.NewJSendDebugInternalServerError(
				ErrNilDestination,
				ErrCodeNilDestination,
			),
		)
	}

	// Get the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		_ = d.encoder.Encode(
			w,
			gonethttpresponse.NewJSendDebugInternalServerError(
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
