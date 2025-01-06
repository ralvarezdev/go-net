package response

import (
	gonethttp "github.com/ralvarezdev/go-net/http"
	"net/http"
)

var (
	InternalServerError = NewErrorResponse(
		gonethttp.InternalServerError,
		nil,
		nil,
		http.StatusInternalServerError,
	)
)
