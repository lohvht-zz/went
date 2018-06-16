package utils

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Most of this package is listed and adapted from https://golang.org/src/text/template/parse/lex.go

// Token : String token read in and formed by the lexer
type Token struct {
	typ   tokenType // Type of this token
	pos   int       // Starting position, in bytes of this item in the input string
	value string    // value of this item
	line  int       // Line number at the start of this item
}

func (tok Token) String() string {
	switch {
	case tok.typ == tokenEOF:
		return "EOF"
	case tok.typ == tokenError:
		return tok.value
	case tok.typ > tokenKeyword:
		return fmt.Sprintf("<%s>", tok.value)
	case len(tok.value) > 10:
		return fmt.Sprintf("%.10q...", tok.value)
	}
	return fmt.Sprintf("%q", tok.value)
}

type tokenType int

// What do we want to tokenise?
// variables
// Numbers Constants
// Operations (+, -, *, /)
// String constants
// booleans
// objects
// arrays

// itemPipe       // pipe symbol
// itemRawString  // raw quoted string (includes quotes)
// itemRightDelim // right action delimiter
// itemRightParen // ')' inside action
// itemSpace      // run of spaces separating arguments
// itemString     // quoted string (includes quotes)
// itemText       // plain text
// itemVariable   // variable starting with '$', such as '$' or  '$1' or '$hello'
// // Keywords appear after all the rest.
// itemKeyword  // used only to delimit the keywords
// itemBlock    // block keyword
// itemDot      // the cursor, spelled '.'
// itemDefine   // define keyword
// itemElse     // else keyword
// itemEnd      // end keyword
// itemIf       // if keyword
// itemNil      // the untyped nil constant, easiest to treat as a keyword
// itemRange    // range keyword
// itemTemplate // template keyword
// itemWith     // with keyword

const (
	tokenError  tokenType = iota // error occurred; value is the text of error
	tokenBool                    // boolean constant
	tokenEquals                  // Equals ('=') sign to introduce a declaration
	tokenEOF
	tokenProperty   // alphanumeric identifier starting with '.', used in accessing object properties
	tokenIdentifier // alphanumeric identifier not starting with '.' may be a variable/function/class/struct
	tokenLeftParen  // left parenthesis '('
	tokenNumber     // simple number, including floating points
	tokenOp         // basic math operation ('+', '-', '/', '*')
	tokenSpace      // literally a space
	tokenRawString  // raw quoted string including the quotes (""")
	tokenOr
	tokenAnd
	// Keywords after all the rest
	tokenKeyword // Only used to delimit the keywords below
	tokenFunctionDef
	tokenVar // variable declaration using the keyword 'var'
	tokenIf
	tokenElse
	tokenElseIf
	tokenFor
)

var keyMap = map[string]tokenType{
	"funcDef": tokenFunctionDef,
	"if":      tokenIf,
	"else":    tokenElse,
	"elseIf":  tokenElseIf,
	"for":     tokenFor,
}

func (tok Token) toString() string {
	return string(tok.value)
}

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state
type stateFunc func(*lexer) stateFunc

type lexer struct {
	name  string // name of the input; used only for error reporting
	input string // string being scanned
	// leftDelim        string // start of Action (based on template)
	// rightDelim       string
	pos              int        // current position
	start            int        // start position of this token
	width            int        // width of the last rune read from input
	tokens           chan Token // channel of the scanned items
	paranthesisDepth int        // nesting depth of () brackets
	line             int        // 1 + number of newlines seen
}

// next returns the next rune in the input
func (l *lexer) next() rune {
	if l.pos > len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeLastRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	if r == '\n' {
		l.line++
	}
	return r
}

// peek returns but does not consume next rune in the input
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune, can only be called once per call of next
func (l *lexer) backup() {
	l.pos -= l.width
	// correct newline count
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

// emit passes a token back to the client
func (l *lexer) emit(typ tokenType) {
	l.tokens <- Token{typ, l.start, l.input[l.start:l.pos], l.line}
	switch typ {
	case tokenRawString:
		l.line += strings.Count(l.input[l.start:l.pos], "\n")
	}
	l.start = l.pos
}

// skips over the pending input before this point
func (l *lexer) ignore() {
	l.line += strings.Count(l.input[l.start:l.pos], "\n")
	l.start = l.pos
}

// accept consumes the next rune if its from the valid set
func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set
func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

// errorf returns an error token and terminates the scan by passing back a nil
// pointer that will be the next state, terminating l.nextToken.
// also emits an error token.
func (l *lexer) errorf(format string, args ...interface{}) stateFunc {
	l.tokens <- Token{tokenError, l.start, fmt.Sprintf(format, args...), l.line}
	return nil
}

// nextToken returns the next token from the input
// called by the parser, not in the lexing goroutine
func (l *lexer) nextToken() Token {
	return <-l.tokens
}

// drain drains the output so that the lexing goroutine will exit
// Called by the parser, not in lexing goroutine
func (l *lexer) drain() {
	for range l.tokens {
	}
}

// lex creates a new scanner for the input string
func lex(name, input string) *lexer {
	l := &lexer{
		name:   name,
		input:  input,
		tokens: make(chan Token),
		line:   1,
	}
	go l.run()
	return l
}

// run starts the state machine for the lexer
func (l *lexer) run() {
	for state := lexCode; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

// State functions

const (
	leftComment  = "/*"
	rightComment = "*/"
)

// lexCode scans the main body of the code
func lexCode(l *lexer) stateFunc {
	switch r := l.next(); {
	case r == eof:
		// if l.paranthesisDepth != 0 {
		// 	return l.errorf("Unclosed left paranthesis '('")
		// }
		l.emit(tokenEOF)
	case isEndOfLine(r) || isSpace(r):
		return lexSpace
	case r == '=':
		// If double equals ('==') not an assigment
		if l.next() != '=' {
			l.emit(tokenEquals)
			// go back
			l.backup()
		}
	case r == '|':
		if l.next() != '|' {
			l.emit(tokenOr)
		}

	}
	// if reached at this part, lexer has correctly reached the EOF
	l.emit(tokenEOF)
	return nil
}

// TODO: implement me!
func lexSpace(l *lexer) stateFunc {
	return nil
}

// Utility Functions

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
