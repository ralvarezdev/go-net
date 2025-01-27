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
func GetCtxBody(r *http.Request) interface{} {
	return r.Context().Value(CtxBodyKey)
}
