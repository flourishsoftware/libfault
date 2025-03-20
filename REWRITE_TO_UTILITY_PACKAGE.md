

I want to rewrite this package to be a toolkit for krafting your own custom "errors" package instead of using fault package directly as a replacement of standard libarary "errors" package.

Let's rename the package from "fault" to "libfault".

The API should be something like:


```go
// Package myerrors is my custom errors package.
package myerrors

import (
	"github.com/josephbuchma/libfault"
)

var fault = libfault.Config{
	// Set custom error message formatter.
	FormatErrorMessage: func(chain fault.Chain) string {...},

	// GetLocation overrides default getLocation func.
	// Also accepts skipFramesDelta parameter to adjust which stack frame
	// should be picked if we additionally wrap the Wrap or New method call.
	GetLocation: func(skipFramesDelta int) string {...},
}

type Wrapper = libfault.Wrapper

func New(msg string, wrappers ...Wrapper) error {
	const callStackDelta = 1
	return fault.New(msg, callStackDelta, wrappers)
}

func Wrap(err errorr, wrappers ...Wrapper) error {
	const callStackDelta = 1
	return fault.Wrap(err, wrappers...)
}
```

In the example above you can see an example "myerrors" package created using the libfault package.