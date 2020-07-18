package errors

import (
	"fmt"
)

//DuplicatedRecordDescriber is the interface that contains Describer and indicates it's a not found error
type DuplicatedRecordDescriber interface {
	Describer
	DuplicatedRecord()
}

//DuplicatedRecordCast try to cast the base error into the DuplicatedRecord one
func DuplicatedRecordCast(err error) (DuplicatedRecordDescriber, bool) {
	if err == nil {
		return nil, false
	}
	e, ok := err.(DuplicatedRecordDescriber)
	return e, ok
}

//DuplicatedRecord should be used when the given resource was not found
type DuplicatedRecord struct {
	Info
}

//NewDuplicatedRecord creates a new instance of private DuplicatedRecord
func NewDuplicatedRecord(code Code) *DuplicatedRecord {
	err := &DuplicatedRecord{}
	err.Code = code
	err.defaultMessage = "duplicated record"
	return err
}

//WithMessage sets the message and returns the DuplicatedRecord
func (err *DuplicatedRecord) WithMessage(msg string) *DuplicatedRecord {
	err.Message = msg
	return err
}

//WithMessagef formats the message according to args paramets and set the message
func (err *DuplicatedRecord) WithMessagef(msg string, args ...interface{}) *DuplicatedRecord {
	err.Message = fmt.Sprintf(msg, args...)
	return err
}

//WithArg sets the single argument into error's arguments and returns the DuplicatedRecord
func (err *DuplicatedRecord) WithArg(key string, value interface{}) *DuplicatedRecord {
	if err.Args == nil {
		err.Args = make(map[string]interface{})
	}
	err.Args[key] = value
	return err
}

//GetMessage returns the error message
func (err *DuplicatedRecord) GetMessage() string {
	return err.Message
}

//GetCode returns the custom code of error
func (err *DuplicatedRecord) GetCode() Code {
	return err.Code
}

//Error builds the error according its message and code
func (err *DuplicatedRecord) Error() string {
	return err.Info.Error()
}

//GetArgs retrieves the arguments that belongs to error
func (err *DuplicatedRecord) GetArgs() Args {
	return err.Args
}

//DuplicatedRecord is a method that make the error to be an implementation of DuplicatedRecord interface
func (err *DuplicatedRecord) DuplicatedRecord() {}
