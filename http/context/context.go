package context

import (
	"context"
	"net/http"
)

// SetCtxBody sets the body in the context
func SetCtxBody(r *http.Request, body interface{}) *http.Request {
	ctx := context.WithValue(r.Context(), CtxBodyKey, body)
	return r.WithContext(ctx)
}

// GetCtxBody tries to get the body from the context
func GetCtxBody(r *http.Request) (interface{}, error) {
	// Get the token claims from the context
	value := r.Context().Value(CtxBodyKey)
	if value == nil {
		return nil, ErrMissingBodyInContext
	}
	return value, nil
}
