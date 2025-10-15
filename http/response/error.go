package response

type (
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
