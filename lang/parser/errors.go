package parser

import "github.com/lohvht/went/lang/token"

// SyntaxError refers to the error of lexing and parsing a piece of went syntax
type SyntaxError struct {
	token.GenericError
	errorname string
}

// NewSyntaxError returns a went syntax error
func NewSyntaxError(inputName string, pos token.Pos, msg string) *SyntaxError {
	return &SyntaxError{
		GenericError: token.GenericError{Input: inputName, Pos: pos, Msg: msg},
		errorname:    "SyntaxError",
	}
}

func (e SyntaxError) Error() string {
	return e.StandardErrorMessageFormat(e.errorname)
}
