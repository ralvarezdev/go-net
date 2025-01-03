package handler

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttpjson "github.com/ralvarezdev/go-net/http/json"
	"net/http"
)

type (
	// Handler interface for handling the responses
	Handler interface {
		HandleRequest(
			w http.ResponseWriter,
			r *http.Request,
			data interface{},
		) (err error)
		HandleSuccessResponse(w http.ResponseWriter, response *Response)
		HandleErrorProneResponse(
			w http.ResponseWriter,
			successResponse *Response,
			errorResponse *Response,
		)
		HandleErrorResponse(w http.ResponseWriter, response *Response)
	}

	// DefaultHandler struct
	DefaultHandler struct {
		mode        *goflagsmode.Flag
		jsonEncoder gonethttpjson.Encoder
		jsonDecoder gonethttpjson.Decoder
	}
)

// NewDefaultHandler creates a new default request handler
func NewDefaultHandler(
	mode *goflagsmode.Flag,
	jsonEncoder gonethttpjson.Encoder,
	jsonDecoder gonethttpjson.Decoder,
) (*DefaultHandler, error) {
	// Check if the flag mode, the JSON encoder or the JSON decoder is nil
	if mode == nil {
		return nil, goflagsmode.ErrNilModeFlag
	}
	if jsonEncoder == nil {
		return nil, gonethttpjson.ErrNilJSONEncoder
	}
	if jsonDecoder == nil {
		return nil, gonethttpjson.ErrNilJSONDecoder
	}

	return &DefaultHandler{
		mode:        mode,
		jsonEncoder: jsonEncoder,
		jsonDecoder: jsonDecoder,
	}, nil
}

// HandleRequest handles the request
func (d *DefaultHandler) HandleRequest(
	w http.ResponseWriter,
	r *http.Request,
	data interface{},
) (err error) {
	return d.jsonDecoder.Decode(w, r, data)
}

// HandleSuccessResponse handles the response
func (d *DefaultHandler) HandleSuccessResponse(
	w http.ResponseWriter,
	response *Response,
) {
	if response != nil && response.Code != nil {
		_ = d.jsonEncoder.Encode(w, response.Data, *response.Code)
	} else {
		SendInternalServerError(w)
	}
}

// HandleErrorProneResponse handles the response that may contain an error
func (d *DefaultHandler) HandleErrorProneResponse(
	w http.ResponseWriter,
	successResponse *Response,
	errorResponse *Response,
) {
	// Check if the error response is nil
	if errorResponse != nil {
		d.HandleErrorResponse(w, errorResponse)
		return
	}

	// Handle the success response
	d.HandleSuccessResponse(w, successResponse)
}

// HandleErrorResponse handles the error response
func (d *DefaultHandler) HandleErrorResponse(
	w http.ResponseWriter,
	response *Response,
) {
	if response != nil && response.Code != nil {
		_ = d.jsonEncoder.Encode(w, response.Data, *response.Code)
	} else {
		SendInternalServerError(w)
	}
}
