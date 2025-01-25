package factory

import (
	gonethttproute "github.com/ralvarezdev/go-net/http/route"
)

type (
	// ControllerWrapper is the interface for the route controller
	ControllerWrapper interface {
		RegisterRoutes()
		RegisterGroups()
		gonethttproute.RouterWrapper
	}
)
