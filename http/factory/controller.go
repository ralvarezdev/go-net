package factory

import (
	gonethttproute "github.com/ralvarezdev/go-net/http/route"
)

type (
	// ControllerWrapper is the interface for the route controller
	ControllerWrapper interface {
		RegisterRoutes()
		RegisterGroups()
	}

	// Controller is the struct for the route controller
	Controller struct {
		gonethttproute.RouterWrapper
	}

	// RegisterRoutes registers the router routes
	func (c *Controller) RegisterRoutes() {}

	// RegisterGroups registers the route groups
	func (c *Controller) RegisterGroups() {}
)
