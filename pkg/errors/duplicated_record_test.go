package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/evzpav/user-auth/pkg/errors"
)

func TestNewDuplicatedRecord(t *testing.T) {
	var DuplicatedRecordCode errors.Code
	DuplicatedRecordCode = "DUPLICATED_RECORD"

	DuplicatedRecord := errors.NewDuplicatedRecord(DuplicatedRecordCode)
	assert.Equal(t, "<DUPLICATED_RECORD> duplicated record", DuplicatedRecord.Error())
	assert.Equal(t, DuplicatedRecordCode, DuplicatedRecord.GetCode())
	assert.Equal(t, "", DuplicatedRecord.GetMessage())

	DuplicatedRecord = errors.NewDuplicatedRecord(DuplicatedRecordCode).WithArg("id", "123abc123abc")
	assert.Equal(t, "<DUPLICATED_RECORD> duplicated record (id: 123abc123abc)", DuplicatedRecord.Error())

	DuplicatedRecord = errors.NewDuplicatedRecord(DuplicatedRecordCode).WithMessage("custom message duplicated record").WithArg("id", "123abc123abc")
	assert.Equal(t, "<DUPLICATED_RECORD> custom message duplicated record (id: 123abc123abc)", DuplicatedRecord.Error())
	assert.Equal(t, "custom message duplicated record", DuplicatedRecord.GetMessage())
}
