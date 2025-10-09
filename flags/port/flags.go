package port

import (
	goflags "github.com/ralvarezdev/go-flags"
)

type (
	// Flag is a custom flag type for port, embedding goflags.Flag.
	Flag struct {
		goflags.Flag
	}
)

// NewFlag creates a new Flag with allowed values.
//
// Parameters:
//
//	defaultValue - the default value for the flag.
//
// Returns:
//
//	A pointer to the created Flag.
func NewFlag(
	defaultValue *string,
) *Flag {
	return &Flag{
		Flag: *goflags.NewFlag(defaultValue, nil, FlagName, FlagUsage),
	}
}

// Default returns the default value of the flag.
//
// Returns:
//
//	The default value.
func (f *Flag) Default() string {
	if f == nil {
		return ""
	}
	return f.Default()
}

// SetFlag initializes the port flag.
func SetFlag(flag *Flag) {
	if flag != nil {
		flag.SetFlag()
	}
}
