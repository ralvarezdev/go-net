package response

type (
	// FailRequestError struct
	FailRequestError interface {
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

	// CookieError struct
	CookieError struct {
		name       string
		err        string
		httpStatus int
		errorCode  *string
	}
)

// NewFieldError creates a new field error
func NewFieldError(
	field, err string, httpStatus int, errorCode *string,
) *FieldError {
	return &FieldError{
		field,
		err,
		httpStatus,
		errorCode,
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
	header, err string, httpStatus int, errorCode *string,
) *HeaderError {
	return &HeaderError{
		header,
		err,
		httpStatus,
		errorCode,
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

// NewCookieError creates a new cookie error
func NewCookieError(
	name, err string, httpStatus int, errorCode *string,
) *CookieError {
	return &CookieError{
		name,
		err,
		httpStatus,
		errorCode,
	}
}

// Name returns the cookie name
func (c *CookieError) Name() string {
	return c.name
}

// Key returns the cookie name
func (c *CookieError) Key() string {
	return c.Name()
}

// Error returns the cookie error as a string
func (c *CookieError) Error() string {
	return c.err
}

// HTTPStatus returns the HTTP status
func (c *CookieError) HTTPStatus() int {
	return c.httpStatus
}

// ErrorCode returns the error code
func (c *CookieError) ErrorCode() *string {
	return c.errorCode
}

// NewRequestErrorsBodyData creates a new request errors body data
func NewRequestErrorsBodyData(
	requestErrors ...FailRequestError,
) *map[string]*[]string {
	// Initialize the request errors map
	requestErrorsMap := make(map[string]*[]string)

	// Add the request error to the request errors map
	for _, requestError := range requestErrors {
		requestErrorsMap[requestError.Key()] = &[]string{requestError.Error()}
	}

	return &requestErrorsMap
}

// NewResponseFromFailRequestError creates a new fail response from a fail request error
func NewResponseFromFailRequestError(
	requestError FailRequestError,
) *FailResponse {
	return NewJSendFailResponse(
		NewRequestErrorsBodyData(requestError),
		requestError.ErrorCode(),
		requestError.HTTPStatus(),
	)
}
