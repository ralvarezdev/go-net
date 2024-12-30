package response

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpjson "github.com/ralvarezdev/go-net/http/json"
	"net/http"
)

type (
	// Handler interface for handling the responses
	Handler interface {
		HandleSuccess(w http.ResponseWriter, response *Response)
		HandleErrorProne(
			w http.ResponseWriter,
			successResponse *Response,
			errorResponse *Response,
		)
		HandleError(w http.ResponseWriter, response *Response)
	}

	// DefaultHandler struct
	DefaultHandler struct {
		mode        *goflagsmode.Flag
		jsonEncoder gonethttpjson.Encoder
	}
)

// NewDefaultHandler creates a new default request handler
func NewDefaultHandler(
	mode *goflagsmode.Flag,
	jsonEncoder gonethttpjson.Encoder,
) (*DefaultHandler, error) {
	// Check if the flag mode or the JSON encoder is nil
	if mode == nil {
		return nil, goflagsmode.ErrNilModeFlag
	}
	if jsonEncoder == nil {
		return nil, gonethttpjson.ErrNilJSONEncoder
	}
	return &DefaultHandler{mode: mode, jsonEncoder: jsonEncoder}, nil
}

// HandleSuccess handles the response
func (d *DefaultHandler) HandleSuccess(
	w http.ResponseWriter,
	response *Response,
) {
	if response != nil && response.Code != nil {
		_ = d.jsonEncoder.Encode(w, response.Data, *response.Code)
	} else {
		http.Error(
			w,
			gonethttp.InternalServerError,
			http.StatusInternalServerError,
		)
	}
}

// HandleErrorProne handles the response that may contain an error
func (d *DefaultHandler) HandleErrorProne(
	w http.ResponseWriter,
	successResponse *Response,
	errorResponse *Response,
) {
	// Check if the error response is nil
	if errorResponse != nil {
		d.HandleError(w, errorResponse)
		return
	}

	// Handle the success response
	d.HandleSuccess(w, successResponse)
}

// HandleError handles the error response
func (d *DefaultHandler) HandleError(
	w http.ResponseWriter,
	response *Response,
) {
	if response != nil && response.Code != nil {
		_ = d.jsonEncoder.Encode(w, response.Data, *response.Code)
	} else {
		http.Error(
			w,
			gonethttp.InternalServerError,
			http.StatusInternalServerError,
		)
	}
}
