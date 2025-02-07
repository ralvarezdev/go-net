package response

type (
	// FailResponseError struct
	FailResponseError struct {
		httpStatus int
		FailBodyError
	}
)

// NewFailResponseError creates a new fail response error
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
func (f *FailResponseError) Key() string {
	return f.key
}

// Error returns the error of the fail response error
func (f *FailResponseError) Error() string {
	return f.err
}

// ErrorCode returns the error code of the fail response error
func (f *FailResponseError) ErrorCode() *string {
	return f.errorCode
}

// HTTPStatus returns the http status of the fail response error
func (f *FailResponseError) HTTPStatus() int {
	return f.httpStatus
}

// NewResponseFromFailResponseError creates a new fail response from a fail response error
func NewResponseFromFailResponseError(
	failResponseError FailResponseError,
) Response {
	return NewResponse(
		NewBodyFromFailBodyError(failResponseError.FailBodyError),
		failResponseError.HTTPStatus(),
	)
}
