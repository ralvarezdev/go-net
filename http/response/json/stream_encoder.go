package json

import (
	"encoding/json"
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
		mode, logger,
	}
}

// Encode encodes the body into JSON and writes it to the response
//
// Parameters:
//
//   - w: The HTTP response writer
//   - response: The response to encode
//
// Returns:
//
//   - error: The error if any
func (s StreamEncoder) Encode(
	w http.ResponseWriter,
	response gonethttpresponse.Response,
) (err error) {
	// Get the body and HTTP status from the response
	body := response.Body(s.mode)
	httpStatus := response.HTTPStatus()

	// Set the Content-Type header if it hasn't been set already
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	// Write the HTTP status if it hasn't been written already
	if w.Header().Get("X-Status-Written") == "" {
		w.Header().Set("X-Status-Written", "true")
		w.WriteHeader(httpStatus)
	}

	// Encode the JSON body
	if err = json.NewEncoder(w).Encode(body); err != nil {
		// Overwrite the status on error
		w.Header().Set("X-Status-Written", "true")
		w.WriteHeader(http.StatusInternalServerError)

		if s.logger != nil {
			s.logger.Error(
				"failed to marshal response body",
				slog.String("error", err.Error()),
			)
		}
		http.Error(
			w,
			gonethttp.InternalServerError,
			http.StatusInternalServerError,
		)
		return err
	}
	return nil
}
