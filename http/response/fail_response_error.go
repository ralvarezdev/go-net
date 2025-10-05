package response

type (
	// FailResponseError struct
	FailResponseError struct {
		httpStatus int
		FailBodyError
	}
)

// NewFailResponseError creates a new fail response error
//
// Parameters:
//
//   - key: The key
//   - err: The error
//   - errorCode: The error code
//   - httpStatus: The HTTP status
//
// Returns:
//
//   - *FailResponseError: The fail response error
func NewFailResponseError(
	key, err string, errorCode *string, httpStatus int,
) *FailResponseError {
	return &FailResponseError{
		FailBodyError: FailBodyError{
			key,
			err,
			errorCode,
		},
		httpStatus: httpStatus,
	}
}

// Key returns the key of the fail response error
//
// Returns:
//
//   - string: The key
func (f FailResponseError) Key() string {
	return f.key
}

// Error returns the error of the fail response error
//
// Returns:
//
//   - string: The error
func (f FailResponseError) Error() string {
	return f.err
}

// ErrorCode returns the error code of the fail response error
//
// Returns:
//
//   - *string: The error code
func (f FailResponseError) ErrorCode() *string {
	return f.errorCode
}

// HTTPStatus returns the http status of the fail response error
//
// Returns:
//
//   - int: The http status
func (f FailResponseError) HTTPStatus() int {
	return f.httpStatus
}

// Response creates a new response from a fail response error
//
// Returns:
//
//   - Response: The response
func (f FailResponseError) Response() Response {
	return NewResponse(
		f.FailBodyError.Body(),
		f.HTTPStatus(),
	)
}
