package fault

import (
	"github.com/Southclaws/libfault"
)

var fault = libfault.Config{
	// Use default error message formatter
	FormatErrorMessage: nil,

	// Use default location getter
	GetLocation: nil,
}

type Wrapper = libfault.Wrapper

// New creates a new error with the given message and optional wrappers
func New(msg string, wrappers ...Wrapper) error {
	const callStackDelta = 1
	return fault.New(msg, callStackDelta, wrappers...)
}

// Wrap wraps an existing error with optional wrappers
func Wrap(err error, wrappers ...Wrapper) error {
	const callStackDelta = 1
	return fault.Wrap(err, callStackDelta, wrappers...)
}

// Flatten attempts to derive more useful structured information from an error chain
func Flatten(err error) libfault.Chain {
	return libfault.Flatten(err)
}
