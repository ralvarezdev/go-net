package json

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ralvarezdev/go-flags/mode"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	// Encoder struct
	Encoder struct {
		mode   *mode.Flag
		logger *slog.Logger
	}
)

// NewEncoder creates a new default JSON encoder
//
// Parameters:
//
//   - mode: The flag mode
//   - logger: The logger
//
// Returns:
//
//   - *Encoder: The default encoder
func NewEncoder(
	mode *mode.Flag,
	logger *slog.Logger,
) *Encoder {
	if logger != nil {
		logger = logger.With(
			slog.String("component", "http_response_json_encoder"),
		)
	}
	return &Encoder{mode, logger}
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
func (e Encoder) Encode(
	w http.ResponseWriter,
	response gonethttpresponse.Response,
) (err error) {
	// Get the response body and HTTP status
	body := response.Body(e.mode)
	httpStatus := response.HTTPStatus()

	// Encode the JSON body
	jsonBody, err := json.Marshal(body)
	if err != nil {
		if e.logger != nil {
			e.logger.Error(
				"failed to marshal response body",
				slog.String("error", err.Error()),
			)
		}
		http.Error(
			w,
			gonethttp.InternalServerError,
			http.StatusInternalServerError,
		)
	}

	// Set the Content-Type header if it hasn't been set already
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	// Write the HTTP status if it hasn't been written already
	if w.Header().Get("X-Status-Written") == "" {
		w.Header().Set("X-Status-Written", "true")
		w.WriteHeader(httpStatus)
	}

	// Write the JSON body to the response
	_, err = w.Write(jsonBody)
	return err
}
