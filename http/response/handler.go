package response

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttpjson "github.com/ralvarezdev/go-net/http/json"
	"net/http"
)

type (
	// Handler interface for handling the responses
	Handler interface {
		HandleResponse(w http.ResponseWriter, response interface{}, code int)
		HandleErrorProneResponse(
			w http.ResponseWriter,
			response interface{},
			successCode int,
			err error,
			errorCode int,
		)
		HandleErrorResponse(w http.ResponseWriter, err error, errorCode int)
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

// HandleResponse handles the response
func (d *DefaultHandler) HandleResponse(
	w http.ResponseWriter,
	response interface{},
	code int,
) {
	_ = d.jsonEncoder.Encode(w, response, code)
}

// HandleErrorProneResponse handles the response that may contain an error
func (d *DefaultHandler) HandleErrorProneResponse(
	w http.ResponseWriter,
	response interface{},
	successCode int,
	err error,
	errorCode int,
) {
	// Check if the error is nil
	if err == nil {
		d.HandleResponse(w, response, successCode)
		return
	}

	// Handle the error response
	d.HandleErrorResponse(w, err, errorCode)
}

// HandleErrorResponse handles the error response
func (d *DefaultHandler) HandleErrorResponse(
	w http.ResponseWriter,
	err error,
	errorCode int,
) {
	_ = d.jsonEncoder.Encode(
		w,
		NewErrorResponse(err),
		errorCode,
	)
}
