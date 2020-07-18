package errors_test

import (
	"testing"

	"gitlab.com/evzpav/user-auth/pkg/errors"

	"github.com/stretchr/testify/assert"
)

func TestNewNotFound(t *testing.T) {
	notFound := errors.NewNotFound("JOB_NOT_FOUND")
	assert.Equal(t, "<JOB_NOT_FOUND> not found", notFound.Error())

	notFound = errors.NewNotFound("JOB_NOT_FOUND").WithArg("id", "123abc123abc")
	assert.Equal(t, "<JOB_NOT_FOUND> not found (id: 123abc123abc)", notFound.Error())

	notFound = errors.NewNotFound("JOB_NOT_FOUND").WithMessage("job not found").WithArg("id", "123abc123abc")
	assert.Equal(t, "<JOB_NOT_FOUND> job not found (id: 123abc123abc)", notFound.Error())
}
