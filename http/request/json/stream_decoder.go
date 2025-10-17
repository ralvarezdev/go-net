package json

import (
	"encoding/json"
	"io"
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttprequest "github.com/ralvarezdev/go-net/http/request"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	// StreamDecoder is the JSON decoder struct
	StreamDecoder struct {
		mode *goflagsmode.Flag
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
	return &StreamDecoder{
		mode,
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
	// Check the body type
	reader, err := gonethttprequest.ToReader(body)
	if err != nil {
		return gonethttpresponse.NewDebugErrorWithCode(
			gonethttprequest.ErrInvalidBodyType,
			gonethttp.ErrInternalServerError,
			gonethttprequest.ErrCodeInvalidBodyType,
			http.StatusInternalServerError,
		)
	}
	return s.DecodeReader(reader, dest)
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
	// Check the decoder destination
	if dest == nil {
		return gonethttpresponse.NewDebugErrorWithCode(
			ErrNilDestination,
			gonethttp.ErrInternalServerError,
			ErrCodeNilDestination,
			http.StatusInternalServerError,
		)
	}

	// Create the stream decoder
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	// Decode JSON body into destination
	if err := decoder.Decode(dest); err != nil {
		return BodyDecodeErrorHandler(err)
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
