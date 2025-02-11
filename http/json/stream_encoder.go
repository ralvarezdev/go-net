package json

import (
	"encoding/json"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpstatusresponse "github.com/ralvarezdev/go-net/http/status/response"
	"net/http"
)

type (
	// DefaultStreamEncoder is the JSON encoder struct
	DefaultStreamEncoder struct {
		mode *goflagsmode.Flag
	}
)

// NewDefaultStreamEncoder creates a new JSON encoder
func NewDefaultStreamEncoder(mode *goflagsmode.Flag) *DefaultStreamEncoder {
	return &DefaultStreamEncoder{
		mode,
	}
}

// Encode encodes the body into JSON and writes it to the response
func (d *DefaultStreamEncoder) Encode(
	w http.ResponseWriter,
	response gonethttpresponse.Response,
) (err error) {
	// Get the body and HTTP status from the response
	body := response.Body(d.mode)
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

		_ = d.Encode(
			w,
			gonethttpstatusresponse.NewJSendDebugInternalServerError(
				err,
				ErrCodeMarshalResponseBodyFailed,
			),
		)
		return err
	}
	return nil
}
