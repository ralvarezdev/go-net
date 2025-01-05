package route

import (
	gologgermode "github.com/ralvarezdev/go-logger/mode"
	gologgermodenamed "github.com/ralvarezdev/go-logger/mode/named"
)

// Logger is the logger for router
type Logger struct {
	logger gologgermodenamed.Logger
}

// NewLogger is the logger for the router
func NewLogger(header string, modeLogger gologgermode.Logger) (*Logger, error) {
	// Initialize the mode named logger
	namedLogger, err := gologgermodenamed.NewDefaultLogger(header, modeLogger)
	if err != nil {
		return nil, err
	}

	return &Logger{logger: namedLogger}, nil
}

// RegisterRouteGroup registers a route group
func (l *Logger) RegisterRouteGroup(routerPath string, routerGroupPath string) {
	l.logger.Debug(
		"registering route group",
		"router path: "+routerPath,
		"router group path: "+routerGroupPath,
	)
}

// RegisterRoute registers a route
func (l *Logger) RegisterRoute(routerPath string, routePath string) {
	l.logger.Debug(
		"registering route",
		"router path: "+routerPath,
		"route path: "+routePath,
	)
}
