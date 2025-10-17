package protojson

import (
	"io"
	"net/http"
	"reflect"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttp "github.com/ralvarezdev/go-net/http"
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
//   - logger: the logger (optional, can be nil)
//
// Returns:
//
// - *Encoder: the new Encoder instance
func NewEncoder(
	mode *goflagsmode.Flag,
) *Encoder {
	// Initialize the JSON encoder
	jsonEncoder := gonethttpresponsejson.NewEncoder(mode)

	// Initialize unmarshal options
	marshalOptions := protojson.MarshalOptions{
		AllowPartial: true,
	}

	return &Encoder{
		mode:           mode,
		jsonEncoder:    jsonEncoder,
		marshalOptions: marshalOptions,
	}
}

// PrecomputeMarshal precomputes the marshaled body by reflecting on the instance
//
// Parameters:
//
// - body: The body to precompute the marshaled body for
//
// Returns:
//
// - (map[string]interface{}, error): The precomputed marshaled body and the error if any
func (e Encoder) PrecomputeMarshal(
	body interface{},
) (map[string]interface{}, error) {
	// Reflect on the instance to get its fields
	v := reflect.ValueOf(body)

	// Precompute the marshaled body
	precomputedMarshal, err := PrecomputeMarshalByReflection(
		v,
		&e.marshalOptions,
	)
	if err != nil {
		return nil, gonethttpresponse.NewDebugErrorWithCode(
			err,
			gonethttp.ErrInternalServerError,
			ErrCodeProtoJSONMarshalFailed,
			http.StatusInternalServerError,
		)
	}
	return precomputedMarshal, nil
}

// Encode encodes the given body to JSON
//
// Parameters:
//
//   - body: The body to encode
//
// Returns:
//
//   - ([]byte, error): The encoded body and the error if any
func (e Encoder) Encode(
	body interface{},
) ([]byte, error) {
	// Marshal the instance to get the precomputed body
	precomputedMarshal, err := e.PrecomputeMarshal(body)
	if err != nil {
		return nil, err
	}
	return e.jsonEncoder.Encode(precomputedMarshal)
}

// EncodeResponse encodes the given response to JSON
//
// Parameters:
//
//   - response: The response to encode
//
// Returns:
//
//   - ([]byte, error): The encoded response and the error if any
func (e Encoder) EncodeResponse(
	response gonethttpresponse.Response,
) ([]byte, error) {
	// Get the response body
	body := response.Body(e.mode)

	// Marshal the instance to get the precomputed body
	precomputedMarshal, err := e.PrecomputeMarshal(body)
	if err != nil {
		return nil, err
	}

	// Create a new response with the precomputed body
	precomputedResponse := gonethttpresponse.NewResponse(
		precomputedMarshal,
		response.HTTPStatus(),
	)

	// Marshal to JSON
	return e.jsonEncoder.EncodeResponse(precomputedResponse)
}

// EncodeAndWrite encodes and writes the given body to the writer
//
// Parameters:
//
//   - writer: The writer to write the encoded body to
//   - beforeWriteFn: The function to call before writing the body
//   - body: The body to encode
//
// Returns:
//
// - error: The error if any
func (e Encoder) EncodeAndWrite(
	writer io.Writer,
	beforeWriteFn func() error,
	body interface{},
) error {
	// Marshal the instance to get the precomputed body
	precomputedMarshal, err := e.PrecomputeMarshal(body)
	if err != nil {
		return err
	}
	return e.jsonEncoder.EncodeAndWrite(
		writer,
		beforeWriteFn,
		precomputedMarshal,
	)
}

// EncodeAndWriteResponse encodes and writes the response to the HTTP response writer
//
// Parameters:
//
//   - writer: The HTTP response writer
//   - response: The response to encode, must have a data field that is a proto.Message
//
// Returns:
//
//   - error: The error if any
func (e Encoder) EncodeAndWriteResponse(
	writer http.ResponseWriter,
	response gonethttpresponse.Response,
) error {
	// Get the response body
	body := response.Body(e.mode)

	// Marshal the instance to get the precomputed body
	precomputedMarshal, err := e.PrecomputeMarshal(body)
	if err != nil {
		return err
	}

	// Create a new response with the precomputed body
	precomputedResponse := gonethttpresponse.NewResponse(
		precomputedMarshal,
		response.HTTPStatus(),
	)

	// Marshal to JSON
	return e.jsonEncoder.EncodeAndWriteResponse(writer, precomputedResponse)
}
