package jsend

import (
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	// ErrorBody struct
	ErrorBody struct {
		Status  Status `json:"status"`
		Message string `json:"message,omitempty"`
		Code    string `json:"code,omitempty"`
	}

	// FailError struct
	FailError struct {
		Field      string
		Err        error
		ErrorCode  string
		HTTPStatus int
	}

	// Error struct (would be an error for Error body responses, but to avoid repetition, it's named Error)
	Error struct {
		DebugErr   error
		Err        error
		ErrorCode  string
		HTTPStatus int
	}
)

// NewErrorBodyWithCode creates a new JSend error response body with error code
//
// Parameters:
//
//   - message: The error message
//   - code: The error code
//
// Returns:
//
//   - *ErrorBody: The JSend error body
func NewErrorBodyWithCode(
	message string,
	code string,
) *ErrorBody {
	return &ErrorBody{
		Status:  StatusError,
		Message: message,
		Code:    code,
	}
}

// NewErrorBody creates a new JSend error response body
//
// Parameters:
//
//   - message: The error message
//
// Returns:
//
//   - *ErrorBody: The JSend error body
func NewErrorBody(
	message string,
) *ErrorBody {
	return NewErrorBodyWithCode(message, "")
}

// NewErrorResponseWithCode creates a new JSend error response with error code
//
// Parameters:
//
//   - message: The error message
//   - code: The error code
//   - httpStatus: The HTTP status code
//
// Returns:
//
//   - Response: The response
func NewErrorResponseWithCode(
	message string,
	code string,
	httpStatus int,
) gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(
		NewErrorBodyWithCode(message, code),
		httpStatus,
	)
}

// NewErrorResponse creates a new JSend error response
//
// Parameters:
//
//   - message: The error message
//   - httpStatus: The HTTP status code
//
// Returns:
//
//   - Response: The response
func NewErrorResponse(
	message string,
	httpStatus int,
) gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(
		NewErrorBody(message),
		httpStatus,
	)
}

// NewDebugErrorResponseWithCode creates a new JSend error response with debug information and error code
//
// Parameters:
//
//   - message: The error message
//   - debugMessage: The debug message
//   - code: The error code
//   - httpStatus: The HTTP status code
//
// Returns:
//
//   - Response: The response
func NewDebugErrorResponseWithCode(
	message string,
	debugMessage string,
	code string,
	httpStatus int,
) gonethttpresponse.Response {
	if debugMessage == "" {
		return NewErrorResponseWithCode(message, code, httpStatus)
	}
	return gonethttpresponse.NewDebugResponse(
		NewErrorBodyWithCode(message, code),
		NewErrorBodyWithCode(debugMessage, code),
		httpStatus,
	)
}

// NewDebugErrorResponse creates a new JSend error response with debug information
//
// Parameters:
//
//   - message: The error message
//   - debugMessage: The debug message
//   - httpStatus: The HTTP status code
//
// Returns:
//
//   - Response: The response
func NewDebugErrorResponse(
	message string,
	debugMessage string,
	httpStatus int,
) gonethttpresponse.Response {
	return NewDebugErrorResponseWithCode(message, debugMessage, "", httpStatus)
}

// NewFailErrorWithCode creates a new JSend fail error with error code
//
// Parameters:
//
//   - field: The field
//   - err: The field error
//   - errorCode: The error code
//   - httpStatus: The HTTP status
//
// Returns:
//
//   - *FailError: The JSend fail error
func NewFailErrorWithCode(
	field string, err error, errorCode string, httpStatus int,
) *FailError {
	return &FailError{
		Field:      field,
		Err:        err,
		ErrorCode:  errorCode,
		HTTPStatus: httpStatus,
	}
}

// NewFailError creates a new JSend fail error
//
// Parameters:
//
//   - field: The field
//   - err: The field error
//   - httpStatus: The HTTP status
//
// Returns:
//
//   - *FailError: The JSend fail error
//   - string: The error code
func NewFailError(
	field string, err error, httpStatus int,
) *FailError {
	return NewFailErrorWithCode(field, err, "", httpStatus)
}

// Error returns the error message from the fail body error
//
// Returns:
//
//   - string: The error message
func (f FailError) Error() string {
	return f.Err.Error()
}

// Data returns a response data map from the fail body error
//
// Returns:
//
//   - map[string][]string: The response data map
func (f FailError) Data() map[string][]string {
	// Initialize the data map
	data := make(map[string][]string)

	// Add the fail body error to the data map
	data[f.Field] = []string{f.Err.Error()}

	return data
}

// Response creates a new response from a JSend fail error
//
// Returns:
//
//   - Response: The response
func (f FailError) Response() gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(
		NewFailBodyWithCode(
			f.Data(),
			f.ErrorCode,
		),
		f.HTTPStatus,
	)
}

// NewErrorWithCode creates a new JSend error with error code
//
// Parameters:
//
//   - err: The error
//   - errCode: The error code
//   - httpStatus: The HTTP status code
//
// Returns:
//
//   - *Error: The JSend error
func NewErrorWithCode(
	err error,
	errCode string,
	httpStatus int,
) *Error {
	return &Error{
		DebugErr:   nil,
		Err:        err,
		ErrorCode:  errCode,
		HTTPStatus: httpStatus,
	}
}

// NewError creates a new JSend error
//
// Parameters:
//
//   - err: The error
//   - httpStatus: The HTTP status code
//
// Returns:
//
//   - *Error: The JSend error
func NewError(
	err error,
	httpStatus int,
) *Error {
	return NewErrorWithCode(err, "", httpStatus)
}

// NewDebugErrorWithCode creates a new JSend error with debug information and error code
//
// Parameters:
//
//   - debugErr: The debug error
//   - err: The error
//   - errCode: The error code
//   - httpStatus: The HTTP status code
//
// Returns:
func NewDebugErrorWithCode(
	debugErr error,
	err error,
	errCode string,
	httpStatus int,
) *Error {
	if debugErr == nil {
		return NewErrorWithCode(err, errCode, httpStatus)
	}
	return &Error{
		DebugErr:   debugErr,
		Err:        err,
		ErrorCode:  errCode,
		HTTPStatus: httpStatus,
	}
}

// NewDebugError creates a new JSend error with debug information
//
// Parameters:
//
//   - debugErr: The debug error
//   - err: The error
//   - httpStatus: The HTTP status code
//
// Returns:
//
//   - *Error: The JSend error
func NewDebugError(
	debugErr error,
	err error,
	httpStatus int,
) *Error {
	return NewDebugErrorWithCode(debugErr, err, "", httpStatus)
}

// Error returns the error message from the error
//
// Returns:
//
//   - string: The error message
func (e Error) Error() string {
	return e.Err.Error()
}

// Response creates a new response from a JSend error
//
// Returns:
//
//   - Response: The response
func (e Error) Response() gonethttpresponse.Response {
	if e.DebugErr != nil {
		return NewDebugErrorResponseWithCode(
			e.Err.Error(),
			e.DebugErr.Error(),
			e.ErrorCode,
			e.HTTPStatus,
		)
	}
	return NewErrorResponseWithCode(
		e.Err.Error(),
		e.ErrorCode,
		e.HTTPStatus,
	)
}
