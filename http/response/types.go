package response

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
)

type (
	// Response struct
	Response struct {
		body       interface{}
		debugBody  interface{}
		httpStatus int
	}

	// JSendBody struct
	JSendBody struct {
		Status  string      `json:"status"`
		Data    interface{} `json:"data"`
		Message *string     `json:"message,omitempty"`
		Code    *int        `json:"code,omitempty"`
	}

	// FieldError struct
	FieldError struct {
		Field string
		Err   error
	}
)

// NewJSendSuccessBody creates a new success response body
func NewJSendSuccessBody(
	data interface{},
) *JSendBody {
	return &JSendBody{
		Status: "success",
		Data:   data,
	}
}

// NewJSendFailBody creates a new fail response body
func NewJSendFailBody(
	data interface{},
	code *int,
) *JSendBody {
	return &JSendBody{
		Status: "fail",
		Data:   data,
		Code:   code,
	}
}

// NewJSendErrorBody creates a new error response body
func NewJSendErrorBody(
	message string,
	data interface{},
	code *int,
) *JSendBody {
	return &JSendBody{
		Status:  "error",
		Data:    data,
		Message: &message,
		Code:    code,
	}
}

// newResponse creates a new response
func newResponse(
	body interface{},
	debugBody interface{},
	httpStatus int,
) *Response {
	return &Response{
		body:       body,
		debugBody:  debugBody,
		httpStatus: httpStatus,
	}
}

// GetBody returns the response body
func (r *Response) GetBody(mode *goflagsmode.Flag) interface{} {
	// Check if the response contains the debug response body
	if r.debugBody != nil && mode != nil && mode.IsDebug() {
		return r.debugBody
	}
	return r.body
}

// GetHTTPStatus returns the HTTP status
func (r *Response) GetHTTPStatus() int {
	return r.httpStatus
}

// NewDebugSuccessResponse creates a new success response
func NewDebugSuccessResponse(
	data interface{},
	debugData interface{},
	httpStatus int,
) *Response {
	return newResponse(
		NewJSendSuccessBody(data),
		NewJSendSuccessBody(debugData),
		httpStatus,
	)
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(
	data interface{},
	httpStatus int,
) *Response {
	return NewDebugSuccessResponse(data, data, httpStatus)
}

// NewDebugFailResponse creates a new fail response
func NewDebugFailResponse(
	data interface{},
	debugData interface{},
	errorCode *int,
	httpStatus int,
) *Response {
	return newResponse(
		NewJSendFailBody(data, errorCode),
		NewJSendFailBody(debugData, errorCode),
		httpStatus,
	)
}

// NewFailResponse creates a new fail response
func NewFailResponse(
	data interface{},
	errorCode *int,
	httpStatus int,
) *Response {
	return NewDebugFailResponse(data, data, errorCode, httpStatus)
}

// NewDebugErrorResponse creates a new error response
func NewDebugErrorResponse(
	err error,
	debugErr error,
	data interface{},
	errorCode *int,
	httpStatus int,
) *Response {
	return newResponse(
		NewJSendErrorBody(err.Error(), data, errorCode),
		NewJSendErrorBody(debugErr.Error(), data, errorCode),
		httpStatus,
	)
}

// NewErrorResponse creates a new error response
func NewErrorResponse(
	err error,
	data interface{},
	errorCode *int,
	httpStatus int,
) *Response {
	return NewDebugErrorResponse(err, err, data, errorCode, httpStatus)
}

// NewFieldError creates a new field error
func NewFieldError(
	field string,
	err error,
) *FieldError {
	return &FieldError{
		Field: field,
		Err:   err,
	}
}

// String returns the field error as a string
func (f *FieldError) String() string {
	return f.Err.Error()
}

// NewFieldErrorsBodyData creates a new field errors body data
func NewFieldErrorsBodyData(
	fieldErrors ...FieldError,
) *map[string]*[]string {
	// Check if there are field errors
	if len(fieldErrors) == 0 {
		return nil
	}

	// Initialize the field errors map
	fieldErrorsMap := make(map[string]*[]string)

	// Iterate over the field errors
	for _, fieldError := range fieldErrors {
		// Check if the field name exists in the map
		if _, ok := fieldErrorsMap[fieldError.Field]; !ok {
			// Initialize the field errors slice
			fieldErrorsMap[fieldError.Field] = &[]string{fieldError.Err.Error()}
		} else {
			// Append the error to the field errors slice
			*fieldErrorsMap[fieldError.Field] = append(
				*fieldErrorsMap[fieldError.Field],
				fieldError.Err.Error(),
			)
		}
	}

	return &fieldErrorsMap
}
