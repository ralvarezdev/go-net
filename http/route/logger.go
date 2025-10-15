package route

import (
	"log/slog"
)

// AddRouter registers a route
//
// Parameters:
//
//   - fullPath: The full path of the route
//   - pattern: The pattern of the route
//   - logger: The logger
func AddRouter(fullPath, pattern string, logger *slog.Logger) {
	if logger != nil {
		logger.Debug(
			"Registering route",
			slog.String("full_path", fullPath),
			slog.String("pattern", pattern),
		)
	}
}
