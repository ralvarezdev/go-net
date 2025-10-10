package validator

import (
	"log/slog"
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
		handler     gonethttphandler.Handler
		validator   govalidatorstructmappervalidator.Service
		generator   govalidatorstructmapper.Generator
		validateFns map[string]func(next http.Handler) http.Handler
		logger      *slog.Logger
	}
)

// NewMiddleware creates a new Middleware instance
//
// Parameters:
//
//   - handler: The HTTP handler to parse the request body
//   - validator: The struct validator service
//   - generator: The struct mapper generator
//   - logger: The logger (can be nil)
//
// Returns:
//
//   - *Middleware: The middleware instance
//   - error: The error if any
func NewMiddleware(
	handler gonethttphandler.Handler,
	validator govalidatorstructmappervalidator.Service,
	generator govalidatorstructmapper.Generator,
	logger *slog.Logger,
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

	if logger != nil {
		logger = logger.With(
			slog.String("component", "http_middleware_validator"),
		)
	}

	return &Middleware{
		handler:   handler,
		validator: validator,
		generator: generator,
		logger:    logger,
	}, nil
}

// createMapper creates a mapper for a given struct
//
// Parameters:
//
//   - structInstance: the struct instance to create the mapper for
//
// Returns:
//
//   - *govalidatormapper.Mapper: the mapper
//   - error: if there was an error creating the mapper
func (m Middleware) createMapper(
	structInstance interface{},
) (*govalidatorstructmapper.Mapper, reflect.Type, error) {
	// Get the type of the request
	structInstanceType := goreflect.GetTypeOf(structInstance)

	// Dereference the request type if it is a pointer
	if structInstanceType.Kind() == reflect.Pointer {
		structInstanceType = structInstanceType.Elem()
	} else {
		structInstance = &structInstance
	}

	// Create the mapper
	mapper, err := m.generator.NewMapper(structInstance)
	if err != nil {
		if m.logger != nil {
			m.logger.Error(
				"Failed to create mapper",
				slog.String("type", structInstanceType.String()),
				slog.Any("error", err),
			)
		}
		return nil, structInstanceType, err
	}
	return mapper, structInstanceType, nil
}

// CreateValidateFn validates the request body and stores it in the context
//
// Parameters:
//
//   - bodyExample: An example of the body to validate
//   - cache: Whether to cache the validation function or not
//   - auxiliaryValidatorFns: Optional auxiliary validator functions
//
// Returns:
//
//   - func(next http.Handler) http.Handler: the validation middleware
//   - error: if there was an error creating the validation function
func (m Middleware) CreateValidateFn(
	bodyExample interface{},
	cache bool,
	auxiliaryValidatorFns ...interface{},
) (func(next http.Handler) http.Handler, error) {
	// Create the mapper
	mapper, bodyType, err := m.createMapper(bodyExample)
	if err != nil {
		return nil, err
	}

	// Check if the validate function is already cached
	if cache && m.validateFns != nil {
		if validateFn, ok := m.validateFns[goreflect.UniqueTypeReference(mapper.GetStructInstance())]; ok {
			return validateFn, nil
		}
	}

	// Create the inner validate function
	innerValidateFn, err := m.validator.CreateValidateFn(
		mapper,
		cache,
		auxiliaryValidatorFns...,
	)
	if err != nil {
		panic(err)
	}

	// Create the validate function
	validateFn := func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Get a new instance of the body
				dest := goreflect.NewInstanceFromType(bodyType)

				// Decode the request body
				if !m.handler.Parse(
					w,
					r,
					dest,
					innerValidateFn,
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

	// Cache the validate function
	if cache {
		if m.validateFns == nil {
			m.validateFns = make(map[string]func(next http.Handler) http.Handler)
		}
		m.validateFns[goreflect.UniqueTypeReference(mapper.GetStructInstance())] = validateFn
	}

	return validateFn, nil
}

// Validate validates the request body and stores it in the context
//
// Parameters:
//
//   - body: The body to validate
//   - auxiliaryValidatorFns: Optional auxiliary validator functions
//
// Returns:
//
//   - func(next http.Handler) http.Handler: the validation middleware
func (m Middleware) Validate(
	body interface{},
	auxiliaryValidatorFns ...interface{},
) func(next http.Handler) http.Handler {
	validateFn, err := m.CreateValidateFn(
		body,
		true,
		auxiliaryValidatorFns...,
	)
	if err != nil {
		// Log the error and panic
		if m.logger != nil {
			m.logger.Error(
				"Failed to create validate function",
				slog.Any("error", err),
			)
		}
		panic(err)
	}
	return validateFn
}
