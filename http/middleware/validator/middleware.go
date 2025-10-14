package validator

import (
	"log/slog"
	"net/http"
	"reflect"

	gonethttpctx "github.com/ralvarezdev/go-net/http/context"
	gonethttprequesthandler "github.com/ralvarezdev/go-net/http/request/handler"
	goreflect "github.com/ralvarezdev/go-reflect"
	govalidatormapper "github.com/ralvarezdev/go-validator/mapper"
	govalidatormapperparser "github.com/ralvarezdev/go-validator/mapper/parser"
	govalidatormapperparserjson "github.com/ralvarezdev/go-validator/mapper/parser/json"
	govalidatormappervalidator "github.com/ralvarezdev/go-validator/mapper/validator"
)

type (
	// Middleware struct is the validation middleware
	Middleware struct {
		handler          gonethttprequesthandler.Handler
		validatorService govalidatormappervalidator.Service
		generator        govalidatormapper.Generator
		validateFns      map[string]func(next http.Handler) http.Handler
		logger           *slog.Logger
	}
)

// NewMiddleware creates a new Middleware instance
//
// Parameters:
//
//   - handler: The HTTP handler to parse the request body
//   - birthdateOptions: The birthdate options (can be nil)
//   - passwordOptions: The password options (can be nil)
//   - logger: The logger (can be nil)
//
// Returns:
//
//   - *Middleware: The middleware instance
//   - error: The error if any
func NewMiddleware(
	handler gonethttprequesthandler.Handler,
	birthdateOptions *govalidatormappervalidator.BirthdateOptions,
	passwordOptions *govalidatormappervalidator.PasswordOptions,
	logger *slog.Logger,
) (*Middleware, error) {
	// Check if the handler is nil
	if handler == nil {
		return nil, gonethttprequesthandler.ErrNilHandler
	}

	// Initialize the raw parser
	rawParser := govalidatormapperparser.NewDefaultRawParser(logger)

	// Initialize the end parser
	endParser := govalidatormapperparserjson.NewDefaultEndParser()

	// Initialize the validator
	validator := govalidatormappervalidator.NewDefaultValidator(logger)

	// Initialize the validator service
	validatorService, err := govalidatormappervalidator.NewDefaultService(
		rawParser,
		endParser,
		validator,
		birthdateOptions,
		passwordOptions,
		logger,
	)
	if err != nil {
		return nil, err
	}

	// Initialize the generator
	generator := govalidatormapper.NewJSONGenerator(logger)

	if logger != nil {
		logger = logger.With(
			slog.String("component", "http_middleware_validator"),
		)
	}

	return &Middleware{
		handler:          handler,
		validatorService: validatorService,
		generator:        generator,
		logger:           logger,
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
) (*govalidatormapper.Mapper, reflect.Type, error) {
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
//   - decode: Whether to decode the body or not
//   - cache: Whether to cache the validation function or not
//   - auxiliaryValidatorFns: Optional auxiliary validator functions
//
// Returns:
//
//   - func(next http.Handler) http.Handler: the validation middleware
//   - error: if there was an error creating the validation function
func (m Middleware) CreateValidateFn(
	bodyExample interface{},
	decode bool,
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
	innerValidateFn, err := m.validatorService.CreateValidateFn(
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

				// Decode the request body if needed, and validate it
				if decode {
					if !m.handler.DecodeAndValidate(
						w,
						r,
						dest,
						innerValidateFn,
					) {
						return
					}
				} else if !m.handler.Validate(
					w,
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
		false,
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

// DecodeAndValidate decodes and validates the request body and stores it in the context
//
// Parameters:
//
//   - body: The body to decode and validate
//   - auxiliaryValidatorFns: Optional auxiliary validator functions
//
// Returns:
//
//   - func(next http.Handler) http.Handler: the validation middleware
func (m Middleware) DecodeAndValidate(
	body interface{},
	auxiliaryValidatorFns ...interface{},
) func(next http.Handler) http.Handler {
	validateFn, err := m.CreateValidateFn(
		body,
		true,
		true,
		auxiliaryValidatorFns...,
	)
	if err != nil {
		// Log the error and panic
		if m.logger != nil {
			m.logger.Error(
				"Failed to create decode and validate function",
				slog.Any("error", err),
			)
		}
		panic(err)
	}
	return validateFn
}
