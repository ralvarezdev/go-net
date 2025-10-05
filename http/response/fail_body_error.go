package response

type (
	// FailBodyError struct
	FailBodyError struct {
		key       string
		err       string
		errorCode *string
	}
)

// NewFailBodyError creates a new fail body error
//
// Parameters:
//
//   - key: The key of the fail body error
//   - err: The error of the fail body error
//   - errorCode: The error code of the fail body error
//
// Returns:
//
//   - *FailBodyError: The fail body error
func NewFailBodyError(
	key, err string, errorCode *string,
) *FailBodyError {
	return &FailBodyError{
		key,
		err,
		errorCode,
	}
}

// Key returns the key of the fail body error
//
// Returns:
//
//   - string: The key of the fail body error
func (f FailBodyError) Key() string {
	return f.key
}

// Error returns the error of the fail body error
//
// Returns:
//
//   - string: The error of the fail body error
func (f FailBodyError) Error() string {
	return f.err
}

// ErrorCode returns the error code of the fail body error
//
// Returns:
//
//   - *string: The error code of the fail body error
func (f FailBodyError) ErrorCode() *string {
	return f.errorCode
}

// Data returns a response data map from the fail body error
//
// Returns:
//
//   - map[string][]string: The response data map
func (f FailBodyError) Data() map[string][]string {
	// Initialize the data map
	data := make(map[string][]string)

	// Add the fail body error to the data map
	data[f.Key()] = []string{f.Error()}

	return data
}

// Body returns a response body from the fail body error
//
// Returns:
//
//   - *JSendFailBody: The response body
func (f FailBodyError) Body() *JSendFailBody {
	return NewJSendFailBody(
		f.Data(),
		f.ErrorCode(),
	)
}
