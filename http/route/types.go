package route

import (
	"net/http"
)

type (
	// EndpointHandler is the type for the HTTP endpoint handler
	EndpointHandler func(w http.ResponseWriter, r *http.Request) error
)
