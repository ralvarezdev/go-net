package protojson

import (
	"fmt"
	"io"
	"net/http"

	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpresponsehandler "github.com/ralvarezdev/go-net/http/response/handler"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type (
	Decoder struct {
		unmarshalOptions protojson.UnmarshalOptions
		handler          gonethttpresponsehandler.Handler
	}
)

// NewDecoder creates a new Decoder instance
//
// Parameters:
//
//   - handler: The HTTP response handler (optional, can be nil)
//   - logger: The logger (optional, can be nil)
//
// Returns:
//
//   - *Decoder: The decoder instance
//   - error: The error if any
func NewDecoder(
	handler gonethttpresponsehandler.Handler,
) (*Decoder, error) {
	// Check if the handler is nil
	if handler == nil {
		return nil, gonethttpresponsehandler.ErrNilHandler
	}

	// Initialize unmarshal options
	unmarshalOptions := protojson.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}

	return &Decoder{
		unmarshalOptions: unmarshalOptions,
		handler:          handler,
	}, nil
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
		d.handler.HandleDebugErrorResponseWithCode(
			w,
			ErrInvalidProtoMessage,
			gonethttp.ErrInternalServerError,
			ErrCodeInvalidProtoMessage,
			http.StatusInternalServerError,
		)
		return ErrInvalidProtoMessage
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		d.handler.HandleDebugErrorResponseWithCode(
			w,
			fmt.Errorf(ErrReadBodyFailed, err.Error()),
			gonethttp.ErrInternalServerError,
			ErrCodeReadBodyFailed,
			http.StatusInternalServerError,
		)
		return err
	}

	// Decode the request body into the proto message
	if err = protojson.Unmarshal(body, msg); err != nil {
		d.handler.HandleDebugErrorResponseWithCode(
			w,
			fmt.Errorf(ErrUnmarshalProtoJSONFailed, err.Error()),
			gonethttp.ErrInternalServerError,
			ErrCodeUnmarshalProtoJSONFailed,
			http.StatusInternalServerError,
		)
		return err
	}
	return nil
}
