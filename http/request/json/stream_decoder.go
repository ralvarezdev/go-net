package json

import (
	"io"
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gojsondecoder "github.com/ralvarezdev/go-json/decoder"
	gojsondecoderjson "github.com/ralvarezdev/go-json/decoder/json"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttprequest "github.com/ralvarezdev/go-net/http/request"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	// StreamDecoder is the JSON decoder struct
	StreamDecoder struct {
		decoder gojsondecoder.Decoder
		mode    *goflagsmode.Flag
	}
)

// NewStreamDecoder creates a new JSON decoder
//
// Parameters:
//
//   - mode: The flag mode
//   - responsesHandler: The HTTP handler to handle errors
//
// Returns:
//
//   - *StreamDecoder: The default decoder
func NewStreamDecoder(
	mode *goflagsmode.Flag,
) *StreamDecoder {
	// Create the JSON decoder
	decoder := gojsondecoderjson.NewDecoder()

	return &StreamDecoder{
		decoder: decoder,
		mode:    mode,
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
func (s StreamDecoder) Decode(
	body interface{},
	dest interface{},
) error {
	if err := s.decoder.Decode(body, dest); err != nil {
		return gonethttpresponse.NewDebugErrorWithCode(
			gonethttprequest.ErrInvalidBodyType,
			gonethttp.ErrInternalServerError,
			gonethttprequest.ErrCodeInvalidBodyType,
			http.StatusInternalServerError,
		)
	}
	return nil
}

// DecodeReader decodes a JSON body from a reader into a destination
//
// Parameters:
//
//   - reader: The reader to read the body from
//   - dest: The destination to store the decoded body
//
// Returns:
//
//   - error: The error if any
func (s StreamDecoder) DecodeReader(
	reader io.Reader,
	dest interface{},
) error {
	if err := s.decoder.DecodeReader(
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
//   - dest: The destination to store the decoded body
//
// Returns:
//
//   - error: The error if any
func (s StreamDecoder) DecodeRequest(
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

	// Check the content type
	if !gonethttprequest.CheckContentType(r) {
		return gonethttpresponse.NewFailFieldErrorWithCode(
			gonethttprequest.ErrInvalidContentTypeField,
			gonethttprequest.ErrInvalidContentType,
			gonethttprequest.ErrCodeInvalidContentType,
			http.StatusUnsupportedMediaType,
		)
	}

	// Decode the body from the request
	return s.DecodeReader(
		r.Body,
		dest,
	)
}
