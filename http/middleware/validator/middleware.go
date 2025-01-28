package validator

import (
	gonethttpctx "github.com/ralvarezdev/go-net/http/context"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
	goreflect "github.com/ralvarezdev/go-reflect"
	govalidatorstructmapper "github.com/ralvarezdev/go-validator/struct/mapper"
	"net/http"
	"reflect"
)

// Middleware struct
type Middleware struct {
	handler   gonethttphandler.Handler
	generator govalidatorstructmapper.Generator
}

// NewMiddleware creates a new Middleware instance
func NewMiddleware(
	handler gonethttphandler.Handler,
	generator govalidatorstructmapper.Generator,
) (*Middleware, error) {
	// Check if the handler or the generator is nil
	if handler == nil {
		return nil, gonethttphandler.ErrNilHandler
	}
	if generator == nil {
		return nil, govalidatorstructmapper.ErrNilGenerator
	}

	return &Middleware{
		handler,
		generator,
	}, nil
}

// Validate validates the request body and stores it in the context
func (m *Middleware) Validate(
	body,
	createValidateFn interface{},
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

	// Check if the createValidateFn is valid
	if createValidateFn == nil {
		panic(ErrNilCreateValidateFn)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Get a new instance of the body
				dest := goreflect.NewInstanceFromType(bodyType)

				// Call the validate function
				results, err := goreflect.SafeCallFunction(
					createValidateFn,
					dest,
					mapper,
				)
				if err != nil {
					panic(err)
				}
				validateFn := results[0]

				// Parse the validate function
				if validateFn == nil {
					panic(ErrNilValidateFn)
				}
				parsedValidateFn, ok := validateFn.(func() (interface{}, error))
				if !ok {
					panic(ErrInvalidValidateFn)
				}

				// Decode the request body
				if !m.handler.Parse(
					w,
					r,
					dest,
					parsedValidateFn,
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
