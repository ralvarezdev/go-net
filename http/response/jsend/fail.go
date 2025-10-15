package jsend

import (
	"net/http"
	"strings"

	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
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

// NewFailBodyFromErrorDetailsBadRequest creates a new JSend fail response body from error details of type BadRequest
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - *FailBody: The JSend fail body
func NewFailBodyFromErrorDetailsBadRequest(
	errorDetails *errdetails.BadRequest,
) *FailBody {
	// Initialize the data map
	data := make(map[string]interface{})

	// Loop through the field violations and add them to the data map
	for _, violation := range errorDetails.GetFieldViolations() {
		// Split the field by dot notation
		parts := strings.Split(violation.GetField(), ".")

		// If there are no parts, add the violation to the data map
		if len(parts) == 0 {
			// If the field doesn't exist, create a new slice
			if _, ok := data[violation.GetField()]; !ok {
				data[violation.GetField()] = []string{}
			}

			// Parse the data map to a slice
			parsedSlice, ok := data[violation.GetField()].([]string)
			if !ok {
				panic(
					gonethttpresponse.NewDebugErrorWithCode(
						ErrExpectedStringSliceOnDataMap,
						gonethttp.ErrInternalServerError,
						ErrCodeExpectedStringSliceOnDataMap,
						http.StatusInternalServerError,
					),
				)
			}

			// Add the violation description to the slice
			data[violation.GetField()] = append(
				parsedSlice,
				violation.GetDescription(),
			)
			continue
		}

		// If the field has dot notation, create a nested map
		nestedMap := data
		for i, part := range parts {
			if i == len(parts)-1 {
				// If the part is the last part, create a new slice if it doesn't exist
				if _, ok := nestedMap[part]; !ok {
					nestedMap[part] = []string{}
				}

				// Parse the nested map to a slice
				parsedSlice, ok := nestedMap[part].([]string)
				if !ok {
					panic(
						gonethttpresponse.NewDebugErrorWithCode(
							ErrExpectedStringSliceOnNestedDataMap,
							gonethttp.ErrInternalServerError,
							ErrCodeExpectedStringSliceOnNestedDataMap,
							http.StatusInternalServerError,
						),
					)
				}

				// If it's the last part, add the violation description
				nestedMap[part] = append(
					parsedSlice,
					violation.GetDescription(),
				)
				continue
			}

			// If the part doesn't exist, create a new map
			if _, ok := nestedMap[part]; !ok {
				nestedMap[part] = make(map[string]interface{})
			}

			// Parse the nested map
			parsedNestedMap, ok := nestedMap[part].(map[string]interface{})
			if !ok {
				panic(
					gonethttpresponse.NewDebugErrorWithCode(
						ErrExpectedMapOnNestedDataMap,
						gonethttp.ErrInternalServerError,
						ErrCodeExpectedMapOnNestedDataMap,
						http.StatusInternalServerError,
					),
				)
			}

			// Move to the next nested map
			nestedMap = parsedNestedMap
		}
	}
	return NewFailBody(data)
}

// NewFailResponseFromErrorDetailsBadRequest creates a new JSend fail response from error details of type BadRequest
//
// Parameters:
//
//   - errorDetails: The error details
//   - httpStatus: The HTTP status code
//
// Returns:
//
//   - Response: The response
func NewFailResponseFromErrorDetailsBadRequest(
	errorDetails *errdetails.BadRequest,
	httpStatus int,
) gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(
		NewFailBodyFromErrorDetailsBadRequest(errorDetails),
		httpStatus,
	)
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
