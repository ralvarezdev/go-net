package protojson

import (
	"io"
	"net/http"

	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttprequest "github.com/ralvarezdev/go-net/http/request"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"google.golang.org/protobuf/encoding/protojson"
)

type (
	Decoder struct {
		unmarshalOptions protojson.UnmarshalOptions
	}
)

// NewDecoder creates a new Decoder instance
//
// Returns:
//
//   - *Decoder: The decoder instance
func NewDecoder() *Decoder {
	// Initialize unmarshal options
	unmarshalOptions := protojson.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}

	return &Decoder{
		unmarshalOptions: unmarshalOptions,
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
	return d.DecodeReader(reader, dest)
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
	return UnmarshalByReflection(
		reader,
		dest,
		&d.unmarshalOptions,
	)
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
