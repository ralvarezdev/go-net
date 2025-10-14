package jsend

import (
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	// SuccessBody struct
	SuccessBody[T interface{}] struct {
		Status Status `json:"status"`
		Data   T      `json:"data,omitempty"`
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
func NewSuccessBody[T interface{}](
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
	data interface{},
	httpStatus int,
) gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(NewSuccessBody(data), httpStatus)
}
