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

/*
// NewFailBodyFromErrorDetails creates a new JSend fail response body from error details
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - *FailBody: The JSend fail body
func NewFailBodyFromErrorDetails(
	errorDetails *errdetails.BadRequest,
) *FailBody {
	// Initialize the data map
	data := make(map[string]interface{})

	// Loop through the field violations and add them to the data map
	for _, violation := range errorDetails.GetFieldViolations() {
		// Split the field by dot notation
		parts := strings.Split(violation.GetField(), ".")
		if len(parts) > 1 {
			// If the field has dot notation, create a nested map
			nestedMap := data
			for i, part := range parts {
				if i == len(parts)-1 {
					// If the part is the last part, create a new slice if it doesn't exist
					if _, ok := nestedMap[part]; !ok {
						nestedMap[part] = []string{}
					}

					// If it's the last part, add the violation description
					nestedMap[part] = append(
						nestedMap[part].([]string),
						violation.GetDescription(),
					)
				} else {
					// If the part doesn't exist, create a new map
					if _, ok := nestedMap[part]; !ok {
						nestedMap[part] = make(map[string]interface{})
					}

					// Move to the next nested map
					parsedNestedMap, ok := nestedMap[part].(map[string]interface{})
					if !ok {
						panic(
							NewFailError()
							)
					}
				}
			}
			continue
		}

		// If the field doesn't exist, create a new slice
		if _, ok := data[violation.GetField()]; !ok {
			data[violation.GetField()] = []string{}
		}

		// Add the violation description to the slice
		data[violation.GetField()] = append(
			data[violation.GetField()].([]string),
			violation.GetDescription(),
		)
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
*/

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
