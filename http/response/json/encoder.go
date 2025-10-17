package json

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ralvarezdev/go-flags/mode"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	// Encoder struct
	Encoder struct {
		mode *mode.Flag
	}
)

// NewEncoder creates a new default JSON encoder
//
// Parameters:
//
//   - mode: The flag mode
//
// Returns:
//
//   - *Encoder: The default encoder
func NewEncoder(
	mode *mode.Flag,
) *Encoder {
	return &Encoder{mode}
}

// Encode encodes the body into JSON bytes
//
// Parameters:
//
//   - body: The body to encode
//
// Returns:
//
//   - []byte: The encoded JSON bytes
//   - error: The error if any
func (e Encoder) Encode(
	body interface{},
) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, gonethttpresponse.NewDebugErrorWithCode(
			err,
			gonethttp.ErrInternalServerError,
			ErrCodeJSONMarshalFailed,
			http.StatusInternalServerError,
		)
	}
	return jsonBody, nil
}

// EncodeAndWrite encodes the body and writes it to the writer
//
// Parameters:
//
// - writer: The writer to write the response to
// - beforeWriteFn: The function to call before writing the response
// - body: The body to encode
//
// Returns:
//
// - error: The error if any
func (e Encoder) EncodeAndWrite(
	writer io.Writer,
	beforeWriteFn func() error,
	body interface{},
) error {
	// Encode the body into JSON
	jsonBody, err := e.Encode(body)
	if err != nil {
		return gonethttpresponse.NewDebugErrorWithCode(
			err,
			gonethttp.ErrInternalServerError,
			ErrCodeJSONMarshalFailed,
			http.StatusInternalServerError,
		)
	}

	// Call the before write function if provided
	if beforeWriteFn != nil {
		if err = beforeWriteFn(); err != nil {
			return err
		}
	}

	// Write the JSON body to the writer
	_, err = writer.Write(jsonBody)
	return err
}

// EncodeResponse encodes the response into JSON bytes
//
// Parameters:
//
//   - response: The response to encode
//
// Returns:
//
//   - []byte: The encoded JSON bytes
//   - error: The error if any
func (e Encoder) EncodeResponse(
	response gonethttpresponse.Response,
) ([]byte, error) {
	// Get the response body and HTTP status
	body := response.Body(e.mode)

	// Encode the body into JSON
	jsonBody, err := e.Encode(body)

	return jsonBody, err
}

// writeHeaders writes the headers to the http.ResponseWriter
//
// Parameters:
//
//   - ww: The http.ResponseWriter
//   - httpStatus: The HTTP status to write
//
// Returns:
//
//   - error: The error if any
func (e Encoder) writeHeaders(
	writer http.ResponseWriter,
	httpStatus int,
) error {
	// Set the Content-Type header if it hasn't been set already
	if writer.Header().Get("Content-Type") == "" {
		writer.Header().Set("Content-Type", "application/json")
	}

	// Write the HTTP status if it hasn't been written already
	if writer.Header().Get("X-Status-Written") == "" {
		writer.Header().Set("X-Status-Written", "true")
		writer.WriteHeader(httpStatus)
	}
	return nil
}

// EncodeAndWriteResponse encodes the response and writes it to the http.ResponseWriter
//
// Parameters:
//
//   - writer: The http.ResponseWriter
//   - response: The response to encode and write
//
// Returns:
//
//   - error: The error if any
func (e Encoder) EncodeAndWriteResponse(
	writer http.ResponseWriter,
	response gonethttpresponse.Response,
) error {
	// Get the response body and HTTP status
	body := response.Body(e.mode)

	// Build the before write function
	beforeWriteFn := func() error {
		return e.writeHeaders(writer, response.HTTPStatus())
	}

	return e.EncodeAndWrite(
		writer,
		beforeWriteFn,
		body,
	)
}
