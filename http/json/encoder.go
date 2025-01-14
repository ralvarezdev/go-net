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
			response *gonethttpresponse.Response,
		) error
	}

	// DefaultEncoder struct
	DefaultEncoder struct {
		mode *mode.Flag
	}
)

// NewDefaultEncoder creates a new default JSON encoder
func NewDefaultEncoder(mode *mode.Flag) *DefaultEncoder {
	return &DefaultEncoder{mode: mode}
}

// Encode encodes the body into JSON and writes it to the response
func (d *DefaultEncoder) Encode(
	w http.ResponseWriter,
	response *gonethttpresponse.Response,
) (err error) {
	// Get the response body and HTTP status
	body := response.GetBody(d.mode)
	httpStatus := response.GetHTTPStatus()

	// Encode the JSON body
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return d.Encode(
			w,
			gonethttpstatusresponse.NewDebugInternalServerError(
				err,
				ErrCodeMarshalResponseBodyFailed,
			),
		)
	}

	// Write the JSON body to the response
	w.WriteHeader(httpStatus)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonBody)
	if err != nil {
		return err
	}
	return nil
}
