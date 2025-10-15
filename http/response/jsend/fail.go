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
