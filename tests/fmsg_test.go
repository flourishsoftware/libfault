package tests

import (
	"errors"
	"testing"

	"github.com/josephbuchma/libfault/fmsg"
	"github.com/josephbuchma/libfault/tests/internal/fault"
	"github.com/stretchr/testify/assert"
)

func TestWithOne(t *testing.T) {
	err := fault.Wrap(errors.New("a problem"), fmsg.WithDesc("shit happened", "Shit happened."))
	out := fmsg.GetIssue(err)

	assert.Equal(t, "Shit happened.", out)
}

func TestWithNone(t *testing.T) {
	err := errors.New("a problem")
	out := fmsg.GetIssue(err)

	assert.Equal(t, "", out)
}

func TestWithMany(t *testing.T) {
	err := fault.Wrap(errors.New("the original problem"))
	err = fault.Wrap(err, fmsg.WithDesc("layer 1", "The post was not found."))
	err = fault.Wrap(err, fmsg.WithDesc("layer 2", "Unable to reply to post."))
	err = fault.Wrap(err, fmsg.WithDesc("layer 3", "Your reply draft has been saved however we could not publish it."))
	out := fmsg.GetIssue(err)

	assert.Equal(t, "Your reply draft has been saved however we could not publish it. Unable to reply to post. The post was not found.", out)
}

func TestWithManySlice(t *testing.T) {
	err := errors.New("the original problem")

	err = fault.Wrap(err, fmsg.WithDesc("layer 1", "The post was not found."))
	err = fault.Wrap(err, fmsg.WithDesc("layer 2", "Unable to reply to post."))
	err = fault.Wrap(err, fmsg.WithDesc("layer 3", "Your reply draft has been saved however we could not publish it."))
	out := fmsg.GetIssues(err)

	assert.Len(t, out, 3)
	assert.Equal(t, []string{"Your reply draft has been saved however we could not publish it.", "Unable to reply to post.", "The post was not found."}, out)
}

func TestInternalMessageFallback(t *testing.T) {
	err := fault.Wrap(errors.New("underlying"), fmsg.WithDesc("", "External message."))
	assert.Equal(t, "External message: underlying", err.Error())
}

func TestCollapsingMessages(t *testing.T) {
	err := errors.New("underlying")
	err = fault.Wrap(err,
		fmsg.WithDesc("internal", ""),
		fmsg.WithDesc("", "External message."))

	assert.Equal(t, "internal: underlying", err.Error())
	assert.Equal(t, "External message.", fmsg.GetIssue(err))
}
