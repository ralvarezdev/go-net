package route

type (
	// ControllerWrapper is the interface for the route controller
	ControllerWrapper interface {
		RegisterRoutes()
		RegisterGroups()
	}

	// Controller is the struct for the route controller
	Controller struct {
		RouterWrapper
	}
)
