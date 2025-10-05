package response

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
)

type (
	// DefaultResponse struct
	DefaultResponse struct {
		body       interface{}
		debugBody  interface{}
		httpStatus int
	}
)

// NewResponse creates a new response
//
// Parameters:
//
//   - body: The response body
//   - httpStatus: The HTTP status
//
// Returns:
//
//   - *DefaultResponse: The default response
func NewResponse(
	body interface{},
	httpStatus int,
) *DefaultResponse {
	return &DefaultResponse{
		body,
		body,
		httpStatus,
	}
}

// NewDebugResponse creates a new debug response
//
// Parameters:
//
//   - body: The response body
//   - debugBody: The debug response body
//   - httpStatus: The HTTP status
//
// Returns:
//
//   - *DefaultResponse: The default response
func NewDebugResponse(
	body interface{},
	debugBody interface{},
	httpStatus int,
) *DefaultResponse {
	if debugBody == nil {
		debugBody = body
	}

	return &DefaultResponse{
		body,
		debugBody,
		httpStatus,
	}
}

// Body returns the response body
//
// Parameters:
//
//   - mode: The flag mode
//
// Returns:
//
//   - interface{}: The response body
func (d DefaultResponse) Body(mode *goflagsmode.Flag) interface{} {
	if mode != nil && mode.IsDebug() {
		return d.debugBody
	}
	return d.body
}

// HTTPStatus returns the HTTP status
//
// Returns:
//
//   - int: The HTTP status
func (d DefaultResponse) HTTPStatus() int {
	return d.httpStatus
}
