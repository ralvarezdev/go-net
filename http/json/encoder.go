package json

import (
	"encoding/json"
	"github.com/ralvarezdev/go-flags/mode"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"net/http"
)

type (
	// Encoder interface
	Encoder interface {
		Encode(
			w http.ResponseWriter,
			data interface{},
			code int,
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
	data interface{},
	code int,
) (err error) {
	// Check the data type
	if err = checkJSONData(w, data, d.mode); err != nil {
		return err
	}

	// Encode the data
	jsonData, err := json.Marshal(data)
	if err != nil {
		_ = d.Encode(
			w,
			gonethttpresponse.NewJSONErrorResponse(err),
			http.StatusInternalServerError,
		)
		return err
	}

	// Write the JSON data
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}
