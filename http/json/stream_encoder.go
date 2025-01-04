package json

import (
	"encoding/json"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttperrors "github.com/ralvarezdev/go-net/http/errors"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
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

// Encode encodes the data into JSON
func (d *DefaultStreamEncoder) Encode(
	w http.ResponseWriter,
	data interface{},
	code int,
) (err error) {
	// Check the data type
	if err = checkJSONData(w, data, d.mode, d); err != nil {
		return err
	}

	// Encode JSON data and write it to the response
	if err = json.NewEncoder(w).Encode(data); err != nil {
		if d.mode != nil && d.mode.IsDebug() {
			_ = d.Encode(
				w,
				gonethttpresponse.NewJSONErrorResponse(err),
				http.StatusInternalServerError,
			)
		} else {
			_ = d.Encode(
				w,
				gonethttperrors.InternalServerError,
				http.StatusInternalServerError,
			)
		}
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	return nil
}
