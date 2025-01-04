package handler

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttpjson "github.com/ralvarezdev/go-net/http/json"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
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
		HandleResponse(
			w http.ResponseWriter,
			response *gonethttpresponse.Response,
		)
		HandleErrorProneResponse(
			w http.ResponseWriter,
			successResponse *gonethttpresponse.Response,
			errorResponse *gonethttpresponse.Response,
		)
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

// HandleResponse handles the response
func (d *DefaultHandler) HandleResponse(
	w http.ResponseWriter,
	response *gonethttpresponse.Response,
) {
	if response == nil {
		SendInternalServerError(w)
		return
	}

	if gonethttpresponse.Code != nil {
		if gonethttpresponse.DebugData != nil && d.mode != nil && d.mode.IsDebug() {
			_ = d.jsonEncoder.Encode(
				w,
				gonethttpresponse.DebugData,
				*gonethttpresponse.Code,
			)
			return
		}
		_ = d.jsonEncoder.Encode(
			w,
			gonethttpresponse.Data,
			*gonethttpresponse.Code,
		)
	} else {
		if gonethttpresponse.DebugData != nil && d.mode != nil && d.mode.IsDebug() {
			_ = d.jsonEncoder.Encode(
				w,
				gonethttpresponse.DebugData,
				*gonethttpresponse.Code,
			)
			return
		}
		SendInternalServerError(w)
	}
}

// HandleErrorProneResponse handles the response that may contain an error
func (d *DefaultHandler) HandleErrorProneResponse(
	w http.ResponseWriter,
	successResponse *gonethttpresponse.Response,
	errorResponse *gonethttpresponse.Response,
) {
	// Check if the error response is nil
	if errorResponse != nil {
		d.HandleResponse(w, errorResponse)
		return
	}

	// Handle the success response
	d.HandleResponse(w, successResponse)
}
