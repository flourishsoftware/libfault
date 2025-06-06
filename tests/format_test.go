package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/flourishsoftware/libfault/fctx"
	"github.com/flourishsoftware/libfault/ftag"
	"github.com/flourishsoftware/libfault/tests/internal/fault"
	"github.com/stretchr/testify/assert"
)

func TestFormatStdlibSentinelError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(1)

	a.Equal("failed to call function: stdlib sentinel error", err.Error())
	a.Equal("failed to call function: stdlib sentinel error", fmt.Sprintf("%s", err))
	a.Equal("failed to call function: stdlib sentinel error", fmt.Sprintf("%v", err))
	a.Regexp(`stdlib sentinel error
\s+.+fault/tests/test_callers.go:29
failed to call function
\s+.+fault/tests/test_callers.go:20
`, fmt.Sprintf("%+v", err))
}

func TestFormatFaultSentinelError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(2)

	a.Equal("failed to call function: fault sentinel error", err.Error())
	a.Equal("failed to call function: fault sentinel error", fmt.Sprintf("%s", err))
	a.Equal("failed to call function: fault sentinel error", fmt.Sprintf("%v", err))
	a.Regexp(`fault sentinel error
\s+.+fault/tests/root.go:15
\s+.+fault/tests/test_callers.go:29
failed to call function
\s+.+fault/tests/test_callers.go:20
\s+.+fault/tests/test_callers.go:11
`, fmt.Sprintf("%+v", err))
}

func TestFormatStdlibInlineError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(3)

	a.Equal("failed to call function: stdlib root cause error", err.Error())
	a.Equal("failed to call function: stdlib root cause error", fmt.Sprintf("%s", err))
	a.Equal("failed to call function: stdlib root cause error", fmt.Sprintf("%v", err))
	a.Regexp(`stdlib root cause error
\s+.+fault/tests/test_callers.go:29
failed to call function
\s+.+fault/tests/test_callers.go:20
`, fmt.Sprintf("%+v", err))
}

func TestFormatFaultInlineError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(4)

	a.Equal("failed to call function: fault root cause error", err.Error())
	a.Equal("failed to call function: fault root cause error", fmt.Sprintf("%s", err))
	a.Equal("failed to call function: fault root cause error", fmt.Sprintf("%v", err))
	a.Regexp(`fault root cause error
\s+.+fault/tests/root.go:28
\s+.+fault/tests/test_callers.go:29
failed to call function
\s+.+fault/tests/test_callers.go:20
\s+.+fault/tests/test_callers.go:11
`, fmt.Sprintf("%+v", err))
}

func TestFormatStdlibSentinelErrorWrappedWithoutMessage(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()

	err := errorCaller(1)
	err = fault.Wrap(err, fctx.With(ctx))
	err = fault.Wrap(err, ftag.With(ftag.Internal))

	a.NotContains(err.Error(), "<fctx>", "filtered out by .Error()")
	a.NotContains(err.Error(), "<ftag>", "filtered out by .Error()")

	a.Equal("failed to call function: stdlib sentinel error", err.Error())
	a.Equal("failed to call function: stdlib sentinel error", fmt.Sprintf("%s", err))
	a.Equal("failed to call function: stdlib sentinel error", fmt.Sprintf("%v", err))
	a.Regexp(`stdlib sentinel error
\s+.+fault/tests/test_callers.go:29
failed to call function
\s+.+fault/tests/test_callers.go:20
`, fmt.Sprintf("%+v", err))
}
