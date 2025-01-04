package handler

import (
	"errors"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpjson "github.com/ralvarezdev/go-net/http/json"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"net/http"
)

// SendInternalServerError sends an internal server error response
func SendInternalServerError(
	w http.ResponseWriter,
	encoder gonethttpjson.Encoder,
) {
	if encoder != nil {
		_ = encoder.Encode(
			w,
			gonethttpresponse.NewJSONErrorResponse(errors.New(gonethttp.InternalServerError)),
			http.StatusInternalServerError,
		)
	} else {
		http.Error(
			w,
			gonethttp.InternalServerError,
			http.StatusInternalServerError,
		)
	}
}
