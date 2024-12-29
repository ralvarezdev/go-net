package json

import (
	"encoding/json"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	"net/http"
)

type (
	// StreamDecoder is the JSON decoder interface
	StreamDecoder interface {
		Decode(
			w http.ResponseWriter,
			r *http.Request,
			data interface{},
		) (err error)
	}

	// DefaultStreamDecoder is the JSON decoder struct
	DefaultStreamDecoder struct {
		mode *goflagsmode.Flag
	}
)

// NewDefaultStreamDecoder creates a new JSON decoder
func NewDefaultStreamDecoder(mode *goflagsmode.Flag) *DefaultStreamDecoder {
	return &DefaultStreamDecoder{
		mode: mode,
	}
}

// Decode decodes the JSON data
func (d *DefaultStreamDecoder) Decode(
	w http.ResponseWriter,
	r *http.Request,
	data interface{},
) (err error) {
	// Check the data type
	if err = checkJSONData(w, data, d.mode); err != nil {
		return err
	}

	// Decode JSON data
	if err = json.NewDecoder(r.Body).Decode(data); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
	}
	return err
}
