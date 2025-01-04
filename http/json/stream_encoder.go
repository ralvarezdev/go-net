package json

import (
	"encoding/json"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"net/http"
)

type (
	// StreamEncoder is an interface for encoding JSON data
	StreamEncoder interface {
		Encode(w http.ResponseWriter, data interface{}, code int) (err error)
	}

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
	if err = checkJSONData(w, data, d.mode); err != nil {
		return err
	}

	// Encode JSON data and write it to the response
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(data); err != nil {
		_ = d.Encode(
			w,
			gonethttpresponse.NewJSONErrorResponse(err),
			http.StatusInternalServerError,
		)
		return err
	}
	return nil
}
