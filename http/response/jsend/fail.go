package jsend

import (
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	// FailBody struct
	FailBody struct {
		Status Status      `json:"status"`
		Code   string      `json:"code,omitempty"`
		Data   interface{} `json:"data,omitempty"`
	}

	// FailError struct
	FailError struct {
		httpStatus int
		field      string
		err        string
		errorCode  string
	}
)

// NewFailBodyWithCode creates a new JSend fail response body with error code
//
// Parameters:
//
//   - data: The data
//   - code: The error code
//
// Returns:
//
//   - *FailBody: The JSend fail body
func NewFailBodyWithCode(
	data interface{},
	code string,
) *FailBody {
	return &FailBody{
		Status: StatusFail,
		Code:   code,
		Data:   data,
	}
}

// NewFailBody creates a new JSend fail response body
//
// Parameters:
//
//   - data: The data
//
// Returns:
func NewFailBody(
	data interface{},
) *FailBody {
	return NewFailBodyWithCode(data, "")
}

// NewFailResponseWithCode creates a new JSend fail response with error code
//
// Parameters:
//
//   - data: The data
//   - code: The error code
//   - httpStatus: The HTTP status code
//
// Returns:
//
//   - Response: The response
func NewFailResponseWithCode(
	data interface{},
	code string,
	httpStatus int,
) gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(
		NewFailBodyWithCode(data, code),
		httpStatus,
	)
}

// NewFailResponse creates a new JSend fail response
//
// Parameters:
//
//   - data: The data
//   - httpStatus: The HTTP status code
//
// Returns:
//
//   - Response: The response
func NewFailResponse(
	data interface{},
	httpStatus int,
) gonethttpresponse.Response {
	return NewFailResponseWithCode(data, "", httpStatus)
}

// NewFailErrorWithCode creates a new fail response error with error code
//
// Parameters:
//
//   - field: The field
//   - err: The error
//   - errorCode: The error code
//   - httpStatus: The HTTP status
//
// Returns:
//
//   - *FailError: The fail response error
func NewFailErrorWithCode(
	field, err string, errorCode string, httpStatus int,
) *FailError {
	return &FailError{
		field:      field,
		err:        err,
		errorCode:  errorCode,
		httpStatus: httpStatus,
	}
}

// NewFailError creates a new fail response error
//
// Parameters:
//
//   - field: The field
//   - err: The error
//   - httpStatus: The HTTP status
//
// Returns:
//
//   - *FailError: The fail response error
//   - string: The error code
func NewFailError(
	field, err string, httpStatus int,
) *FailError {
	return NewFailErrorWithCode(field, err, "", httpStatus)
}

// Field returns the field of the fail response error
//
// Returns:
//
//   - string: The field
func (f FailError) Field() string {
	return f.field
}

// Error returns the error of the fail response error
//
// Returns:
//
//   - string: The error
func (f FailError) Error() string {
	return f.err
}

// ErrorCode returns the error code of the fail response error
//
// Returns:
//
//   - string: The error code
func (f FailError) ErrorCode() string {
	return f.errorCode
}

// HTTPStatus returns the http status of the fail response error
//
// Returns:
//
//   - int: The http status
func (f FailError) HTTPStatus() int {
	return f.httpStatus
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
	data[f.Field()] = []string{f.Error()}

	return data
}

// Body returns a response body from the fail body error
//
// Returns:
//
//   - *JSendFailBody: The response body
func (f FailError) Body() *FailBody {
	return NewFailBodyWithCode(
		f.Data(),
		f.ErrorCode(),
	)
}

// Response creates a new response from a fail response error
//
// Returns:
//
//   - Response: The response
func (f FailError) Response() gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(
		f.Body(),
		f.HTTPStatus(),
	)
}
