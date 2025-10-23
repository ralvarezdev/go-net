package response

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
)

type (
	// DefaultResponse struct
	DefaultResponse struct {
		body       any
		debugBody  any
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
	body any,
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
	body any,
	debugBody any,
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
//   - any: The response body
func (d DefaultResponse) Body(mode *goflagsmode.Flag) any {
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
