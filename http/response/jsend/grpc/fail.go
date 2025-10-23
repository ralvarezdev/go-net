package grpc

import (
	"net/http"
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"

	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpresponsejsend "github.com/ralvarezdev/go-net/http/response/jsend"
)

// NewFailBodyFromErrorDetailsBadRequest creates a new JSend fail response body from error details of type BadRequest
//
// Parameters:
//
//   - errorDetails: The error details
//   - parseAsValidations: Whether to parse the error details as validation errors
//
// Returns:
//
//   - *gonethttpresponsejsend.FailBody: The JSend fail body
func NewFailBodyFromErrorDetailsBadRequest(
	errorDetails *errdetails.BadRequest,
	parseAsValidations bool,
) *gonethttpresponsejsend.FailBody {
	// Initialize the data map
	data := make(map[string]any)

	// Check if we should parse as validation errors
	if !parseAsValidations {
		// Loop through the field violations and add them to the data map
		for _, violation := range errorDetails.GetFieldViolations() {
			// If the field doesn't exist, create a new slice
			violationField := violation.GetField()
			if _, ok := data[violationField]; !ok {
				data[violationField] = []*errdetails.BadRequest_FieldViolation{}
			}

			// Parse the data map to a slice
			parsedSlice, ok := data[violationField].([]*errdetails.BadRequest_FieldViolation)
			if !ok {
				panic(
					gonethttpresponse.NewDebugErrorWithCode(
						ErrExpectedFieldViolationSliceOnDataMap,
						gonethttp.ErrInternalServerError,
						ErrCodeExpectedFieldViolationSliceOnDataMap,
						http.StatusInternalServerError,
					),
				)
			}

			// Add the violation description to the slice
			data[violationField] = append(
				parsedSlice,
				violation,
			)
		}

		return gonethttpresponsejsend.NewFailBodyWithCode(
			data,
			ErrCodeBadRequest,
		)
	}

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
				nestedMap[part] = make(map[string]any)
			}

			// Parse the nested map
			parsedNestedMap, ok := nestedMap[part].(map[string]any)
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
	return gonethttpresponsejsend.NewFailBodyWithCode(
		data,
		ErrCodeBadRequest,
	)
}

// NewFailResponseFromErrorDetailsBadRequest creates a new JSend fail response from error details of type BadRequest
//
// Parameters:
//
//   - errorDetails: The error details
//   - parseAsValidations: Whether to parse the error details as validation errors
//
// Returns:
//
//   - Response: The response
func NewFailResponseFromErrorDetailsBadRequest(
	errorDetails *errdetails.BadRequest,
	parseAsValidations bool,
) gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(
		NewFailBodyFromErrorDetailsBadRequest(errorDetails, parseAsValidations),
		http.StatusBadRequest,
	)
}

// NewFailErrorFromErrorDetailsBadRequest creates a new fail error from error details of type BadRequest
//
// Parameters:
//
//   - errorDetails: The error details
//   - parseAsValidations: Whether to parse the error details as validation errors
//
// Returns:
//
//   - error: The error
func NewFailErrorFromErrorDetailsBadRequest(
	errorDetails *errdetails.BadRequest,
	parseAsValidations bool,
) error {
	return gonethttpresponse.NewFailDataError(
		NewFailBodyFromErrorDetailsBadRequest(errorDetails, parseAsValidations),
		http.StatusBadRequest,
	)
}

// NewFailBodyFromErrorDetailsPreconditionFailure creates a new JSend fail response body from error details of type
// PreconditionFailure
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - *gonethttpresponsejsend.FailBody: The JSend fail body
func NewFailBodyFromErrorDetailsPreconditionFailure(
	errorDetails *errdetails.PreconditionFailure,
) *gonethttpresponsejsend.FailBody {
	// Initialize the data map
	data := make(map[string]any)

	// Loop through the violations and add them to the data map
	for _, violation := range errorDetails.GetViolations() {
		// If the type doesn't exist, create a new slice
		key := violation.GetType()
		if _, ok := data[key]; !ok {
			data[key] = []*errdetails.PreconditionFailure_Violation{}
		}

		// Parse the data map to a slice
		parsedSlice, ok := data[key].([]*errdetails.PreconditionFailure_Violation)
		if !ok {
			panic(
				gonethttpresponse.NewDebugErrorWithCode(
					ErrExpectedPreconditionFailureViolationSliceOnDataMap,
					gonethttp.ErrInternalServerError,
					ErrCodeExpectedPreconditionFailureViolationSliceOnDataMap,
					http.StatusInternalServerError,
				),
			)
		}

		// Add the violation description to the slice
		data[key] = append(
			parsedSlice,
			violation,
		)
	}
	return gonethttpresponsejsend.NewFailBodyWithCode(
		data,
		ErrCodePreconditionFailure,
	)
}

