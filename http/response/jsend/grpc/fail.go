package grpc

import (
	"net/http"
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"

	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

// NewFailDataErrorFromErrorDetailsBadRequest creates a new JSend fail data error from error details of type BadRequest
//
// Parameters:
//
//   - errorDetails: The error details
//   - parseAsValidations: Whether to parse the error details as validation errors
//
// Returns:
//
//   - *gonethttpresponse.FailData: The JSend fail data error
func NewFailDataErrorFromErrorDetailsBadRequest(
	errorDetails *errdetails.BadRequest,
	parseAsValidations bool,
) *gonethttpresponse.FailDataError {
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

		return gonethttpresponse.NewFailDataErrorWithCode(
			data,
			ErrCodeBadRequest,
			http.StatusBadRequest,
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
	return gonethttpresponse.NewFailDataErrorWithCode(
		data,
		ErrCodeBadRequest,
		http.StatusBadRequest,
	)
}

// NewFailDataErrorFromErrorDetailsPreconditionFailure creates a new JSend fail data error from error details of type
// PreconditionFailure
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - *gonethttpresponse.FailDataError: The JSend fail data error
func NewFailDataErrorFromErrorDetailsPreconditionFailure(
	errorDetails *errdetails.PreconditionFailure,
) *gonethttpresponse.FailDataError {
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
	return gonethttpresponse.NewFailDataErrorWithCode(
		data,
		ErrCodePreconditionFailure,
		http.StatusPreconditionFailed,
	)
}

// NewFailDataErrorFromErrorDetailsQuotaFailure creates a new JSend fail data error from error details of type
// QuotaFailure
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - *gonethttpresponse.FailDataError: The JSend fail data error
func NewFailDataErrorFromErrorDetailsQuotaFailure(
	errorDetails *errdetails.QuotaFailure,
) *gonethttpresponse.FailDataError {
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
	return gonethttpresponse.NewFailDataErrorWithCode(data, ErrCodeQuotaFailure, http.StatusTooManyRequests)
}

// NewFailDataErrorFromErrorDetailsRequestInfo creates a JSend fail data error from RequestInfo error details
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - *gonethttpresponse.FailDataError: The JSend fail data error
func NewFailDataErrorFromErrorDetailsRequestInfo(
	errorDetails *errdetails.RequestInfo,
) *gonethttpresponse.FailDataError {
	// Create the data map
	data := map[string]any{
		"request_id":   errorDetails.GetRequestId(),
		"serving_data": errorDetails.GetServingData(),
	}
	return gonethttpresponse.NewFailDataErrorWithCode(data, ErrCodeRequestInfo, http.StatusBadRequest)
}

// NewFailDataErrorFromErrorDetailsResourceInfo creates a JSend fail data error from ResourceInfo error details
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - *gonethttpresponse.FailDataError: The JSend fail data error
func NewFailDataErrorFromErrorDetailsResourceInfo(
	errorDetails *errdetails.ResourceInfo,
) *gonethttpresponse.FailDataError {
	// Create the data map
	data := map[string]any{
		"resource_type": errorDetails.GetResourceType(),
		"resource_name": errorDetails.GetResourceName(),
		"owner":         errorDetails.GetOwner(),
		"description":   errorDetails.GetDescription(),
	}
	return gonethttpresponse.NewFailDataErrorWithCode(data, ErrCodeResourceInfo, http.StatusBadRequest)
}

// NewFailDataErrorFromErrorDetailsHelp creates a JSend fail data error from Help error details
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - *gonethttpresponse.FailDataError: The JSend fail data error
func NewFailDataErrorFromErrorDetailsHelp(
	errorDetails *errdetails.Help,
) *gonethttpresponse.FailDataError {
	// Initialize the links slice
	links := make([]*errdetails.Help_Link, len(errorDetails.GetLinks()))

	// Loop through the links and add them to the data map
	links = append(links, errorDetails.GetLinks()...)

	// Initialize the data map
	data := map[string]any{
		"links": links,
	}
	return gonethttpresponse.NewFailDataErrorWithCode(data, ErrCodeHelp, http.StatusBadRequest)
}

// NewFailDataErrorFromErrorDetailsLocalizedMessage creates a JSend fail data error from LocalizedMessage error details
//
// Parameters:
//
//   - errorDetails: The error details
//
// Returns:
//
//   - *gonethttpresponse.FailDataError: The JSend fail data error
func NewFailDataErrorFromErrorDetailsLocalizedMessage(
	errorDetails *errdetails.LocalizedMessage,
) *gonethttpresponse.FailDataError {
	// Create the data map
	data := map[string]any{
		"locale":  errorDetails.GetLocale(),
		"message": errorDetails.GetMessage(),
	}
	return gonethttpresponse.NewFailDataErrorWithCode(
		data,
		ErrCodeLocalizedMessage,
		http.StatusBadRequest,
	)
}