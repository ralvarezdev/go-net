package json

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	// StreamEncoder is the JSON encoder struct
	StreamEncoder struct {
		mode   *goflagsmode.Flag
		logger *slog.Logger
	}
)

// NewStreamEncoder creates a new JSON encoder
//
// Parameters:
//
//   - mode: The flag mode
//   - logger: The logger
//
// Returns:
//
//   - *StreamEncoder: The default encoder
func NewStreamEncoder(
	mode *goflagsmode.Flag,
	logger *slog.Logger,
) *StreamEncoder {
	if logger != nil {
		logger = logger.With(
			slog.String("component", "http_response_json_stream_encoder"),
		)
	}

	return &StreamEncoder{
		mode,
		logger,
	}
}

// Encode encodes the body into JSON
//
// Parameters:
//
//   - body: The body to encode
//
// Returns:
//
//   - ([]byte): The encoded JSON
//   - error: The error if any
func (s StreamEncoder) Encode(
	body interface{},
) ([]byte, error) {
	// Create a buffer to write to
	buffer := new(bytes.Buffer)

	// Wrap it with a bufio.Writer
	writer := bufio.NewWriter(buffer)

	// Encode the body into JSON
	if err := json.NewEncoder(writer).Encode(body); err != nil {
		return nil, gonethttpresponse.NewDebugErrorWithCode(
			err,
			gonethttp.ErrInternalServerError,
			ErrCodeJSONMarshalFailed,
			http.StatusInternalServerError,
		)
	}

	// Flush to ensure all data is written to the underlying buffer
	if err := writer.Flush(); err != nil {
		return nil, gonethttpresponse.NewDebugErrorWithCode(
			err,
			gonethttp.ErrInternalServerError,
			ErrCodeJSONMarshalFailed,
			http.StatusInternalServerError,
		)
	}

	return buffer.Bytes(), nil
}

// EncodeAndWrite encodes the body into JSON and writes it to the writer
//
// Parameters:
//
//   - writer: The writer
//   - beforeWriteFn: The function to call before writing
//   - body: The body to encode
//
// Returns:
//
//   - error: The error if any
func (s StreamEncoder) EncodeAndWrite(
	writer io.Writer,
	beforeWriteFn func() error,
	body interface{},
) (err error) {
	// Call the before write function if provided
	if beforeWriteFn != nil {
		if err = beforeWriteFn(); err != nil {
			return err
		}
	}

	// Encode the body into JSON
	if err = json.NewEncoder(writer).Encode(body); err != nil {
		return gonethttpresponse.NewDebugErrorWithCode(
			err,
			gonethttp.ErrInternalServerError,
			ErrCodeJSONMarshalFailed,
			http.StatusInternalServerError,
		)
	}

	return nil
}

// EncodeResponse encodes the body into JSON
//
// Parameters:
//
//   - response: The response to encode
//
// Returns:
//
//   - ([]byte): The encoded JSON
//   - error: The error if any
func (s StreamEncoder) EncodeResponse(
	response gonethttpresponse.Response,
) ([]byte, error) {
	// Get the body from the response
	body := response.Body(s.mode)

	// Encode the JSON body
	encodedBody, err := s.Encode(body)
	if err != nil {
		return nil, err
	}
	return encodedBody, nil
}

// writeHeaders writes the headers to the response writer
//
// Parameters:
//
//   - writer: The HTTP response writer
//   - httpStatus: The HTTP status code
//
// Returns:
//
//   - error: The error if any
func (s StreamEncoder) writeHeaders(
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

// EncodeAndWriteResponse encodes the body into JSON and writes it to the response
//
// Parameters:
//
//   - writer: The HTTP response writer
//   - response: The response to encode
//
// Returns:
//
//   - error: The error if any
func (s StreamEncoder) EncodeAndWriteResponse(
	writer http.ResponseWriter,
	response gonethttpresponse.Response,
) error {
	// Get the body and HTTP status from the response
	body := response.Body(s.mode)

	// Create a before write function to write the headers
	beforeWriteFn := func() error {
		return s.writeHeaders(writer, response.HTTPStatus())
	}

	// Encode the JSON body
	if err := s.EncodeAndWrite(
		writer,
		beforeWriteFn,
		body,
	); err != nil {
		if s.logger != nil {
			s.logger.Error(
				"Failed to encode and write response",
				slog.String("error", err.Error()),
			)
		}
		return err
	}
	return nil
}
