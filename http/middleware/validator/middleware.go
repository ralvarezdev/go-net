package validator

import (
	"net/http"
	"reflect"

	gonethttpctx "github.com/ralvarezdev/go-net/http/context"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
	goreflect "github.com/ralvarezdev/go-reflect"
	govalidatorstructmapper "github.com/ralvarezdev/go-validator/struct/mapper"
	govalidatorstructmappervalidator "github.com/ralvarezdev/go-validator/struct/mapper/validator"
)

type (
	// Middleware struct is the validation middleware
	Middleware struct {
		handler   gonethttphandler.Handler
		validator govalidatorstructmappervalidator.Service
		generator govalidatorstructmapper.Generator
	}
)

// NewMiddleware creates a new Middleware instance
//
// Parameters:
//
//   - handler: The HTTP handler to parse the request body
//   - validator: The struct validator service
//   - generator: The struct mapper generator
//
// Returns:
//
//   - *Middleware: The middleware instance
//   - error: The error if any
func NewMiddleware(
	handler gonethttphandler.Handler,
	validator govalidatorstructmappervalidator.Service,
	generator govalidatorstructmapper.Generator,
) (*Middleware, error) {
	// Check if the handler, validator or the generator is nil
	if handler == nil {
		return nil, gonethttphandler.ErrNilHandler
	}
	if validator == nil {
		return nil, govalidatorstructmappervalidator.ErrNilService
	}
	if generator == nil {
		return nil, govalidatorstructmapper.ErrNilGenerator
	}

	return &Middleware{
		handler,
		validator,
		generator,
	}, nil
}

// Validate validates the request body and stores it in the context
//
// Parameters:
//
//   - body: An instance of the body to validate
//   - auxiliaryValidatorFns: Optional auxiliary validator functions
//
// Returns:
//
//   - func(next http.Handler) http.Handler: The middleware function
func (m Middleware) Validate(
	body interface{},
	auxiliaryValidatorFns ...govalidatorstructmappervalidator.AuxiliaryValidatorFn,
) func(next http.Handler) http.Handler {
	// Get the type of the body
	bodyType := goreflect.GetTypeOf(body)

	// Dereference the body type if it is a pointer
	if bodyType.Kind() == reflect.Pointer {
		bodyType = bodyType.Elem()
	} else {
		body = &body
	}

	// Create the mapper
	mapper, err := m.generator.NewMapper(body)
	if err != nil {
		panic(err)
	}

	// Create the validate function
	validateFn, err := m.validator.CreateValidateFn(
		mapper,
		auxiliaryValidatorFns...,
	)
	if err != nil {
		panic(err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Get a new instance of the body
				dest := goreflect.NewInstanceFromType(bodyType)

				// Decode the request body
				if !m.handler.Parse(
					w,
					r,
					dest,
					validateFn,
				) {
					return
				}

				// Store the validated body in the context
				r = gonethttpctx.SetCtxBody(r, dest)

				// Call the next handler
				next.ServeHTTP(w, r)
			},
		)
	}
}
