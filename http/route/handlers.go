package route

import (
	"net/http"
)

// ChainHandlers chains the handlers
//
// Parameters:
//
//   - lastHandler: The last handler to be executed
//   - handlers: The handlers to be chained
//
// Returns:
//
//   - http.Handler: The chained handler
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
