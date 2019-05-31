package runtime

import "github.com/lohvht/went/lang/token"

// RuntimeError refers to the error occured during the runtime of a went program
type RuntimeError struct {
	token.GenericError
	errorname string
}

// NewRuntimeError returns a went syntax error
func NewRuntimeError(inputName string, pos token.Pos, msg string) *RuntimeError {
	return &RuntimeError{
		GenericError: token.GenericError{Input: inputName, Pos: pos, Msg: msg},
		errorname:    "RuntimeError",
	}
}

func (e RuntimeError) Error() string {
	return e.StandardErrorMessageFormat(e.errorname)
}
