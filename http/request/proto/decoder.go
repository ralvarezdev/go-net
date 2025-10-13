package proto

import (
	"fmt"
	"io"
	"net/http"

	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type (
	Decoder struct {
		unmarshalOptions protojson.UnmarshalOptions
		encoder          gonethttpresponse.Encoder
	}
)

// NewDecoder creates a new Decoder instance
//
// Parameters:
//
//   - encoder: The HTTP response encoder
//
// Returns:
//
//   - *Decoder: The decoder instance
//   - error: The error if any
func NewDecoder(
	encoder gonethttpresponse.Encoder,
) (*Decoder, error) {
	// Check if the encoder is nil
	if encoder == nil {
		return nil, gonethttpresponse.ErrNilEncoder
	}

	// Initialize unmarshal options
	unmarshalOptions := protojson.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}

	return &Decoder{
		unmarshalOptions: unmarshalOptions,
		encoder:          encoder,
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
		_ = d.encoder.Encode(
			w,
			gonethttpresponse.NewJSendDebugInternalServerError(
				ErrInvalidProtoMessage,
				ErrCodeInvalidProtoMessage,
			),
		)
		return ErrInvalidProtoMessage
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		_ = d.encoder.Encode(
			w,
			gonethttpresponse.NewJSendDebugInternalServerError(
				fmt.Errorf(ErrReadBodyFailed, err.Error()),
				ErrCodeReadBodyFailed,
			),
		)
		return err
	}

	// Decode the request body into the proto message
	if err = protojson.Unmarshal(body, msg); err != nil {
		_ = d.encoder.Encode(
			w,
			gonethttpresponse.NewJSendDebugBadRequest(
				fmt.Errorf(ErrUnmarshalProtoJSONFailed, err.Error()),
				ErrCodeUnmarshalProtoJSONFailed,
			),
		)
		return err
	}
	return nil
}
