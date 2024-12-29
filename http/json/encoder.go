package json

import (
	"encoding/json"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	"net/http"
)

type (
	// Encoder is an interface for encoding JSON data
	Encoder interface {
		Encode(w http.ResponseWriter, data interface{}) (err error)
	}

	// DefaultEncoder is the JSON encoder struct
	DefaultEncoder struct {
		mode *goflagsmode.Flag
	}
)

// NewDefaultEncoder creates a new JSON encoder
func NewDefaultEncoder(mode *goflagsmode.Flag) *DefaultEncoder {
	return &DefaultEncoder{
		mode: mode,
	}
}

// Encode encodes the data into JSON
func (d *DefaultEncoder) Encode(
	w http.ResponseWriter,
	data interface{},
) (err error) {
	// Check the data type
	if err = checkJSONData(w, data, d.mode); err != nil {
		return err
	}

	// Encode JSON data
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}
