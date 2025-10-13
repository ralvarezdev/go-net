package route

import (
	"net/http"
)

// AddHandlersToStart adds a handler to the start of the handlers slice
//
// Parameters:
//
//   - originalHandlers: The original handlers slice
//   - handlersToAdd: The handlers to be added
func AddHandlersToStart(
	originalHandlers *[]func(http.Handler) http.Handler,
	handlersToAdd ...func(http.Handler) http.Handler,
) {
	// Check if the handlers slice is nil or the handlers to add are empty
	if originalHandlers == nil {
		return
	}
	if len(handlersToAdd) == 0 {
		return
	}

	// Add the handler to the start of the handlers slice
	*originalHandlers = append(handlersToAdd, *originalHandlers...)
}

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
