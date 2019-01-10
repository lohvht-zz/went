package token

import (
	"fmt"
	"io"
	"sort"
)

// Error has the position Pos, which if valid, points to beginning of offending token
// and error condition as described by the message
type Error struct {
	Filename string
	Pos      Pos
	Msg      string
}

func (e Error) Error() string {
	s := e.Filename
	if e.Pos.IsValid() {
		if s != "" {
			s += ":"
		}
		s += e.Pos.String()
	}
	if s == "" {
		// return Msg if empty filename and invalid Pos
		return e.Msg
	}
	return s + ": " + e.Msg
}

// ErrorList is a list of *Errors
type ErrorList []*Error

// Add adds an Error with given position and error message to an ErrorList.
func (p *ErrorList) Add(filename string, pos Pos, msg string) {
	*p = append(*p, &Error{filename, pos, msg})
}

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
	if p[i].Filename != p[j].Filename {
		return p[i].Filename < p[j].Filename
	}
	el, ec := p[i].Pos.decompose()
	fl, fc := p[j].Pos.decompose()

	if el != fl {
		return el < fl
	}
	if ec != fc {
		return ec < fc
	}
	return p[i].Msg < p[j].Msg
}

// Sort sorts an ErrorList. *Error entries are sorted by position,
// other errors are sorted by error message, and before any *Error
// entry.
func (p ErrorList) Sort() { sort.Sort(p) }

// RemoveMultiples sorts an ErrorList and removes all but the first error per line.
func (p *ErrorList) RemoveMultiples() {
	sort.Sort(p)
	var lastFn string
	var lastPos Pos
	i := 0
	for _, e := range *p {
		if e.Filename != lastFn || e.Pos.Line() != lastPos.Line() {
			lastPos = e.Pos
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
