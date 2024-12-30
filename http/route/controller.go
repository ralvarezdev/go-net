package route

type (
	// ControllerWrapper is the interface for the route controller
	ControllerWrapper interface {
		RegisterRoutes()
		RegisterRouteGroups()
	}

	// Controller is the struct for the route controller
	Controller struct {
		Service Service
		RouterWrapper
	}
)
