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
		Code    *string     `json:"code,omitempty"`
	}

	// RequestError struct
	RequestError interface {
		Key() string
		Error() string
		HTTPStatus() int
		ErrorCode() *string
	}

	// FieldError struct
	FieldError struct {
		field      string
		err        string
		httpStatus int
		errorCode  *string
	}

	// HeaderError struct
	HeaderError struct {
		header     string
		err        string
		httpStatus int
		errorCode  *string
	}

	// ParameterError struct
	ParameterError struct {
		parameter  string
		err        string
		httpStatus int
		errorCode  *string
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
	code *string,
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
	code *string,
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
	errorCode *string,
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
	errorCode *string,
	httpStatus int,
) *Response {
	return NewDebugFailResponse(data, data, errorCode, httpStatus)
}

// NewDebugErrorResponse creates a new error response
func NewDebugErrorResponse(
	err error,
	debugErr error,
	data interface{},
	errorCode *string,
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
	errorCode *string,
	httpStatus int,
) *Response {
	return NewDebugErrorResponse(err, err, data, errorCode, httpStatus)
}

// NewFieldError creates a new field error
func NewFieldError(
	field, err string, httpStatus int, errorCode ...string,
) *FieldError {
	// Check if the error code is provided
	if len(errorCode) > 0 {
		return &FieldError{
			field:      field,
			err:        err,
			httpStatus: httpStatus,
			errorCode:  &errorCode[0],
		}
	}

	return &FieldError{
		field:      field,
		err:        err,
		httpStatus: httpStatus,
	}
}

// Key returns the field name
func (f *FieldError) Key() string {
	return f.field
}

// Error returns the field error as a string
func (f *FieldError) Error() string {
	return f.err
}

// HTTPStatus returns the HTTP status
func (f *FieldError) HTTPStatus() int {
	return f.httpStatus
}

// ErrorCode returns the error code
func (f *FieldError) ErrorCode() *string {
	return f.errorCode
}

// NewHeaderError creates a new header error
func NewHeaderError(
	header, err string, httpStatus int, errorCode ...string,
) *HeaderError {
	// Check if the error code is provided
	if len(errorCode) > 0 {
		return &HeaderError{
			header:     header,
			err:        err,
			httpStatus: httpStatus,
			errorCode:  &errorCode[0],
		}
	}

	return &HeaderError{
		header:     header,
		err:        err,
		httpStatus: httpStatus,
	}
}

// Key returns the header name
func (h *HeaderError) Key() string {
	return h.header
}

// Error returns the header error as a string
func (h *HeaderError) Error() string {
	return h.err
}

// HTTPStatus returns the HTTP status
func (h *HeaderError) HTTPStatus() int {
	return h.httpStatus
}

// ErrorCode returns the error code
func (h *HeaderError) ErrorCode() *string {
	return h.errorCode
}

// NewParameterError creates a new parameter error
func NewParameterError(
	parameter, err string, httpStatus int, errorCode ...string,
) *ParameterError {
	// Check if the error code is provided
	if len(errorCode) > 0 {
		return &ParameterError{
			parameter:  parameter,
			err:        err,
			httpStatus: httpStatus,
			errorCode:  &errorCode[0],
		}
	}

	return &ParameterError{
		parameter:  parameter,
		err:        err,
		httpStatus: httpStatus,
	}
}

// Key returns the parameter name
func (p *ParameterError) Key() string {
	return p.parameter
}

// Error returns the parameter error as a string
func (p *ParameterError) Error() string {
	return p.err
}

// HTTPStatus returns the HTTP status
func (p *ParameterError) HTTPStatus() int {
	return p.httpStatus
}

// ErrorCode returns the error code
func (p *ParameterError) ErrorCode() *string {
	return p.errorCode
}

// NewRequestErrorBodyData creates a new request errors body data
func NewRequestErrorBodyData(
	requestError RequestError,
) *map[string]*[]string {
	// Initialize the request errors map
	requestErrorsMap := make(map[string]*[]string)

	// Add the request error to the request errors map
	requestErrorsMap[requestError.Key()] = &[]string{requestError.Error()}

	return &requestErrorsMap
}

// NewFailResponseFromRequestError creates a new fail response from a request error
func NewFailResponseFromRequestError(
	requestError RequestError,
) *Response {
	return NewFailResponse(
		NewRequestErrorBodyData(requestError),
		requestError.ErrorCode(),
		requestError.HTTPStatus(),
	)
}
