package protojson

import (
	"io"
	"net/http"

	gojsondecoder "github.com/ralvarezdev/go-json/decoder"
	gojsondecoderprotojson "github.com/ralvarezdev/go-json/decoder/protojson"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttprequest "github.com/ralvarezdev/go-net/http/request"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	Decoder struct {
		decoder gojsondecoder.Decoder
	}
)

// NewDecoder creates a new Decoder instance
//
// Returns:
//
//   - *Decoder: The decoder instance
func NewDecoder() *Decoder {
	// Create the JSON decoder
	decoder := gojsondecoderprotojson.NewDecoder()

	return &Decoder{
		decoder: decoder,
	}
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
	body interface{},
	dest interface{},
) error {
	if err := d.decoder.Decode(
		body,
		dest,
	); err != nil {
		return gonethttprequest.BodyDecodeErrorHandler(err)
	}
	return nil
}

// DecodeReader  decodes a JSON body from a reader into a destination
//
// Parameters:
//
//   - reader: The io.Reader to read the body from
//   - dest: The destination to decode the body into
//
// Returns:
//
//   - error: The error if any
func (d Decoder) DecodeReader(
	reader io.Reader,
	dest interface{},
) error {
	if err := d.decoder.DecodeReader(
		reader,
		dest,
	); err != nil {
		return gonethttprequest.BodyDecodeErrorHandler(err)
	}
	return nil
}

// DecodeRequest decodes a JSON body from an HTTP request into a destination
//
// Parameters:
//
//   - request: The HTTP request to read the body from
//   - dest: The destination to decode the body into
//
// Returns:
//
//   - error: The error if any
func (d Decoder) DecodeRequest(
	r *http.Request,
	dest interface{},
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
