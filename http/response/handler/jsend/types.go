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
	// ResponsesHandler struct
	ResponsesHandler struct {
		mode *goflagsmode.Flag
		gonethttpresponse.Encoder
	}
)

// NewResponsesHandler creates a new default response handler
//
// Parameters:
//
//   - mode: The flag mode
//   - encoder: The HTTP response encoder
//
// Returns:
//
//   - *ResponsesHandler: The default handler
//   - error: The error if any
func NewResponsesHandler(
	mode *goflagsmode.Flag,
	encoder gonethttpresponse.Encoder,
) (*ResponsesHandler, error) {
	// Check if the flag mode or the encoder is nil
	if mode == nil {
		return nil, goflagsmode.ErrNilModeFlag
	}
	if encoder == nil {
		return nil, gonethttpresponse.ErrNilEncoder
	}

	return &ResponsesHandler{
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
func (r ResponsesHandler) HandleResponse(
	w http.ResponseWriter,
	response gonethttpresponse.Response,
) {
	// Check if the response is nil
	if response == nil {
		r.HandleDebugErrorResponseWithCode(
			w,
			gonethttpresponse.ErrNilResponse,
			gonethttp.ErrInternalServerError,
			ErrCodeNilResponse,
			http.StatusInternalServerError,
		)
		return
	}

	// Call the encoder
	if err := r.Encode(w, response); err != nil {
		r.HandleError(w, err)
		return
	}
}

// HandleError handles the error response
//
// Parameters:
//
//   - w: The HTTP response writer
//   - err: The error to handle
func (r ResponsesHandler) HandleError(
	w http.ResponseWriter,
	err error,
) {
	// Check if the errors is a JSend fail error
	var failErrorTarget *gonethttpresponse.FailError
	if errors.As(err, &failErrorTarget) {
		r.HandleResponse(
			w,
			gonethttpresponsejsend.NewResponseFromFailError(failErrorTarget),
		)
		return
	}

	// Check if the errors is a JSend internal error
	var errorTarget *gonethttpresponse.Error
	if errors.As(err, &errorTarget) {
		r.HandleResponse(
			w,
			gonethttpresponsejsend.NewResponseFromError(errorTarget),
		)
		return
	}

	r.HandleDebugErrorResponseWithCode(
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
func (r ResponsesHandler) HandleErrorResponse(
	w http.ResponseWriter,
	err error,
	httpStatus int,
) {
	r.HandleResponse(
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
func (r ResponsesHandler) HandleErrorResponseWithCode(
	w http.ResponseWriter,
	err error,
	errCode string,
	httpStatus int,
) {
	r.HandleResponse(
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
func (r ResponsesHandler) HandleDebugErrorResponse(
	w http.ResponseWriter,
	debugErr error,
	err error,
	httpStatus int,
) {
	r.HandleResponse(
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
func (r ResponsesHandler) HandleDebugErrorResponseWithCode(
	w http.ResponseWriter,
	debugErr error,
	err error,
	errCode string,
	httpStatus int,
) {
	r.HandleResponse(
		w,
		gonethttpresponsejsend.NewDebugErrorResponseWithCode(
			err.Error(),
			debugErr.Error(),
			errCode,
			httpStatus,
		),
	)
}

// HandleFailErrorResponse handles the fail error response
//
// Parameters:
//
//   - w: The HTTP response writer
//   - field: The field that failed
//   - err: The error to handle
//   - httpStatus: The HTTP status code to return
func (r ResponsesHandler) HandleFailErrorResponse(
	w http.ResponseWriter,
	field string,
	err error,
	httpStatus int,
) {
	r.HandleError(
		w,
		gonethttpresponse.NewFailError(
			field,
			err,
			httpStatus,
		),
	)
}

// HandleFailErrorResponseWithCode handles the fail error response with an error code
//
// Parameters:
//
//   - w: The HTTP response writer
//   - field: The field that failed
//   - err: The error to handle
//   - errCode: The error code to return
//   - httpStatus: The HTTP status code to return
func (r ResponsesHandler) HandleFailErrorResponseWithCode(
	w http.ResponseWriter,
	field string,
	err error,
	errCode string,
	httpStatus int,
) {
	r.HandleError(
		w,
		gonethttpresponse.NewFailErrorWithCode(
			field,
			err,
			errCode,
			httpStatus,
		),
	)
}

// HandleFailResponse handles the fail response
//
// Parameters:
//
//   - w: The HTTP response writer
//   - data: The fail data to handle
//   - httpStatus: The HTTP status code to return
func (r ResponsesHandler) HandleFailResponse(
	w http.ResponseWriter,
	data interface{},
	httpStatus int,
) {
	r.HandleResponse(
		w,
		gonethttpresponsejsend.NewFailResponse(
			data,
			httpStatus,
		),
	)
}

// HandleFailResponseWithCode handles the fail response with an error code
//
// Parameters:
//
//   - w: The HTTP response writer
//   - data: The fail data to handle
//   - errCode: The error code to return
//   - httpStatus: The HTTP status code to return
func (r ResponsesHandler) HandleFailResponseWithCode(
	w http.ResponseWriter,
	data interface{},
	errCode string,
	httpStatus int,
) {
	r.HandleResponse(
		w,
		gonethttpresponsejsend.NewFailResponseWithCode(
			data,
			errCode,
			httpStatus,
		),
	)
}
