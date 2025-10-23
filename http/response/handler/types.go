package jsend

import (
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gojsonencoder "github.com/ralvarezdev/go-json/encoder"

	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	// ResponsesHandler struct
	ResponsesHandler struct {
		mode *goflagsmode.Flag
		gonethttpresponse.Encoder
		gonethttphandler.RawErrorHandler
	}
)

// NewResponsesHandler creates a new default response handler
//
// Parameters:
//
//   - mode: The flag mode
//   - encoder: The HTTP response encoder
//   - rawErrorHandler: The raw error handler
//
// Returns:
//
//   - *ResponsesHandler: The default handler
//   - error: The error if any
func NewResponsesHandler(
	mode *goflagsmode.Flag,
	encoder gonethttpresponse.Encoder,
	rawErrorHandler gonethttphandler.RawErrorHandler,
) (*ResponsesHandler, error) {
	// Check if the flag mode, the encoder or the raw error handler is nil
	if mode == nil {
		return nil, goflagsmode.ErrNilModeFlag
	}
	if encoder == nil {
		return nil, gojsonencoder.ErrNilEncoder
	}
	if rawErrorHandler == nil {
		return nil, gonethttphandler.ErrNilRawErrorHandler
	}

	return &ResponsesHandler{
		mode,
		encoder,
		rawErrorHandler,
	}, nil
}

// HandleResponse handles the response
//
// Parameters:
//
//   - w: The HTTP response writer
//   - response: The response to handle
func (r ResponsesHandler) HandleResponse(
	w http.ResponseWriter,
	response gonethttpresponse.Response,
) {
	// Check if the response is nil
	if response == nil {
		r.HandleDebugErrorWithCode(
			w,
			gonethttpresponse.ErrNilResponse,
			gonethttp.ErrInternalServerError,
			gonethttpresponse.ErrCodeNilResponse,
			http.StatusInternalServerError,
		)
		return
	}

	// Call the encoder
	if err := r.EncodeAndWriteResponse(w, response); err != nil {
		r.HandleRawError(w, err)
		return
	}
}

// HandleRawError handles the raw error response
//
// Parameters:
//
//   - w: The HTTP response writer
//   - err: The error to handle
func (r ResponsesHandler) HandleRawError(
	w http.ResponseWriter,
	err error,
) {
	r.RawErrorHandler.HandleRawError(w, err, r.HandleResponse)
}

// HandleError handles the error response
//
// Parameters:
//
//   - w: The HTTP response writer
//   - err: The error to handle
//   - httpStatus: The HTTP status code to return
func (r ResponsesHandler) HandleError(
	w http.ResponseWriter,
	err error,
	httpStatus int,
) {
	r.HandleRawError(
		w,
		gonethttpresponse.NewError(
			err,
			httpStatus,
		),
	)
}

// HandleErrorWithCode handles the error response with an error code
//
// Parameters:
//
//   - w: The HTTP response writer
//   - err: The error to handle
//   - errCode: The error code to return
//   - httpStatus: The HTTP status code to return
func (r ResponsesHandler) HandleErrorWithCode(
	w http.ResponseWriter,
	err error,
	errCode string,
	httpStatus int,
) {
	r.HandleRawError(
		w,
		gonethttpresponse.NewErrorWithCode(
			err,
			errCode,
			httpStatus,
		),
	)
}

// HandleDebugError handles the debug error response
//
// Parameters:
//
//   - w: The HTTP response writer
//   - debugErr: The debug error to handle
//   - err: The error to handle
func (r ResponsesHandler) HandleDebugError(
	w http.ResponseWriter,
	debugErr error,
	err error,
	httpStatus int,
) {
	r.HandleRawError(
		w,
		gonethttpresponse.NewDebugError(
			err,
			debugErr,
			httpStatus,
		),
	)
}

// HandleDebugErrorWithCode handles the debug error response with an error code
//
// Parameters:
//
//   - w: The HTTP response writer
//   - debugErr: The debug error to handle
//   - err: The error to handle
//   - errCode: The error code to return
//   - httpStatus: The HTTP status code to return
func (r ResponsesHandler) HandleDebugErrorWithCode(
	w http.ResponseWriter,
	debugErr error,
	err error,
	errCode string,
	httpStatus int,
) {
	r.HandleRawError(
		w,
		gonethttpresponse.NewDebugErrorWithCode(
			err,
			debugErr,
			errCode,
			httpStatus,
		),
	)
}

// HandleFailFieldError handles the fail error response
//
// Parameters:
//
//   - w: The HTTP response writer
//   - field: The field that failed
//   - err: The error to handle
//   - httpStatus: The HTTP status code to return
func (r ResponsesHandler) HandleFailFieldError(
	w http.ResponseWriter,
	field string,
	err error,
	httpStatus int,
) {
	r.HandleRawError(
		w,
		gonethttpresponse.NewFailFieldError(
			field,
			err,
			httpStatus,
		),
	)
}

// HandleFailFieldErrorWithCode handles the fail field error response with an error code
//
// Parameters:
//
//   - w: The HTTP response writer
//   - field: The field that failed
//   - err: The error to handle
//   - errCode: The error code to return
//   - httpStatus: The HTTP status code to return
func (r ResponsesHandler) HandleFailFieldErrorWithCode(
	w http.ResponseWriter,
	field string,
	err error,
	errCode string,
	httpStatus int,
) {
	r.HandleRawError(
		w,
		gonethttpresponse.NewFailFieldErrorWithCode(
			field,
			err,
			errCode,
			httpStatus,
		),
	)
}

// HandleFailDataError handles the fail data error response
//
// Parameters:
//
//   - w: The HTTP response writer
//   - data: The fail data to handle
//   - httpStatus: The HTTP status code to return
func (r ResponsesHandler) HandleFailDataError(
	w http.ResponseWriter,
	data any,
	httpStatus int,
) {
	r.HandleRawError(
		w,
		gonethttpresponse.NewFailDataError(
			data,
			httpStatus,
		),
	)
}

// HandleFailDataErrorWithCode handles the fail data error response with an error code
//
// Parameters:
//
//   - w: The HTTP response writer
//   - data: The fail data to handle
//   - errCode: The error code to return
//   - httpStatus: The HTTP status code to return
func (r ResponsesHandler) HandleFailDataErrorWithCode(
	w http.ResponseWriter,
	data any,
	errCode string,
	httpStatus int,
) {
	r.HandleRawError(
		w,
		gonethttpresponse.NewFailDataErrorWithCode(
			data,
			errCode,
			httpStatus,
		),
	)
}
