package json

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"net/http"
)

// checkJSONData checks if the given JSON data is nil
func checkJSONData(
	w http.ResponseWriter,
	data interface{},
	mode *goflagsmode.Flag,
	encoder Encoder,
) (err error) {
	// Check if data is nil
	if data == nil {
		err = ErrNilJSONData
		encodeError(w, err, http.StatusInternalServerError, mode, encoder)
	}
	return err
}

// encodeError encodes the given JSON data
func encodeError(
	w http.ResponseWriter,
	err error,
	code int,
	mode *goflagsmode.Flag,
	encoder Encoder,
) {
	var data interface{}
	var dataStr string

	// Check if the mode is debug and the encoder is nil
	if mode != nil && mode.IsDebug() {
		if encoder != nil {
			data = gonethttpresponse.NewJSONErrorResponse(err)
		} else {
			dataStr = err.Error()
		}
	} else if encoder != nil {
		data = gonethttpresponse.InternalServerError
	} else {
		dataStr = gonethttp.InternalServerError
	}

	// Encode the data
	if encoder != nil {
		_ = encoder.Encode(
			w,
			data,
			code,
		)
	} else {
		http.Error(
			w,
			dataStr,
			code,
		)
	}
}
