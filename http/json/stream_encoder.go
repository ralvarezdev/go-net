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
		mode: mode,
	}
}

// Encode encodes the body into JSON and writes it to the response
func (d *DefaultStreamEncoder) Encode(
	w http.ResponseWriter,
	response *gonethttpresponse.Response,
) (err error) {
	// Get the body and HTTP status from the response
	body := response.GetBody(d.mode)
	httpStatus := response.GetHTTPStatus()

	// Encode the JSON body
	if err = json.NewEncoder(w).Encode(body); err != nil {
		_ = d.Encode(
			w,
			gonethttpstatusresponse.NewDebugInternalServerError(
				err,
				ErrCodeMarshalResponseBodyFailed,
			),
		)
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	return nil
}
