package json

import (
	"encoding/json"
	"net/http"

	"github.com/ralvarezdev/go-flags/mode"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpstatusresponse "github.com/ralvarezdev/go-net/http/status/response"
)

type (
	// DefaultEncoder struct
	DefaultEncoder struct {
		mode *mode.Flag
	}
)

// NewDefaultEncoder creates a new default JSON encoder
//
// Parameters:
//
//   - mode: The flag mode
//
// Returns:
//
//   - *DefaultEncoder: The default encoder
func NewDefaultEncoder(mode *mode.Flag) *DefaultEncoder {
	return &DefaultEncoder{mode}
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
func (d DefaultEncoder) Encode(
	w http.ResponseWriter,
	response gonethttpresponse.Response,
) (err error) {
	// Get the response body and HTTP status
	body := response.Body(d.mode)
	httpStatus := response.HTTPStatus()

	// Encode the JSON body
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return d.Encode(
			w,
			gonethttpstatusresponse.NewJSendDebugInternalServerError(
				err,
				ErrCodeMarshalResponseBodyFailed,
			),
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
