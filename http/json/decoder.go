package json

import (
	"encoding/json"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	"io"
	"net/http"
)

type (
	// Decoder interface
	Decoder interface {
		Decode(
			w http.ResponseWriter,
			r *http.Request,
			data interface{},
		) (err error)
	}

	// DefaultDecoder struct
	DefaultDecoder struct {
		mode *goflagsmode.Flag
	}
)

// NewDefaultDecoder creates a new JSON decoder
func NewDefaultDecoder(mode *goflagsmode.Flag) *DefaultDecoder {
	return &DefaultDecoder{
		mode: mode,
	}
}

// Decode decodes the JSON data from the request
func (d *DefaultDecoder) Decode(
	w http.ResponseWriter,
	r *http.Request,
	data interface{},
) (err error) {
	// Check the data type
	if err = checkJSONData(w, data, d.mode); err != nil {
		return err
	}

	// Get the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return err
	}

	// Decode JSON data
	if err = json.Unmarshal(body, data); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
	}
	return err
}
