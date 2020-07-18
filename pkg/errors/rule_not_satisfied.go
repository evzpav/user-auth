package errors

import (
	"fmt"
)

//RuleNotSatisfiedDescriber is the interface that contains Describer and indicates it's a not found error
type RuleNotSatisfiedDescriber interface {
	Describer
	RuleNotSatisfied()
}

//RuleNotSatisfiedCast try to cast the base error into the RuleNotSatisfied one
func RuleNotSatisfiedCast(err error) (RuleNotSatisfiedDescriber, bool) {
	if err == nil {
		return nil, false
	}
	e, ok := err.(RuleNotSatisfiedDescriber)
	return e, ok
}

//RuleNotSatisfied should be used when the given resource was not found
type RuleNotSatisfied struct {
	Info
}

//NewRuleNotSatisfied creates a new instance of private RuleNotSatisfied
func NewRuleNotSatisfied(code Code) *RuleNotSatisfied {
	err := &RuleNotSatisfied{}
	err.Code = code
	err.defaultMessage = "internal rule was not satisfied"
	return err
}

//WithMessage sets the message and returns the RuleNotSatisfied
func (err *RuleNotSatisfied) WithMessage(msg string) *RuleNotSatisfied {
	err.Message = msg
	return err
}

//WithMessagef formats the message according to args paramets and set the message
func (err *RuleNotSatisfied) WithMessagef(msg string, args ...interface{}) *RuleNotSatisfied {
	err.Message = fmt.Sprintf(msg, args...)
	return err
}

//WithArg sets the single argument into error's arguments and returns the RuleNotSatisfied
func (err *RuleNotSatisfied) WithArg(key string, value interface{}) *RuleNotSatisfied {
	if err.Args == nil {
		err.Args = make(map[string]interface{})
	}
	err.Args[key] = value
	return err
}

//GetMessage returns the error message
func (err *RuleNotSatisfied) GetMessage() string {
	return err.Message
}

//GetCode returns the custom code of error
func (err *RuleNotSatisfied) GetCode() Code {
	return err.Code
}

//Error builds the error according its message and code
func (err *RuleNotSatisfied) Error() string {
	return err.Info.Error()
}

//GetArgs retrieves the arguments that belongs to error
func (err *RuleNotSatisfied) GetArgs() Args {
	return err.Args
}

//RuleNotSatisfied is a method that make the error to be an implementation of RuleNotSatisfied interface
func (err *RuleNotSatisfied) RuleNotSatisfied() {}
