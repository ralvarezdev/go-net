package protojson

import (
	"fmt"
	"io"
	"net/http"

	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttprequest "github.com/ralvarezdev/go-net/http/request"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
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
	// Assert that dest is a proto.Message
	msg, ok := dest.(proto.Message)
	if !ok {
		return gonethttpresponse.NewDebugErrorWithCode(
			ErrInvalidProtoMessage,
			gonethttp.ErrInternalServerError,
			ErrCodeInvalidProtoMessage,
			http.StatusInternalServerError,
		)
	}

	// Read the request body
	body, err := io.ReadAll(reader)
	if err != nil {
		return gonethttpresponse.NewDebugErrorWithCode(
			fmt.Errorf(ErrReadBodyFailed, err),
			gonethttp.ErrInternalServerError,
			ErrCodeReadBodyFailed,
			http.StatusInternalServerError,
		)
	}

	// Decode the request body into the proto message
	if err = protojson.Unmarshal(body, msg); err != nil {
		return gonethttpresponse.NewDebugErrorWithCode(
			fmt.Errorf(ErrUnmarshalProtoJSONFailed, err),
			gonethttp.ErrInternalServerError,
			ErrCodeUnmarshalProtoJSONFailed,
			http.StatusInternalServerError,
		)
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
	request *http.Request,
	dest interface{},
) error {
	// Check the request
	if request == nil {
		return gonethttpresponse.NewDebugErrorWithCode(
			gonethttprequest.ErrNilRequest,
			gonethttp.ErrInternalServerError,
			gonethttprequest.ErrCodeNilRequest,
			http.StatusInternalServerError,
		)
	}

	if !gonethttprequest.CheckContentType(request) {
		return gonethttpresponse.NewFailFieldErrorWithCode(
			gonethttprequest.ErrInvalidContentTypeField,
			gonethttprequest.ErrInvalidContentType,
			gonethttprequest.ErrCodeInvalidContentType,
			http.StatusUnsupportedMediaType,
		)
	}

	return d.DecodeReader(request.Body, dest)
}
