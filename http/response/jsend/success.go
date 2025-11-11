package jsend

import (
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	// SuccessBody struct
	SuccessBody[T any] struct {
		// Data of the response
		Data T `json:"data,omitempty"`

		// Status of the response:
		//   - "success"
		Status Status `json:"status" enum:"success"`
	}
)

// NewSuccessBody creates a new JSend success response body
//
// Parameters:
//
//   - data: The data
//
// Returns:
//
//   - *SuccessBody: The JSend success body
func NewSuccessBody[T any](
	data T,
) *SuccessBody[T] {
	return &SuccessBody[T]{
		Status: StatusSuccess,
		Data:   data,
	}
}

// NewSuccessResponse creates a new JSend success response
//
// Parameters:
//
//   - data: The data
//   - httpStatus: The HTTP status code
//
// Returns:
//
//   - Response: The response
func NewSuccessResponse(
	data any,
	httpStatus int,
) gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(NewSuccessBody(data), httpStatus)
}
