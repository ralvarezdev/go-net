package json

import (
	"encoding/json"
	"github.com/ralvarezdev/go-flags/mode"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpstatusresponse "github.com/ralvarezdev/go-net/http/status/response"
	"net/http"
)

var (
	ErrCodeMarshalResponseBodyFailed *string
)

type (
	// Encoder interface
	Encoder interface {
		Encode(
			w http.ResponseWriter,
			response gonethttpresponse.Response,
		) error
	}

	// DefaultEncoder struct
	DefaultEncoder struct {
		mode *mode.Flag
	}
)

// NewDefaultEncoder creates a new default JSON encoder
func NewDefaultEncoder(mode *mode.Flag) *DefaultEncoder {
	return &DefaultEncoder{mode}
}

// Encode encodes the body into JSON and writes it to the response
func (d *DefaultEncoder) Encode(
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
