package json

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttp "github.com/ralvarezdev/go-net/http"
	"net/http"
)

// checkJSONData checks if the given JSON data is nil
func checkJSONData(
	w http.ResponseWriter,
	data interface{},
	mode *goflagsmode.Flag,
) error {
	// Check if data is nil
	if data == nil {
		return handleDataTypeError(w, ErrNilJSONData, mode)
	}
	return nil
}

// handleDataTypeError handles the data type error
func handleDataTypeError(
	w http.ResponseWriter,
	err error,
	mode *goflagsmode.Flag,
) error {
	if mode != nil && mode.IsDebug() {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
	}
	http.Error(
		w,
		gonethttp.InternalServerError,
		http.StatusInternalServerError,
	)
	return err
}
