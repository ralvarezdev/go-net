package json

import (
	"io"
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gojsondecoder "github.com/ralvarezdev/go-json/decoder"

	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttprequest "github.com/ralvarezdev/go-net/http/request"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	// Decoder struct
	Decoder struct {
		decoder gojsondecoder.Decoder
		mode    *goflagsmode.Flag
	}
)

// NewDecoder creates a new JSON decoder
//
// Parameters:
//
//   - mode: The flag mode
//   - decoder: The JSON decoder
//
// Returns:
//
//   - *Decoder: The default decoder
//   - error: The error if any
func NewDecoder(
	mode *goflagsmode.Flag,
	decoder gojsondecoder.Decoder,
) (*Decoder, error) {
	// Create the JSON decoder
	if decoder == nil {
		return nil, gonethttpresponse.NewDebugErrorWithCode(
			gonethttprequest.ErrNilDecoder,
			gonethttp.ErrInternalServerError,
			gonethttprequest.ErrCodeNilDecoder,
			http.StatusInternalServerError,
		)
	}

	return &Decoder{
		decoder: decoder,
		mode:    mode,
	}, nil
}

// Decode decodes the JSON body from an any value and stores it in the destination
//
// Parameters:
//
//   - body: The body to decode
//   - dest: The destination to store the decoded body
//
// Returns:
//
//   - error: The error if any
func (d Decoder) Decode(
	body any,
	dest any,
) error {
	if err := d.decoder.Decode(
		body,
		dest,
	); err != nil {
		return gonethttpresponse.NewDebugErrorWithCode(
			gonethttprequest.ErrInvalidBodyType,
			gonethttp.ErrInternalServerError,
			gonethttprequest.ErrCodeInvalidBodyType,
			http.StatusInternalServerError,
		)
	}
	return nil
}

// DecodeReader decodes the JSON body and stores it in the destination
//
// Parameters:
//
//   - reader: The reader to read the body from
//   - dest: The destination to store the decoded body
//
// Returns:
//
//   - error: The error if any
func (d Decoder) DecodeReader(
	reader io.Reader,
	dest any,
) error {
	if err := d.decoder.DecodeReader(
		reader,
		dest,
	); err != nil {
		return gonethttprequest.BodyDecodeErrorHandler(err)
	}
	return nil
}

// DecodeRequest decodes the JSON request body and stores it in the destination
//
// Parameters:
//
//   - request: The HTTP request
//   - dest: The destination to store the decoded body
//
// Returns:
//
//   - error: The error if any
func (d Decoder) DecodeRequest(
	r *http.Request,
	dest any,
) error {
	// Check the request
	if r == nil {
		return gonethttpresponse.NewDebugErrorWithCode(
			gonethttprequest.ErrNilRequest,
			gonethttp.ErrInternalServerError,
			gonethttprequest.ErrCodeNilRequest,
			http.StatusInternalServerError,
		)
	}

	// Check the content type
	if !gonethttprequest.CheckContentType(r) {
		return gonethttpresponse.NewFailFieldErrorWithCode(
			gonethttprequest.ErrInvalidContentTypeField,
			gonethttprequest.ErrInvalidContentType,
			gonethttprequest.ErrCodeInvalidContentType,
			http.StatusUnsupportedMediaType,
		)
	}
	return d.DecodeReader(r.Body, dest)
}
