package context

import (
	"context"
	"net/http"
)

// SetCtxBody sets the body in the context
//
// Parameters:
//
//   - r: The HTTP request
//   - body: The body to set in the context
//
// Returns:
//
//   - *http.Request: The HTTP request with the body set in the context
func SetCtxBody(r *http.Request, body interface{}) *http.Request {
	ctx := context.WithValue(r.Context(), CtxBodyKey, body)
	return r.WithContext(ctx)
}

// SetBody wraps SetCtxBody
func SetBody(r *http.Request, body interface{}) *http.Request {
	return SetCtxBody(r, body)
}

// GetCtxBody tries to get the body from the context
//
// Parameters:
//
//   - r: The HTTP request
//
// Returns:
//
//   - interface{}: The body from the context, or nil if not found
func GetCtxBody(r *http.Request) interface{} {
	return r.Context().Value(CtxBodyKey)
}

// GetBody wraps GetCtxBody
func GetBody(r *http.Request) interface{} {
	return GetCtxBody(r)
}

// GetCtxWildcards tries to get the wildcards from the context
//
// Parameters:
//
//   - r: The HTTP request
//
// Returns:
//
//   - map[string]string: The wildcards from the context, or nil if not found
func GetCtxWildcards(r *http.Request) map[string]string {
	wildcard, ok := r.Context().Value(CtxWildcardsKey).(map[string]string)
	if !ok {
		return nil
	}
	return wildcard
}

// GetWildcards wraps GetCtxWildcards
func GetWildcards(r *http.Request) map[string]string {
	return GetCtxWildcards(r)
}

// SetCtxWildcards sets the wildcards in the context
//
// Parameters:
//
//   - r: The HTTP request
//   - wildcards: The wildcards to set in the context
//
// Returns:
//
//   - *http.Request: The HTTP request with the wildcards set in the context
func SetCtxWildcards(
	r *http.Request,
	wildcards map[string]string,
) *http.Request {
	ctx := context.WithValue(r.Context(), CtxWildcardsKey, wildcards)
	return r.WithContext(ctx)
}

// SetWildcards wraps SetCtxWildcards
func SetWildcards(
	r *http.Request,
	wildcards map[string]string,
) *http.Request {
	return SetCtxWildcards(r, wildcards)
}

// SetCtxQueryParameters sets the query parameters in the context
//
// Parameters:
//
//   - r: The HTTP request
//   - parameters: The query parameters to set in the context
//
// Returns:
//
//   - *http.Request: The HTTP request with the query parameters set in the context
func SetCtxQueryParameters(
	r *http.Request,
	parameters map[string][]string,
) *http.Request {
	ctx := context.WithValue(r.Context(), CtxQueryParametersKey, parameters)
	return r.WithContext(ctx)
}

// SetQueryParameters wraps SetCtxQueryParameters
func SetQueryParameters(
	r *http.Request,
	parameters map[string][]string,
) *http.Request {
	return SetCtxQueryParameters(r, parameters)
}

// GetCtxQueryParameters tries to get the query parameters from the context
//
// Parameters:
//
//   - r: The HTTP request
//
// Returns:
//
//   - map[string][]string: The query parameters from the context, or nil if not found
func GetCtxQueryParameters(r *http.Request) map[string][]string {
	parameters, ok := r.Context().Value(CtxQueryParametersKey).(map[string][]string)
	if !ok {
		return nil
	}
	return parameters
}

// GetQueryParameters wraps GetCtxQueryParameters
func GetQueryParameters(r *http.Request) map[string][]string {
	return GetCtxQueryParameters(r)
}
