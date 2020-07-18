package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/evzpav/user-auth/pkg/errors"
)

func TestNewNotAuthorized(t *testing.T) {
	var NotAuthorizedCode errors.Code
	NotAuthorizedCode = "NOT_AUTHORIZED"

	NotAuthorized := errors.NewNotAuthorized(NotAuthorizedCode)
	assert.Equal(t, "<NOT_AUTHORIZED> not authorized", NotAuthorized.Error())
	assert.Equal(t, NotAuthorizedCode, NotAuthorized.GetCode())
	assert.Equal(t, "", NotAuthorized.GetMessage())

	NotAuthorized = errors.NewNotAuthorized(NotAuthorizedCode).WithArg("id", "123abc123abc")
	assert.Equal(t, "<NOT_AUTHORIZED> not authorized (id: 123abc123abc)", NotAuthorized.Error())

	NotAuthorized = errors.NewNotAuthorized(NotAuthorizedCode).WithMessage("custom message not authorized").WithArg("id", "123abc123abc")
	assert.Equal(t, "<NOT_AUTHORIZED> custom message not authorized (id: 123abc123abc)", NotAuthorized.Error())
	assert.Equal(t, "custom message not authorized", NotAuthorized.GetMessage())
}
