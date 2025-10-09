package port

import (
	goflags "github.com/ralvarezdev/go-flags"
	gostringsconvert "github.com/ralvarezdev/go-strings/convert"
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

// Port returns the current port value.
//
// Returns:
//
//	The port number as an integer, or an error if conversion fails.
func (f *Flag) Port() (int, error) {
	if f == nil {
		return 0, nil
	}

	// Convert the flag value to an integer
	var port int
	if err := gostringsconvert.ToInt(f.Value(), &port); err != nil {
		return 0, err
	}
	return port, nil
}

// SetFlag initializes the port flag.
func SetFlag(flag *Flag) {
	if flag != nil {
		flag.SetFlag()
	}
}
