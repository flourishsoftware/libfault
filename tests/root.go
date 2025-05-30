package tests

import (
	"errors"
	"fmt"

	"github.com/flourishsoftware/libfault/tests/internal/fault"
)

// the root cause of errors for testing. All line numbers here should remain the
// same so the tests don't have to be edited constantly when new cases are added

var (
	errSentinelStdlib = errors.New("stdlib sentinel error")
	errSentinelFault  = fault.New("fault sentinel error")
)

func rootCause(kind int) error {
	if kind == 0 {
		return nil
	} else if kind == 1 {
		return errSentinelStdlib
	} else if kind == 2 {
		return errSentinelFault
	} else if kind == 3 {
		return errors.New("stdlib root cause error")
	} else if kind == 4 {
		return fault.New("fault root cause error")
	} else if kind == 5 {
		return fmt.Errorf("errorf wrapped: %w", errSentinelStdlib)
	} else if kind == 6 {
		err := externalError()
		return fmt.Errorf("errorf wrapped external: %w", err)
	} else if kind == 7 {
		err := externalWrappedError()
		return fmt.Errorf("errorf wrapped external: %w", err)
	} else if kind == 8 {
		err := externalWrappedError()
		return err
	}
	return nil
}
