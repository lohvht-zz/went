package token

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// WentError is the  error type that is used for all reported went errors
type WentError interface {
	error
	InputName() string    // name of the input string, usually a filename
	Position() (int, int) // the position within the input string, line then column
	Message() string
}

// GenericError is the base error type of all went errors, it should be embedded
// when implementing a new error in went. The position Pos  if valid points to
// beginning of offending token and error condition as described by the message.
type GenericError struct {
	Input string
	Pos   Pos
	Msg   string
}

// InputName for WentError Interface
func (e GenericError) InputName() string { return e.Input }

// Position for WentError Interface
func (e GenericError) Position() (l int, c int) {
	l, c = e.Pos.decompose()
	return
}

// Message for WentError Interface
func (e GenericError) Message() string { return e.Msg }

// InputNamePos returns a string representation of <InputName>:<line#>:<col#>
// it can take the following forms:
// <InputName>:<line#>:<col#>
// <InputName>:<line#>
// <line#>:<col#>
// "" => only happens when InputName is empty, and Pos is not valid
func (e GenericError) InputNamePos() string {
	s := e.InputName()
	if e.Pos.IsValid() {
		if s != "" {
			s += ":"
		}
		s += e.Pos.String()
	}
	return s
}

// StandardErrorMessageFormat returns a string that adheres to the standard error format
// if inputNamePos and errorType are both "", return only the message
// if only inputNamePos is empty, return "[errorType]: message"
// else, return "[errorType]:inputName:l:c: message"
func (e GenericError) StandardErrorMessageFormat(errorType string) string {
	s := e.InputNamePos()
	switch {
	case s == "" && errorType == "":
		return e.Msg
	case s == "" && errorType != "":
		return "[" + errorType + "]: " + e.Msg
	case s != "" && errorType == "":
		return s + ": " + e.Msg
	default:
		return "[" + errorType + "]:" + s + ": " + e.Msg
	}
}

func (e GenericError) Error() string {
	return e.StandardErrorMessageFormat("")
}

// NewGenericError returns a generic went error
func NewGenericError(inputname string, pos Pos, msg string) *GenericError {
	return &GenericError{inputname, pos, msg}
}

// ErrorList is a list of WentErrors
type ErrorList []WentError

// Add adds an Error with given position and error message to an ErrorList.
func (p *ErrorList) Add(e WentError) { *p = append(*p, e) }

// Reset resets an ErrorList to no errors.
func (p *ErrorList) Reset() { *p = (*p)[0:0] }

// ErrorList implements the sort Interface.

// Len for sort interface
func (p ErrorList) Len() int { return len(p) }

//Swap for sort interface
func (p ErrorList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p ErrorList) Less(i, j int) bool {
	// Note that it is not sufficient to simply compare file offsets because
	// the offsets do not reflect modified line information (through //line
	// comments).
	if p[i].InputName() != p[j].InputName() {
		return p[i].InputName() < p[j].InputName()
	}
	el, ec := p[i].Position()
	fl, fc := p[j].Position()

	if el != fl {
		return el < fl
	}
	if ec != fc {
		return ec < fc
	}
	return p[i].Message() < p[j].Message()
}

// Sort sorts an ErrorList. *Error entries are sorted by position,
// other errors are sorted by error message, and before any *Error
// entry.
func (p ErrorList) Sort() { sort.Sort(p) }

// RemoveMultiples sorts an ErrorList and removes all but the first error per line.
func (p *ErrorList) RemoveMultiples() {
	sort.Sort(p)
	var lastFn string
	var lastLine int
	i := 0
	for _, e := range *p {
		if currLine, _ := e.Position(); e.InputName() != lastFn || currLine != lastLine {
			lastLine = currLine
			(*p)[i] = e
			i++
		}
	}
	(*p) = (*p)[0:i]
}

// Error interface, an ErrorList implements it
func (p ErrorList) Error() string {
	switch len(p) {
	case 0:
		return "no errors"
	case 1:
		return p[0].Error()
	}
	// NOTE: Printing here for convenience
	PrintError(os.Stdout, p)
	return fmt.Sprintf("%s (and %d more errors)", p[0], len(p)-1)
}

// Err returns an error equivalent to this error list.
// If the list is empty, Err returns nil.
func (p ErrorList) Err() error {
	if len(p) == 0 {
		return nil
	}
	return p
}

// PrintError is a utility function that prints a list of errors to w,
// one error per line, if the err parameter is an ErrorList. Otherwise
// it prints the err string.
//
func PrintError(w io.Writer, err error) {
	if list, ok := err.(ErrorList); ok {
		for _, e := range list {
			fmt.Fprintf(w, "%s\n", e)
		}
	} else if err != nil {
		fmt.Fprintf(w, "%s\n", err)
	}
}
