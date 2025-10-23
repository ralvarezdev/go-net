package protojson

import (
	"io"
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gojsonencoder "github.com/ralvarezdev/go-json/encoder"
	gojsonencoderprotojson "github.com/ralvarezdev/go-json/encoder/protojson"

	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpresponsejson "github.com/ralvarezdev/go-net/http/response/json"
)

type (
	// Encoder is the implementation of the Encoder interface
	Encoder struct {
		protoJSONEncoder gojsonencoder.ProtoJSONEncoder
		jsonEncoder      gonethttpresponse.Encoder
		mode             *goflagsmode.Flag
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
	// Initialize the ProtoJSON encoder
	protoJSONEncoder := gojsonencoderprotojson.NewEncoder()

	// Initialize the JSON encoder
	jsonEncoder := gonethttpresponsejson.NewEncoder(
		mode,
	)

	return &Encoder{
		mode:             mode,
		protoJSONEncoder: protoJSONEncoder,
		jsonEncoder:      jsonEncoder,
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
// - (map[string]any, error): The precomputed marshaled body and the error if any
func (e Encoder) PrecomputeMarshal(
	body any,
) (map[string]any, error) {
	precomputedMarshal, err := e.protoJSONEncoder.PrecomputeMarshal(body)
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
	body any,
) ([]byte, error) {
	// Marshal the instance to get the precomputed body
	precomputedMarshal, err := e.PrecomputeMarshal(body)
	if err != nil {
		return nil, err
	}

	// Marshal to JSON
	marshaledBody, err := e.protoJSONEncoder.Encode(precomputedMarshal)
	if err != nil {
		return nil, gonethttpresponse.NewDebugErrorWithCode(
			err,
			gonethttp.ErrInternalServerError,
			ErrCodeProtoJSONMarshalFailed,
			http.StatusInternalServerError,
		)
	}
	return marshaledBody, nil
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
	body any,
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
