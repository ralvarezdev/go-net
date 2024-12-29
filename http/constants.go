package http

import (
	"net/http"
)

var (
	InternalServerError = http.StatusText(http.StatusInternalServerError)
	BadRequest          = http.StatusText(http.StatusBadRequest)
)
