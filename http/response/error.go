package response

import (
	"fmt"
)

type (
	// FailFieldError struct
	FailFieldError struct {
		Field      string
		Err        error
		ErrorCode  string
		HTTPStatus int
	}

	// FailDataError struct
	FailDataError struct {
		Data       interface{}
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

// NewFailFieldErrorWithCode creates a new JSend fail field error with error code
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
//   - *FailFieldError: The JSend fail field error
func NewFailFieldErrorWithCode(
	field string, err error, errorCode string, httpStatus int,
) *FailFieldError {
	return &FailFieldError{
		Field:      field,
		Err:        err,
		ErrorCode:  errorCode,
		HTTPStatus: httpStatus,
	}
}

// NewFailFieldError creates a new JSend fail field error
//
// Parameters:
//
//   - field: The field
//   - err: The field error
//   - httpStatus: The HTTP status
//
// Returns:
//
//   - *FailFieldError: The JSend fail field error
//   - string: The error code
func NewFailFieldError(
	field string, err error, httpStatus int,
) *FailFieldError {
	return NewFailFieldErrorWithCode(field, err, "", httpStatus)
}

// Error returns the error message from the fail body error
//
// Returns:
//
//   - string: The error message
func (f FailFieldError) Error() string {
	return f.Err.Error()
}

// Data returns a response data map from the fail body error
//
// Returns:
//
//   - map[string][]string: The response data map
func (f FailFieldError) Data() map[string][]string {
	// Initialize the data map
	data := make(map[string][]string)

	// Add the fail body error to the data map
	data[f.Field] = []string{f.Err.Error()}

	return data
}

// NewFailDataErrorWithCode creates a new JSend fail data error with error code
//
// Parameters:
//
//   - data: The fail data
//   - errorCode: The error code
//   - httpStatus: The HTTP status
//
// Returns:
//
//   - *FailDataError: The JSend fail data error
func NewFailDataErrorWithCode(
	data interface{},
	errorCode string,
	httpStatus int,
) *FailDataError {
	return &FailDataError{
		Data:       data,
		ErrorCode:  errorCode,
		HTTPStatus: httpStatus,
	}
}

// NewFailDataError creates a new JSend fail data error
//
// Parameters:
//
//   - data: The fail data
//   - httpStatus: The HTTP status
//
// Returns:
//
//   - *FailDataError: The JSend fail data error
func NewFailDataError(
	data interface{},
	httpStatus int,
) *FailDataError {
	return NewFailDataErrorWithCode(data, "", httpStatus)
}

// Error returns the data field as a string from the fail data error
//
// Returns:
//
//   - string: The error message
func (f FailDataError) Error() string {
	return fmt.Sprintf("%v", f.Data)
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
