package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/evzpav/documents/pkg/errors"
)

func TestNewInvalidArgument(t *testing.T) {
	var invalidArgumentCode errors.Code
	invalidArgumentCode = "INVALID_ARGUMENT"

	invalidArgument := errors.NewInvalidArgument(invalidArgumentCode)
	assert.Equal(t, "<INVALID_ARGUMENT> invalid argument", invalidArgument.Error())
	assert.Equal(t, invalidArgumentCode, invalidArgument.GetCode())
	assert.Equal(t, "", invalidArgument.GetMessage())

	invalidArgument = errors.NewInvalidArgument(invalidArgumentCode).WithArg("id", "123abc123abc")
	assert.Equal(t, "<INVALID_ARGUMENT> invalid argument (id: 123abc123abc)", invalidArgument.Error())

	invalidArgument = errors.NewInvalidArgument(invalidArgumentCode).WithMessage("custom message invalid argument").WithArg("id", "123abc123abc")
	assert.Equal(t, "<INVALID_ARGUMENT> custom message invalid argument (id: 123abc123abc)", invalidArgument.Error())
	assert.Equal(t, "custom message invalid argument", invalidArgument.GetMessage())
}
