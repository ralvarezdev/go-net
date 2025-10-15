package protojson

import (
	"fmt"
	"io"
	"net/http"

	gonethttp "github.com/ralvarezdev/go-net/http"
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

// Decode decodes the request body into the destination
//
// Parameters:
//
//   - w: The HTTP response writer
//   - r: The HTTP request
//   - dest: The destination to decode the request body into
//
// Returns:
//
//   - error: The error if any
func (d Decoder) Decode(
	w http.ResponseWriter,
	r *http.Request,
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
	body, err := io.ReadAll(r.Body)
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
