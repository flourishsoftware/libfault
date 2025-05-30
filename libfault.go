// Package libfault provides an extensible yet ergonomic mechanism for wrapping
// errors. It implements this as a kind of middleware style pattern by providing
// a simple option-style interface that can be passed to a call to `fault.Wrap`.
//
// See the GitHub repository for full documentation and examples.
package libfault

import (
	"fmt"
	"runtime"
	"strings"
)

type ChainMessageDeduplicationMode uint8

const (
	// ChainMessageDeduplicationModeExactMatch is the default mode.
	// Messages are de-duplicated on exact match only.
	ChainMessageDeduplicationModeExactMatch ChainMessageDeduplicationMode = iota

	// ChainMessageDeduplicationModeSubstringMatch will de-duplicate shared prefix
	// in subsequent error messages in the chain.
	ChainMessageDeduplicationModeSubstringMatch
)

// Fault allows customization of error formatting, location determination etc.
type Fault struct {
	// BuildDefaultErrorMessage allows customizing how error messages are formatted by default.
	BuildDefaultErrorMessage func(chain Chain) string

	// FormatErrorMessage allows customizing how errors are formatted when using fmt package, such as fmt.Sprintf("%+v", err)
	FormatErrorMessage func(chain Chain, s fmt.State, verb rune)

	// GetLocation overrides the default getLocation function.
	GetLocation func(skipFramesDelta int) string

	// ChainDeduplicationMode allows to select a strategy for messages deduplication in the Chain
	// when building a final error message.
	ChainDeduplicationMode ChainMessageDeduplicationMode

	// AllowWrapperToDiscardError allows to "nillify" an error by a Wrapper that returns nil.
	// By default we do not discarn an error in such situation.
	AllowWrapperToDiscardError bool
}

// Wrapper describes a kind of middleware that packages can satisfy in order to
// decorate errors with additional, domain-specific context.
type Wrapper func(err error) error

// Wrap wraps an error with all of the wrappers provided.
func (c *Fault) Wrap(err error, skipFramesDelta int, w ...Wrapper) error {
	if err == nil {
		return nil
	}

	// The passed err might already have a location if it's a 'fault.New' error, or it might not if it's another type of
	// error like one from the standard library. Wrapping it in a container with an empty location ensures that the
	// location will be reset when we flatten the error chain. If the error is a 'fault.New' error, it will itself be
	// wrapped in a container which will have a location.
	if _, ok := err.(*container); !ok {
		err = &container{
			cause:    err,
			location: "",
			config:   c,
		}
	}

	for _, fn := range w {
		err = fn(err)
		if err == nil && c.AllowWrapperToDiscardError {
			return nil
		}
	}

	containerErr := &container{
		cause:    err,
		location: c.getLocation(skipFramesDelta),
		config:   c,
	}

	return containerErr
}

type container struct {
	cause    error
	location string
	config   *Fault
}

// Error behaves like most error wrapping libraries, it gives you all the error
// messages conjoined with ": ". This is useful only for internal error reports,
// never show this to an end-user or include it in responses as it may reveal
// internal technical information about your application stack.
func (f *container) Error() string {
	chain := f.config.Flatten(f)

	if f.config != nil && f.config.BuildDefaultErrorMessage != nil {
		return f.config.BuildDefaultErrorMessage(chain)
	}

	errs := []string{}

	// reverse iterate since the chain is in caller order
	for i := len(chain) - 1; i >= 0; i-- {
		message := chain[i].Message
		if message != "" && !isInternalString(message) {
			errs = append(errs, chain[i].Message)
		}
	}

	message := strings.Join(errs, ": ")
	if message == "" {
		message = "(no error message provided)"
	}

	return message
}

func (f *container) Unwrap() error { return f.cause }

func (f *container) Format(s fmt.State, verb rune) {
	if f.config != nil && f.config.FormatErrorMessage != nil {
		f.config.FormatErrorMessage(f.config.Flatten(f), s, verb)
		return
	}
	switch verb {
	case 'v':
		if s.Flag('+') {
			u := f.config.Flatten(f)
			for _, v := range u {
				if v.Message != "" {
					fmt.Fprintf(s, "%s\n", v.Message)
				}
				if v.Location != "" {
					fmt.Fprintf(s, "\t%s\n", v.Location)
				}
			}
			return
		}
		if s.Flag('#') {
			fmt.Fprintf(s, "&container{cause: %#v, location: %#v}", f.cause, f.location)
			return
		}

		fallthrough

	case 's':
		fmt.Fprint(s, f.Error())
	}
}

// getLocation returns the file and line where the error occurred.
func (c *Fault) getLocation(skipFramesDelta int) string {
	if c != nil && c.GetLocation != nil {
		return c.GetLocation(skipFramesDelta)
	}

	return defaultGetLocation(skipFramesDelta)
}

// defaultGetLocation is the default implementation of getLocation.
func defaultGetLocation(skipFramesDelta int) string {
	pc := make([]uintptr, 1)
	// +4 because:
	// 1: defaultGetLocation
	// 2: c.getLocation
	// 3: Calling function (Wrap/New)
	// 4: User code
	// +skipFramesDelta to adjust for additional wrappers
	runtime.Callers(4+skipFramesDelta, pc)
	cf := runtime.CallersFrames(pc)
	f, _ := cf.Next()

	return fmt.Sprintf("%s:%d", f.File, f.Line)
}

// isInternalString returns true for messages like <fctx> which are placeholders
func isInternalString(s string) bool {
	return strings.HasPrefix(s, "<") && strings.HasSuffix(s, ">")
}
