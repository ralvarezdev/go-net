package route

type (
	// Controller is the struct for the controller
	Controller struct {
		service Service
		GroupWrapper
	}
)
