package errors

import (
	"fmt"
)

//NotAuthorizedDescriber is the interface that contains Describer and indicates it's a not found error
type NotAuthorizedDescriber interface {
	Describer
	NotAuthorized()
}

//NotAuthorizedCast try to cast the base error into the NotAuthorized one
func NotAuthorizedCast(err error) (NotAuthorizedDescriber, bool) {
	if err == nil {
		return nil, false
	}
	e, ok := err.(NotAuthorizedDescriber)
	return e, ok
}

//NotAuthorized should be used when the given resource was not found
type NotAuthorized struct {
	Info
}

//NewNotAuthorized creates a new instance of private NotAuthorized
func NewNotAuthorized(code Code) *NotAuthorized {
	err := &NotAuthorized{}
	err.Code = code
	err.defaultMessage = "not authorized"
	return err
}

//WithMessage sets the message and returns the NotAuthorized
func (err *NotAuthorized) WithMessage(msg string) *NotAuthorized {
	err.Message = msg
	return err
}

//WithMessagef formats the message according to args paramets and set the message
func (err *NotAuthorized) WithMessagef(msg string, args ...interface{}) *NotAuthorized {
	err.Message = fmt.Sprintf(msg, args...)
	return err
}

//WithArg sets the single argument into error's arguments and returns the NotAuthorized
func (err *NotAuthorized) WithArg(key string, value interface{}) *NotAuthorized {
	if err.Args == nil {
		err.Args = make(map[string]interface{})
	}
	err.Args[key] = value
	return err
}

//GetMessage returns the error message
func (err *NotAuthorized) GetMessage() string {
	return err.Message
}

//GetCode returns the custom code of error
func (err *NotAuthorized) GetCode() Code {
	return err.Code
}

//Error builds the error according its message and code
func (err *NotAuthorized) Error() string {
	return err.Info.Error()
}

//GetArgs retrieves the arguments that belongs to error
func (err *NotAuthorized) GetArgs() Args {
	return err.Args
}

//NotAuthorized is a method that make the error to be an implementation of NotAuthorized interface
func (err *NotAuthorized) NotAuthorized() {}