// NewFailResponseFromErrorDetailsPreconditionFailure creates a new JSend fail response from error details of type
// PreconditionFailure
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - Response: The response
func NewFailResponseFromErrorDetailsPreconditionFailure(
	errorDetails *errdetails.PreconditionFailure,
) gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(
		NewFailBodyFromErrorDetailsPreconditionFailure(errorDetails),
		http.StatusPreconditionFailed,
	)
}

// NewFailErrorFromErrorDetailsPreconditionFailure creates a new fail error from error details of type
// PreconditionFailure
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - error: The error
func NewFailErrorFromErrorDetailsPreconditionFailure(
	errorDetails *errdetails.PreconditionFailure,
) error {
	return gonethttpresponse.NewFailDataError(
		NewFailBodyFromErrorDetailsPreconditionFailure(errorDetails),
		http.StatusPreconditionFailed,
	)
}

// NewFailBodyFromErrorDetailsQuotaFailure creates a new JSend fail response body from error details of type
// QuotaFailure
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - *gonethttpresponsejsend.FailBody: The JSend fail body
func NewFailBodyFromErrorDetailsQuotaFailure(
	errorDetails *errdetails.QuotaFailure,
) *gonethttpresponsejsend.FailBody {
	// Initialize the data map
	data := make(map[string]any)

	// Loop through the violations and add them to the data map
	for _, violation := range errorDetails.GetViolations() {
		// If the subject doesn't exist, use "quota" as the key
		key := violation.GetSubject()
		if key == "" {
			key = "quota"
		}
		if _, ok := data[key]; !ok {
			data[key] = []*errdetails.QuotaFailure_Violation{}
		}

		// Parse the data map to a slice
		parsedSlice, err := data[key].([]*errdetails.QuotaFailure_Violation)
		if !err {
			panic(
				gonethttpresponse.NewDebugErrorWithCode(
					ErrExpectedQuotaFailureViolationSliceOnDataMap,
					gonethttp.ErrInternalServerError,
					ErrCodeExpectedQuotaFailureViolationSliceOnDataMap,
					http.StatusInternalServerError,
				),
			)
		}

		// Add the violation description to the slice
		data[key] = append(
			parsedSlice,
			violation,
		)
	}
	return gonethttpresponsejsend.NewFailBodyWithCode(data, ErrCodeQuotaFailure)
}

// NewFailResponseFromErrorDetailsQuotaFailure creates a new JSend fail response from error details of type QuotaFailure
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - Response: The response
func NewFailResponseFromErrorDetailsQuotaFailure(
	errorDetails *errdetails.QuotaFailure,
) gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(
		NewFailBodyFromErrorDetailsQuotaFailure(errorDetails),
		http.StatusTooManyRequests,
	)
}

// NewFailErrorFromErrorDetailsQuotaFailure creates a new fail error from error details of type QuotaFailure
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - error: The error
func NewFailErrorFromErrorDetailsQuotaFailure(
	errorDetails *errdetails.QuotaFailure,
) error {
	return gonethttpresponse.NewFailDataError(
		NewFailBodyFromErrorDetailsQuotaFailure(errorDetails),
		http.StatusTooManyRequests,
	)
}

// NewFailBodyFromErrorDetailsRequestInfo creates a JSend fail body from RequestInfo error details
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - *gonethttpresponsejsend.FailBody: The JSend fail body
func NewFailBodyFromErrorDetailsRequestInfo(
	errorDetails *errdetails.RequestInfo,
) *gonethttpresponsejsend.FailBody {
	// Create the data map
	data := map[string]any{
		"request_id":   errorDetails.GetRequestId(),
		"serving_data": errorDetails.GetServingData(),
	}
	return gonethttpresponsejsend.NewFailBodyWithCode(data, ErrCodeRequestInfo)
}

// NewFailResponseFromErrorDetailsRequestInfo creates a JSend fail response from RequestInfo error details
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - Response: The response
func NewFailResponseFromErrorDetailsRequestInfo(
	errorDetails *errdetails.RequestInfo,
) gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(
		NewFailBodyFromErrorDetailsRequestInfo(errorDetails),
		http.StatusBadRequest,
	)
}

// NewFailErrorFromErrorDetailsRequestInfo creates a fail error from RequestInfo error details
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - error: The error
func NewFailErrorFromErrorDetailsRequestInfo(
	errorDetails *errdetails.RequestInfo,
) error {
	return gonethttpresponse.NewFailDataError(
		NewFailBodyFromErrorDetailsRequestInfo(errorDetails),
		http.StatusBadRequest,
	)
}

