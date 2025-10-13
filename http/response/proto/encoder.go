package proto

import (
	"net/http"
	"reflect"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpresponsejson "github.com/ralvarezdev/go-net/http/response/json"
	"google.golang.org/protobuf/encoding/protojson"
)

type (
	// Encoder is the implementation of the Encoder interface
	Encoder struct {
		jsonEncoder    *gonethttpresponsejson.Encoder
		mode           *goflagsmode.Flag
		marshalOptions protojson.MarshalOptions
	}
)

// NewEncoder creates a new Encoder instance
//
// Parameters:
//
//   - mode: the flag mode
//   - marshalOptions: options for marshaling protobuf messages to JSON
//
// Returns:
//
// - *Encoder: the new Encoder instance
func NewEncoder(
	mode *goflagsmode.Flag,
	marshalOptions protojson.MarshalOptions,
) *Encoder {
	// Initialize the JSON encoder
	jsonEncoder := gonethttpresponsejson.NewEncoder(mode)

	return &Encoder{
		mode:           mode,
		jsonEncoder:    jsonEncoder,
		marshalOptions: marshalOptions,
	}
}

// Encode protobuf to JSON
//
// Parameters:
//
//   - w: The HTTP response writer
//   - response: The response to encode, must have a data field that is a proto.Message
//
// Returns:
//
//   - error: The error if any
func (e Encoder) Encode(
	w http.ResponseWriter,
	response gonethttpresponse.Response,
) error {
	// Reflect on the response body to get its fields
	v := reflect.ValueOf(response.Body(e.mode))

	// Precompute the marshaled body
	precomputedBody, err := PrecomputeMarshalByReflection(v, &e.marshalOptions)
	if err != nil {
		return e.Encode(
			w,
			gonethttpresponse.NewJSendDebugInternalServerError(
				err,
				ErrCodeMarshalResponseBodyFailed,
			),
		)
	}

	// Create a new response with the precomputed body
	precomputedResponse := gonethttpresponse.NewResponse(
		precomputedBody,
		response.HTTPStatus(),
	)

	// Marshal to JSON
	return e.jsonEncoder.Encode(w, precomputedResponse)
}
