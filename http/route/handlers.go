package route

import (
	"net/http"
)

// ChainHandlers chains the handlers
func ChainHandlers(
	lastHandler http.Handler,
	handlers ...func(http.Handler) http.Handler,
) http.Handler {
	// Check if the handlers are empty
	if len(handlers) == 0 {
		return lastHandler
	}

	// Set the chained handlers
	n := len(handlers)
	chainedHandler := lastHandler
	for i := n - 1; i >= 0; i = i - 1 {
		chainedHandler = handlers[i](chainedHandler)
	}
	return chainedHandler
}
