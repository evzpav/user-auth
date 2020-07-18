package errors

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

//Code is a custom type to identify the error code
type Code string

//Args is a custom type to store parameters of an error
type Args map[string]interface{}

func (a Args) getKeyValueAsString() string {
	values := make([]string, 0, len(a))

	for key, value := range a {
		values = append(values, fmt.Sprintf("%s: %v", key, value))
	}

	sort.Strings(values)
	return strings.Join(values, ", ")
}

//Describer is the interface that normalize how the custom error should be
type Describer interface {
	GetArgs() Args
	GetCode() Code
	GetMessage() string
	Error() string
}

//Info contains the base fields of custom error
type Info struct {
	Code           Code
	Args           Args
	Message        string
	defaultMessage string
}

//Error format the message as follow <error code> error message (error argument when it is not empty)
func (e *Info) Error() string {
	var buf bytes.Buffer
	if e.Code != "" {
		if _, err := fmt.Fprintf(&buf, "<%s> ", e.Code); err != nil {
			return fmt.Sprintf("%s : %s", string(e.Code), err.Error())
		}
	}

	if strings.TrimSpace(e.Message) == "" {
		if _, err := buf.WriteString(e.defaultMessage); err != nil {
			return fmt.Sprintf("%s : %s", string(e.Code), err.Error())
		}
	} else {
		if _, err := buf.WriteString(e.Message); err != nil {
			return fmt.Sprintf("%s : %s", string(e.Code), err.Error())
		}
	}

	if len(e.Args) > 0 {
		if _, err := buf.WriteString(" ("); err != nil {
			return fmt.Sprintf("%s : %s", string(e.Code), err.Error())
		}
		if _, err := buf.WriteString(e.Args.getKeyValueAsString()); err != nil {
			return fmt.Sprintf("%s : %s", string(e.Code), err.Error())
		}
		if _, err := buf.WriteString(")"); err != nil {
			return fmt.Sprintf("%s : %s", string(e.Code), err.Error())
		}
	}

	return buf.String()
}

//DescriberCast try to cast the base error into the Describer one
func DescriberCast(err error) (Describer, bool) {
	if err == nil {
		return nil, false
	}
	e, ok := err.(Describer)
	return e, ok
}
