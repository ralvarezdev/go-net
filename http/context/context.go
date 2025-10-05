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
