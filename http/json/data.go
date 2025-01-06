package json

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttp "github.com/ralvarezdev/go-net/http/errors"
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
		encodeError(
			w,
			gonethttp.InternalServerError,
			err,
			http.StatusInternalServerError,
			mode,
			encoder,
		)
	}
	return err
}

// encodeError encodes the given JSON data
func encodeError(
	w http.ResponseWriter,
	err error,
	debugErr error,
	httpStatus int,
	mode *goflagsmode.Flag,
	encoder Encoder,
) {
	var data interface{}
	var dataStr string

	// Check if the mode is debug and the encoder is nil
	if encoder != nil {
		data = gonethttpresponse.NewDebugErrorResponse(
			err,
			debugErr,
			nil,
			nil,
			httpStatus,
		)
	} else if mode != nil && mode.IsDebug() {
		dataStr = debugErr.Error()
	} else {
		dataStr = err.Error()
	}

	// Encode the data
	if encoder != nil {
		_ = encoder.Encode(
			w,
			data,
			httpStatus,
		)
	} else {
		http.Error(
			w,
			dataStr,
			httpStatus,
		)
	}
}
