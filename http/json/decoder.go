package json

import (
	"encoding/json"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	"net/http"
)

type (
	// Decoder is the JSON decoder interface
	Decoder interface {
		Decode(
			w http.ResponseWriter,
			r *http.Request,
			data interface{},
		) (err error)
	}

	// DefaultDecoder is the JSON decoder struct
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

// Decode decodes the JSON data
func (d *DefaultDecoder) Decode(
	w http.ResponseWriter,
	r *http.Request,
	data interface{},
) (err error) {
	// Check the data type
	if err = checkJSONData(w, data, d.mode); err != nil {
		return err
	}

	// Decode JSON data
	err = json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return err
	}
	return nil
}
