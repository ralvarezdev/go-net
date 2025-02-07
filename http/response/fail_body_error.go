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
func (f *FailBodyError) Key() string {
	return f.key
}

// Error returns the error of the fail body error
func (f *FailBodyError) Error() string {
	return f.err
}

// ErrorCode returns the error code of the fail body error
func (f *FailBodyError) ErrorCode() *string {
	return f.errorCode
}

// NewDataFromFailBodyError creates a new data map from fail body error
func NewDataFromFailBodyError(
	failBodyError *FailBodyError,
) *map[string]*[]string {
	// Check if the fail body error is nil
	if failBodyError == nil {
		return nil
	}

	// Initialize the data map
	data := make(map[string]*[]string)

	// Add the fail body error to the data map
	data[failBodyError.Key()] = &[]string{failBodyError.Error()}

	return &data
}

// NewBodyFromFailBodyError creates a new body from fail body error
func NewBodyFromFailBodyError(
	failBodyError *FailBodyError,
) *JSendBody {
	return NewJSendFailBody(
		NewDataFromFailBodyError(failBodyError),
		failBodyError.ErrorCode(),
	)
}
