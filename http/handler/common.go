package handler

import (
	gonethttp "github.com/ralvarezdev/go-net/http"
	"net/http"
)

// SendInternalServerError sends an internal server error response
func SendInternalServerError(w http.ResponseWriter) {
	http.Error(
		w,
		gonethttp.InternalServerError,
		http.StatusInternalServerError,
	)
}
