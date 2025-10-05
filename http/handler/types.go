package handler

type (
	// ToTypeFn type for the wildcard to type function
	ToTypeFn func(wildcard string, dest interface{}) error
)
