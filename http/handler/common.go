package handler

import (
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttperrors "github.com/ralvarezdev/go-net/http/errors"
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
			gonethttpresponse.NewErrorResponse(
				gonethttperrors.InternalServerError,
				nil, nil,
				http.StatusInternalServerError,
			),
		)
	} else {
		http.Error(
			w,
			gonethttp.InternalServerError,
			http.StatusInternalServerError,
		)
	}
}
