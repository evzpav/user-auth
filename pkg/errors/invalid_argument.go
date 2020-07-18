package errors

import (
	"fmt"
)

//InvalidArgumentDescriber is the interface that contains Describer and indicates it's a not found error
type InvalidArgumentDescriber interface {
	Describer
	InvalidArgument()
}

//InvalidArgumentCast try to cast the base error into the InvalidArgument one
func InvalidArgumentCast(err error) (InvalidArgumentDescriber, bool) {
	if err == nil {
		return nil, false
	}
	e, ok := err.(InvalidArgumentDescriber)
	return e, ok
}

//InvalidArgument should be used when the given resource was not found
type InvalidArgument struct {
	Info
}

//NewInvalidArgument creates a new instance of private InvalidArgument
func NewInvalidArgument(code Code) *InvalidArgument {
	err := &InvalidArgument{}
	err.Code = code
	err.defaultMessage = "invalid argument"
	return err
}

//WithMessage sets the message and returns the InvalidArgument
func (err *InvalidArgument) WithMessage(msg string) *InvalidArgument {
	err.Message = msg
	return err
}

//WithMessagef formats the message according to args paramets and set the message
func (err *InvalidArgument) WithMessagef(msg string, args ...interface{}) *InvalidArgument {
	err.Message = fmt.Sprintf(msg, args...)
	return err
}

//WithArg sets the single argument into error's arguments and returns the InvalidArgument
func (err *InvalidArgument) WithArg(key string, value interface{}) *InvalidArgument {
	if err.Args == nil {
		err.Args = make(map[string]interface{})
	}
	err.Args[key] = value
	return err
}

//GetMessage returns the error message
func (err *InvalidArgument) GetMessage() string {
	return err.Message
}

//GetCode returns the custom code of error
func (err *InvalidArgument) GetCode() Code {
	return err.Code
}

//Error builds the error according its message and code
func (err *InvalidArgument) Error() string {
	return err.Info.Error()
}

//GetArgs retrieves the arguments that belongs to error
func (err *InvalidArgument) GetArgs() Args {
	return err.Args
}

//InvalidArgument is a method that make the error to be an implementation of InvalidArgument interface
func (err *InvalidArgument) InvalidArgument() {}
