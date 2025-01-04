package route

import (
	gologger "github.com/ralvarezdev/go-logger"
	gologgerstatus "github.com/ralvarezdev/go-logger/status"
)

// Logger is the logger for router
type Logger struct {
	logger gologger.Logger
}

// NewLogger is the logger for the router
func NewLogger(logger gologger.Logger) (*Logger, error) {
	// Check if the logger is nil
	if logger == nil {
		return nil, gologger.ErrNilLogger
	}

	return &Logger{logger: logger}, nil
}

// RegisterRouteGroup registers a route group
func (l *Logger) RegisterRouteGroup(routerPath string, routerGroupPath string) {
	l.logger.LogMessage(
		gologger.NewLogMessage(
			"Registering route group",
			gologgerstatus.Debug,
			"router path: "+routerPath,
			"router group path: "+routerGroupPath,
		),
	)
}

// RegisterRoute registers a route
func (l *Logger) RegisterRoute(routerPath string, routePath string) {
	l.logger.LogMessage(
		gologger.NewLogMessage(
			"Registering route",
			gologgerstatus.Debug,
			"router path: "+routerPath,
			"route path: "+routePath,
		),
	)
}
