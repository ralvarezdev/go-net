package json

import (
	"encoding/json"
	"github.com/ralvarezdev/go-flags/mode"
	gonethttperrors "github.com/ralvarezdev/go-net/http/errors"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"net/http"
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

// Encode encodes the given data to JSON
func (d *DefaultEncoder) Encode(
	w http.ResponseWriter,
	response *gonethttpresponse.Response,
) (err error) {
	// Get the response body and HTTP status
	body := response.GetBody(d.mode)
	httpStatus := response.GetHTTPStatus()

	// Check the data type
	if err = checkJSONData(w, body, d.mode, d); err != nil {
		return err
	}

	// Encode the data
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return d.Encode(
			w,
			gonethttpresponse.NewDebugErrorResponse(
				gonethttperrors.InternalServerError,
				err,
				nil,
				nil,
				http.StatusBadRequest,
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