// NewFailBodyFromErrorDetailsResourceInfo creates a JSend fail body from ResourceInfo error details
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - *gonethttpresponsejsend.FailBody: The JSend fail body
func NewFailBodyFromErrorDetailsResourceInfo(
	errorDetails *errdetails.ResourceInfo,
) *gonethttpresponsejsend.FailBody {
	// Create the data map
	data := map[string]any{
		"resource_type": errorDetails.GetResourceType(),
		"resource_name": errorDetails.GetResourceName(),
		"owner":         errorDetails.GetOwner(),
		"description":   errorDetails.GetDescription(),
	}
	return gonethttpresponsejsend.NewFailBodyWithCode(data, ErrCodeResourceInfo)
}

// NewFailResponseFromErrorDetailsResourceInfo creates a JSend fail response from ResourceInfo error details
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - Response: The response
func NewFailResponseFromErrorDetailsResourceInfo(
	errorDetails *errdetails.ResourceInfo,
) gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(
		NewFailBodyFromErrorDetailsResourceInfo(errorDetails),
		http.StatusNotFound,
	)
}

// NewFailErrorFromErrorDetailsResourceInfo creates a fail error from ResourceInfo error details
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - error: The error
func NewFailErrorFromErrorDetailsResourceInfo(
	errorDetails *errdetails.ResourceInfo,
) error {
	return gonethttpresponse.NewFailDataError(
		NewFailBodyFromErrorDetailsResourceInfo(errorDetails),
		http.StatusNotFound,
	)
}

// NewFailBodyFromErrorDetailsHelp creates a JSend fail body from Help error details
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - *gonethttpresponsejsend.FailBody: The JSend fail body
func NewFailBodyFromErrorDetailsHelp(
	errorDetails *errdetails.Help,
) *gonethttpresponsejsend.FailBody {
	// Initialize the links slice
	links := make([]*errdetails.Help_Link, len(errorDetails.GetLinks()))

	// Loop through the links and add them to the data map
	links = append(links, errorDetails.GetLinks()...)

	// Initialize the data map
	data := map[string]any{
		"links": links,
	}
	return gonethttpresponsejsend.NewFailBodyWithCode(data, ErrCodeHelp)
}

// NewFailResponseFromErrorDetailsHelp creates a JSend fail response from Help error details
//
// Parameters:
//
// - errorDetails: The error details
//
// Returns:
//
//   - Response: The response
func NewFailResponseFromErrorDetailsHelp(
	errorDetails *errdetails.Help,
) gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(
		NewFailBodyFromErrorDetailsHelp(errorDetails),
		http.StatusBadRequest,
	)
}

// NewFailErrorFromErrorDetailsHelp creates a fail error from Help error details
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - error: The error
func NewFailErrorFromErrorDetailsHelp(
	errorDetails *errdetails.Help,
) error {
	return gonethttpresponse.NewFailDataError(
		NewFailBodyFromErrorDetailsHelp(errorDetails),
		http.StatusBadRequest,
	)
}

// NewFailBodyFromErrorDetailsLocalizedMessage creates a JSend fail body from LocalizedMessage error details
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - *gonethttpresponsejsend.FailBody: The JSend fail body
func NewFailBodyFromErrorDetailsLocalizedMessage(
	errorDetails *errdetails.LocalizedMessage,
) *gonethttpresponsejsend.FailBody {
	// Create the data map
	data := map[string]any{
		"locale":  errorDetails.GetLocale(),
		"message": errorDetails.GetMessage(),
	}
	return gonethttpresponsejsend.NewFailBodyWithCode(
		data,
		ErrCodeLocalizedMessage,
	)
}

// NewFailResponseFromErrorDetailsLocalizedMessage creates a JSend fail response from LocalizedMessage error details
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - Response: The response
func NewFailResponseFromErrorDetailsLocalizedMessage(
	errorDetails *errdetails.LocalizedMessage,
) gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(
		NewFailBodyFromErrorDetailsLocalizedMessage(errorDetails),
		http.StatusBadRequest,
	)
}

// NewFailErrorFromErrorDetailsLocalizedMessage creates a fail error from LocalizedMessage error details
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - error: The error
func NewFailErrorFromErrorDetailsLocalizedMessage(
	errorDetails *errdetails.LocalizedMessage,
) error {
	return gonethttpresponse.NewFailDataError(
		NewFailBodyFromErrorDetailsLocalizedMessage(errorDetails),
		http.StatusBadRequest,
	)
}
