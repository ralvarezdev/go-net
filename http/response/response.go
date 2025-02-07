package response

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
)

type (
	// Response is the interface for the responses
	Response interface {
		Body(mode *goflagsmode.Flag) interface{}
		HTTPStatus() int
	}

	// DefaultResponse struct
	DefaultResponse struct {
		body       interface{}
		debugBody  interface{}
		httpStatus int
	}
)

// NewResponse creates a new response
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
func NewDebugResponse(
	body interface{},
	debugBody interface{},
	httpStatus int,
) *DefaultResponse {
	return &DefaultResponse{
		body,
		debugBody,
		httpStatus,
	}
}

// Body returns the response body
func (d *DefaultResponse) Body(mode *goflagsmode.Flag) interface{} {
	if mode != nil && mode.IsDebug() {
		return d.debugBody
	}
	return d.body
}

// HTTPStatus returns the HTTP status
func (d *DefaultResponse) HTTPStatus() int {
	return d.httpStatus
}
