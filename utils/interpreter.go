package utils

import (
	"fmt"
	"runtime"
	"strconv"
)

// Interpreter contains the already parsed inputs, as well as
// the various already defined variables/functions
// TODO: scopes
type Interpreter struct {
	tokeniser *lexer
	currTkn   token
	peekCount int
}

// type tokenList struct {
// 	tokens
// }

// next returns the next token
// func (i *Interpreter) nextToken() token {

// }

// nextNonSpace returns the next non-space (' ') token, sets currTkn and return
func (i *Interpreter) nextNonSpace() (tkn token) {
	for {
		tkn := i.tokeniser.nextToken()
		if tkn.typ != tokenSpace {
			break
		}
	}
	i.currTkn = tkn
	return tkn
}

// errorf formats the error and terminates processing
func (i *Interpreter) errorf(format string, args ...interface{}) {
	format = fmt.Sprintf("Interpreter: %d: %s", i.currTkn.line, format)
	panic(fmt.Errorf(format, args...))
}

// error terminates processing
func (i *Interpreter) error(err error) {
	i.errorf("%s", err)
}

// expect consumes the next token and guarantees it has the required type
func (i *Interpreter) expect(expectedType tokenType, context string) token {
	tkn := i.nextNonSpace()
	if tkn.typ != expectedType {
		i.unexpected(tkn, context)
	}
	return tkn
}

// expectOneOf consumes the next token and guarantees it has one of the required types
func (i *Interpreter) expectOneOf(expectedType1, expectedType2 tokenType, context string) token {
	tkn := i.nextNonSpace()
	if tkn.typ != expectedType1 && tkn.typ != expectedType2 {
		i.unexpected(tkn, context)
	}
	return tkn
}

// unexpected complains about the token and terminates processing.
func (i *Interpreter) unexpected(tkn token, context string) {
	i.errorf("Unexpected %s in %s", tkn, context)
}

// recover is the handler that turns panics into returns from the top level of Parse
func (i *Interpreter) recover(errp *error) {
	if e := recover(); e != nil {
		if _, ok := e.(runtime.Error); ok {
			panic(e)
		}
		if i != nil {
			i.tokeniser.drain()
			i.stopParse()
		}
		*errp = e.(error)
	}
}

// stopParse terminates parsing
func (i *Interpreter) stopParse() {
	i.tokeniser = nil
}

// Grammars

// Assuming factor is taking in and returning floats for now
// Separate the tokens and nodes
func (i *Interpreter) factor() float64 {
	f, _ := strconv.ParseFloat(i.nextNonSpace().value, 64)
	return f
}

// func (i *Interpreter) expr() float64 {
// 	result := i.factor()
// 	for {
// 		tkn := i.expectOneOf(tokenMultiply, tokenDivide, "expr: multiply, divide")
// 	}
// }
