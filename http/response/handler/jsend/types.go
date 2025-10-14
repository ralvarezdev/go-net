package jsend

import (
	"errors"
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpresponsejsend "github.com/ralvarezdev/go-net/http/response/jsend"
)

type (
	// Handler struct
	Handler struct {
		mode *goflagsmode.Flag
		gonethttpresponse.Encoder
	}
)

// NewHandler creates a new default response handler
//
// Parameters:
//
//   - mode: The flag mode
//   - encoder: The HTTP response encoder
//
// Returns:
//
//   - *Handler: The default handler
//   - error: The error if any
func NewHandler(
	mode *goflagsmode.Flag,
	encoder gonethttpresponse.Encoder,
) (*Handler, error) {
	// Check if the flag mode or the encoder is nil
	if mode == nil {
		return nil, goflagsmode.ErrNilModeFlag
	}
	if encoder == nil {
		return nil, gonethttpresponse.ErrNilEncoder
	}

	return &Handler{
		mode,
		encoder,
	}, nil
}

// HandleResponse handles the response
//
// Parameters:
//
//   - w: The HTTP response writer
//   - response: The response to handle
func (d Handler) HandleResponse(
	w http.ResponseWriter,
	response gonethttpresponse.Response,
) {
	// Check if the response is nil
	if response == nil {
		d.HandleDebugErrorResponseWithCode(
			w,
			gonethttpresponse.ErrNilResponse,
			gonethttp.ErrInternalServerError,
			ErrCodeNilResponse,
			http.StatusInternalServerError,
		)
		return
	}

	// Call the encoder
	_ = d.Encode(w, response)
}

// HandleError handles the error response
//
// Parameters:
//
//   - w: The HTTP response writer
//   - err: The error to handle
func (d Handler) HandleError(
	w http.ResponseWriter,
	err error,
) {
	// Check if the errors is a fail body error or a fail request error
	var failResponseErrorTarget *gonethttpresponsejsend.FailError
	if errors.As(err, &failResponseErrorTarget) {
		d.HandleResponse(
			w,
			failResponseErrorTarget.Response(),
		)
		return
	}

	d.HandleDebugErrorResponseWithCode(
		w,
		err,
		gonethttp.ErrInternalServerError,
		ErrCodeRequestFatalError,
		http.StatusInternalServerError,
	)
}

// HandleErrorResponse handles the error response
//
// Parameters:
//
//   - w: The HTTP response writer
//   - err: The error to handle
//   - httpStatus: The HTTP status code to return
func (d Handler) HandleErrorResponse(
	w http.ResponseWriter,
	err error,
	httpStatus int,
) {
	d.HandleResponse(
		w,
		gonethttpresponsejsend.NewErrorResponse(
			err.Error(),
			httpStatus,
		),
	)
}

// HandleErrorResponseWithCode handles the error response with an error code
//
// Parameters:
//
//   - w: The HTTP response writer
//   - err: The error to handle
//   - errCode: The error code to return
//   - httpStatus: The HTTP status code to return
func (d Handler) HandleErrorResponseWithCode(
	w http.ResponseWriter,
	err error,
	errCode string,
	httpStatus int,
) {
	d.HandleResponse(
		w,
		gonethttpresponsejsend.NewErrorResponseWithCode(
			err.Error(),
			errCode,
			httpStatus,
		),
	)
}

// HandleDebugErrorResponse handles the debug error response
//
// Parameters:
//
//   - w: The HTTP response writer
//   - debugErr: The debug error to handle
//   - err: The error to handle
func (d Handler) HandleDebugErrorResponse(
	w http.ResponseWriter,
	debugErr error,
	err error,
	httpStatus int,
) {
	d.HandleResponse(
		w,
		gonethttpresponsejsend.NewDebugErrorResponse(
			err.Error(),
			debugErr.Error(),
			httpStatus,
		),
	)
}

// HandleDebugErrorResponseWithCode handles the debug error response with an error code
//
// Parameters:
//
//   - w: The HTTP response writer
//   - debugErr: The debug error to handle
//   - err: The error to handle
//   - errCode: The error code to return
//   - httpStatus: The HTTP status code to return
func (d Handler) HandleDebugErrorResponseWithCode(
	w http.ResponseWriter,
	debugErr error,
	err error,
	errCode string,
	httpStatus int,
) {
	d.HandleResponse(
		w,
		gonethttpresponsejsend.NewDebugErrorResponseWithCode(
			err.Error(),
			debugErr.Error(),
			errCode,
			httpStatus,
		),
	)
}

// HandleFieldFailResponse handles the field fail response
//
// Parameters:
//
//   - w: The HTTP response writer
//   - field: The field that failed
//   - err: The error to handle
//   - httpStatus: The HTTP status code to return
func (d Handler) HandleFieldFailResponse(
	w http.ResponseWriter,
	field string,
	err error,
	httpStatus int,
) {
	d.HandleError(
		w,
		gonethttpresponsejsend.NewFailError(
			field,
			err.Error(),
			httpStatus,
		),
	)
}

// HandleFieldFailResponseWithCode handles the field fail response with an error code
//
// Parameters:
//
//   - w: The HTTP response writer
//   - field: The field that failed
//   - err: The error to handle
//   - errCode: The error code to return
//   - httpStatus: The HTTP status code to return
func (d Handler) HandleFieldFailResponseWithCode(
	w http.ResponseWriter,
	field string,
	err error,
	errCode string,
	httpStatus int,
) {
	d.HandleError(
		w,
		gonethttpresponsejsend.NewFailErrorWithCode(
			field,
			err.Error(),
			errCode,
			httpStatus,
		),
	)
}
