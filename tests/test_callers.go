package tests

import (
	"github.com/flourishsoftware/libfault/fmsg"
	"github.com/flourishsoftware/libfault/tests/internal/fault"
)

func errorCaller(kind int) error {
	err := errorCallerFromMiddleOfChain(kind)
	if err != nil {
		return fault.Wrap(err)
	}

	return nil
}

func errorCallerFromMiddleOfChain(kind int) error {
	err := errorProducerFromRootCause(kind)
	if err != nil {
		return fault.Wrap(err, fmsg.With("failed to call function"))
	}

	return nil
}

func errorProducerFromRootCause(kind int) error {
	err := rootCause(kind)
	if err != nil {
		return fault.Wrap(err)
	}

	return nil
}
