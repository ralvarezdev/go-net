package factory

import (
	gonethttproute "github.com/ralvarezdev/go-net/http/route"
)

type (
	// ControllerWrapper is the interface for the route controller
	ControllerWrapper interface {
		CreateRouter(baseRouter gonethttproute.RouterWrapper, path string) error
		GetRouter() gonethttproute.RouterWrapper
		RegisterRoutes()
		RegisterGroups()
		gonethttproute.RouterWrapper
	}

	// Controller is the struct for the route controller
	Controller struct {
		gonethttproute.RouterWrapper
	}
)

// CreateRouter creates the controller
func (c *Controller) CreateRouter(
	baseRouter gonethttproute.RouterWrapper,
	path string,
) error {
	// Check if the base route is nil
	if baseRouter == nil {
		return gonethttproute.ErrNilRouter
	}

	// Set the base route
	c.RouterWrapper = baseRouter.NewGroup(path)

	// Register the controller routes and groups
	c.RegisterRoutes()
	c.RegisterGroups()
	return nil
}

// GetRouter returns the router
func (c *Controller) GetRouter() gonethttproute.RouterWrapper {
	return c.RouterWrapper
}

// RegisterRoutes registers the routes
func (c *Controller) RegisterRoutes() {}

// RegisterGroups registers the groups
func (c *Controller) RegisterGroups() {}
