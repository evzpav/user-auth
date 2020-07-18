package errors

import (
	"fmt"
)

//NotFoundDescriber is the interface that contains Describer and indicates it's a not found error
type NotFoundDescriber interface {
	Describer
	NotFound()
}

//NotFoundCast try to cast the base error into the NotFound one
func NotFoundCast(err error) (NotFoundDescriber, bool) {
	if err == nil {
		return nil, false
	}
	e, ok := err.(NotFoundDescriber)
	return e, ok
}

//NotFound should be used when the given resource was not found
type NotFound struct {
	Info
}

//NewNotFound creates a new instance of private NotFound
func NewNotFound(code Code) *NotFound {
	err := &NotFound{}
	err.Code = code
	err.defaultMessage = "not found"
	return err
}

//WithMessage sets the message and returns the NotFound
func (err *NotFound) WithMessage(msg string) *NotFound {
	err.Message = msg
	return err
}

//WithMessagef formats the message according to args paramets and set the message
func (err *NotFound) WithMessagef(msg string, args ...interface{}) *NotFound {
	err.Message = fmt.Sprintf(msg, args...)
	return err
}

//WithArg sets the single argument into error's arguments and returns the NotFound
func (err *NotFound) WithArg(key string, value interface{}) *NotFound {
	if err.Args == nil {
		err.Args = make(map[string]interface{})
	}
	err.Args[key] = value
	return err
}

//GetMessage returns the error message
func (err *NotFound) GetMessage() string {
	return err.Message
}

//GetCode returns the custom code of error
func (err *NotFound) GetCode() Code {
	return err.Code
}

//Error builds the error according its message and code
func (err *NotFound) Error() string {
	return err.Info.Error()
}

//GetArgs retrieves the arguments that belongs to error
func (err *NotFound) GetArgs() Args {
	return err.Args
}

//NotFound is a method that make the error to be an implementation of NotFound interface
func (err *NotFound) NotFound() {}
