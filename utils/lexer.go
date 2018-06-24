package utils

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Most of this package is listed and adapted from
// https://golang.org/src/text/template/parse/lex.go

/**
 * token Definition
 */
type token struct {
	typ   tokenType // Type of this token
	pos   Pos       // Starting position, in bytes of this item in the input string
	value string    // value of this item
	line  int       // Line number at the start of this item
}

func (tok token) String() string {
	switch {
	case tok.typ == tokenEOF:
		return "EOF"
	case tok.typ == tokenError:
		return tok.value
	case tok.typ > tokenKeyword:
		return fmt.Sprintf("<%s>", tok.value)
		// case len(tok.value) > 10:
		// 	return fmt.Sprintf("%.10q...", tok.value) // commented this, fullprint
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

const (
	tokenError        tokenType = iota // error occurred; value is the text of error
	tokenEquals                        // Equals ('=') sign to introduce a declaration
	tokenDoubleEquals                  // Double Equals ('==') sign to test for equality
	tokenNotEquals                     // Not equals ('!=') sign, to test for not equality
	tokenLogicalNot                    // exclamation mark ('!') as logical not
	tokenEOF
	tokenProperty   // alphanumeric identifier starting with '.', used in accessing object properties
	tokenIdentifier // alphanumeric identifier not starting with '.' may be a variable/function/class/struct
	tokenLeftParen  // left parenthesis '('
	tokenRightParan // right parenthesis ')'
	tokenLeftBrace  // left brace '{'
	tokenRightBrace // right brace '}'
	// tokenLeftSquare // left square bracket '['
	// tokenRightSquare // right square bracket ']'
	// tolenColon // colon symbol ':'
	tokenOp      // math operations such as '+', '/', '*', '-', '%'
	tokenSpace   // literally a space (' ')
	tokenNewline // newline token ('\r\n' or '\n')
	tokenOr      // OR symbol, represented by ('||')
	tokenAnd     // AND sumbol, represented by ('&&')
	// Literal tokens (not including object, array)
	tokenBool         // boolean literal (true, false)
	tokenNumber       // Integer64 or float64 numbers
	tokenQuotedString // Singly quoted ('\'') strings, escaped using a single '\' char
	tokenRawString    // tilde quoted ('`') strings, intepreted as-is, with no way of escaping
	// Keywords after all the rest
	tokenKeyword     // Only used to delimit the keywords below
	tokenFunctionDef // 'func' keyword for function definition
	tokenVar         // variable declaration using the keyword 'var'
	tokenIf          // 'if' keyword
	tokenElse        // 'else' keyword
	tokenElseIf      // 'else if' keyword
	tokenFor         // 'for' keyword, for loops
	tokenRange       // 'range' keyword, to be used with for-loops to create a shortcut
	tokenNull        // 'null' constant, treated as a keyword
)

var keyMap = map[string]tokenType{
	"func":   tokenFunctionDef,
	"var":    tokenVar,
	"if":     tokenIf,
	"else":   tokenElse,
	"elseIf": tokenElseIf,
	"for":    tokenFor,
	"range":  tokenRange,
	"null":   tokenNull,
}

const eof = -1

/**
 * lexer Definition
 */
type lexer struct {
	name  string // name of the input; used only for error reporting
	input string // string being scanned
	// leftDelim        string // start of Action (based on template)
	// rightDelim       string
	pos              Pos        // current position
	start            Pos        // start position of this token
	width            Pos        // width of the last rune read from input
	tokens           chan token // channel of the scanned items
	paranthesisDepth int        // nesting depth of () brackets
	bracesDepth      int        // nesting depth of {} brackets
	line             int        // 1 + number of newlines seen
}

// next returns the next rune in the input
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
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
	l.tokens <- token{typ, l.start, l.input[l.start:l.pos], l.line}
	// Some of the tokens contain text internally, if so, count their newlines
	switch typ {
	case tokenRawString, tokenQuotedString:
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
	l.tokens <- token{tokenError, l.start, fmt.Sprintf(format, args...), l.line}
	return nil
}

// nextToken returns the next token from the input
// called by the parser, not in the lexing goroutine
func (l *lexer) nextToken() token {
	return <-l.tokens
}

// drain drains the output so that the lexing goroutine will exit
// Called by the parser, not in lexing goroutine
func (l *lexer) drain() {
	for range l.tokens {
	}
}

// run starts the state machine for the lexer
func (l *lexer) run() {
	for state := lexCode; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *lexer) scanNumber() bool {
	// Optional leading sign
	digits := "0123456789"
	leadingSigns := "+-"
	l.accept(leadingSigns)
	l.acceptRun(digits)
	// Decimal
	if l.accept(".") {
		l.acceptRun(digits)
	}
	// Powers of 10
	if l.accept("eE") {
		l.accept(leadingSigns)
		l.accept(digits)
	}
	// Next thing after everything else
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}
	return true
}

// atIdentifierTerminator reports whether the input is at valid
// termination character to appear after an identifier
func (l *lexer) atIdentifierTerminator() bool {
	r := l.peek()
	if isSpace(r) || isEndOfLine(r) {
		return true
	}
	switch r {
	case
		eof,      // EOF character
		'.', ',', // DOT ('.') to denote .property, or commas
		'|', '&', // OR ('||'), or AND ('&&')
		'=',      // assignment/declaration ('='), or equality check ('==')
		')', '(', // Parenthesis '(', ')'
		'+', '-', '/', '*', '%': // Math operator signs, or start of a comment ('//', '/*')
		return true
	}
	return false
}

// tokenise creates a new scanner for the input string
func tokenise(name, input string) *lexer {
	l := &lexer{
		name:   name,
		input:  input,
		tokens: make(chan token),
		line:   1,
	}
	go l.run()
	return l
}

// State functions

// stateFn represents the state of the scanner as a function that returns the next state
type stateFunc func(*lexer) stateFunc

// lexCode scans the main body of the code, recursively returning itself
func lexCode(l *lexer) stateFunc {
	switch r := l.next(); {
	case r == eof: // Where the lexCode loop terminates, when it reaches EOF
		return lexEOF
	case isSpace(r):
		return lexSpace
	case isEndOfLine(r): // detects \r OR \n
		l.backup()
		return lexNewline
	case r == '!':
		if l.next() != '=' {
			l.backup()
			l.emit(tokenLogicalNot)
		} else {
			l.emit(tokenNotEquals)
		}
	case r == '=':
		if l.next() != '=' {
			l.backup()
			l.emit(tokenEquals)
		} else {
			l.emit(tokenDoubleEquals)
		}
	case r == '|':
		if l.next() == '|' {
			l.emit(tokenOr)
		} else {
			l.errorf("Expected '|' token")
		}
	case r == '&':
		if l.next() == '&' {
			l.emit(tokenAnd)
		} else {
			l.errorf("Expected '&' token")
		}
	case r == '"':
		return lexQuotedString
	case r == '`':
		return lexRawString
	case r == '.':
		// Special lookahead for ".property" so we don't break l.backup()
		if int(l.pos) < len(l.input) {
			r := l.input[l.pos]
			if r < '0' || r > '9' { // if its not a number
				return lexProperty
			}
		}
		fallthrough // '.' can start a number
	case '0' <= r && r <= '9':
		l.backup()
		return lexNumber
	case r == '+' || r == '-':
		// Special lookahead to look for a number from 0-9 from the next char
		if int(l.pos) < len(l.input) {
			r := l.input[l.pos]
			if '0' <= r && r <= '9' {
				l.backup()
				return lexNumber
			}
		}
		fallthrough // '+', '-' is a math operation instead of the start of a number
	case r == '*' || r == '%':
		l.emit(tokenOp)
	case r == '/':
		// Scans the next rune
		switch r := l.next(); {
		case r == '/':
			return lexSinglelineComment
		case r == '*':
			return lexMultilineComment
		default:
			// Default case, is divider ('/')
			l.backup()
			l.emit(tokenOp)
		}
	case isAlphaNumeric(r):
		l.backup()
		return lexIdentifier
	case r == '(':
		l.emit(tokenLeftParen)
		l.paranthesisDepth++
	case r == ')':
		l.emit(tokenRightParan)
		l.paranthesisDepth--
		if l.paranthesisDepth < 0 {
			return l.errorf("Unexpected right parenthesis %#U", r)
		}
	case r == '{':
		l.emit(tokenLeftBrace)
		l.bracesDepth++
	case r == '}':
		l.emit(tokenRightBrace)
		l.bracesDepth--
		if l.bracesDepth < 0 {
			return l.errorf("Unexpected right brace %#U", r)
		}
	default:
		return l.errorf("Unrecognised character in code: %#U", r)
	}
	return lexCode
}

// lexEOF emits the EOF token and handles the termination of the main lexCode loop
func lexEOF(l *lexer) stateFunc {
	if l.paranthesisDepth != 0 {
		return l.errorf("Unclosed left paranthesis '('")
	} else if l.bracesDepth != 0 {
		return l.errorf("Unclosed right brace '}'")
	}
	l.emit(tokenEOF)
	return nil
}

// lexSpace scans a run of space characters, One space has already been seen
func lexSpace(l *lexer) stateFunc {
	for isSpace(l.peek()) {
		l.next()
	}
	l.emit(tokenSpace)
	return lexCode
}

// lexNewline scans for a run of newline characters (either \r\n OR \n)
func lexNewline(l *lexer) stateFunc {
Loop:
	for {
		switch r := l.next(); {
		case r == '\n':
			// Absorb and go to next iteration
		case r == '\r' && l.peek() == '\n':
			l.next() // advance after \n
			// Absorb and go to next iteration
		default:
			l.backup()
			break Loop
		}
	}
	l.emit(tokenNewline)
	return lexCode
}

// lexQuotedString scans a quoted string, can be escaped using the '\' character
func lexQuotedString(l *lexer) stateFunc {
	startLine := l.line
Loop:
	for {
		switch l.next() {
		case '\\': // single '\' character as escape character
			if r := l.next(); r == eof {
				// restore line number to where the open quote is by replacing the l.line
				// Error out after that
				l.line = startLine
				return l.errorf("Unterminated Quoted String")
			} // Else just absorb and continue consuming the rest of the string
		case '"':
			break Loop
		}
	}
	l.emit(tokenQuotedString)
	return lexCode
}

// lexRawString scans a raw string delimited by '`' character
func lexRawString(l *lexer) stateFunc {
	startLine := l.line
Loop:
	for {
		switch l.next() {
		case eof:
			// restore line number to the location of the opening quote
			// will error out, okay to overwrite l.line
			l.line = startLine
			return l.errorf("Unterminated raw string")
		case '`':
			break Loop
		}
	}
	l.emit(tokenRawString)
	return lexCode
}

// lexProperty scans an object accessed by its property: .Alphanumeric
// the front '.' char has already been scanned
func lexProperty(l *lexer) stateFunc {
	if l.atIdentifierTerminator() { // If nothing interesting follows => '.'
		return l.errorf("Bad character: %#U", l.input[l.pos])
	}
	var r rune
	for {
		if r = l.next(); !isAlphaNumeric(r) {
			l.backup()
			break
		}
	}
	if !l.atIdentifierTerminator() {
		return l.errorf("Bad character: %#U", r)
	}
	l.emit(tokenProperty)
	return lexCode
}

// lexNumber scan for a decimal number, it isn't a perfect number scanner
// for e.g. it accepts '.' and '089', but when its wrong the input is invalid
func lexNumber(l *lexer) stateFunc {
	if !l.scanNumber() {
		return l.errorf("Bad number syntax: %q", l.input[l.start:l.pos])
	}
	l.emit(tokenNumber)
	return lexCode
}

// lexIdentifier scans an alphanumeric word
func lexIdentifier(l *lexer) stateFunc {
Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb until no more next alphanumeric characters
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			if !l.atIdentifierTerminator() {
				return l.errorf("Bad character: %#U", r)
			}
			switch {
			case keyMap[word] > tokenKeyword:
				l.emit(keyMap[word])
			case word == "true", word == "false":
				l.emit(tokenBool)
			default:
				l.emit(tokenIdentifier)
			}
			break Loop
		}
	}
	return lexCode
}

// lexSinglelineComment scans a single line comment ('//') and discards it
// The comment marker ('//') has already been consumed
// This assumes that the entire line is scanned, if no newline is detected, then
// it will basically count to EOF
func lexSinglelineComment(l *lexer) stateFunc {
	if i := strings.Index(l.input[l.pos:], "\n"); i < 0 {
		// Major assumption, if the index of newline ("\n") is not found, then the input
		// has only 1 single line with a comment somewhere on the line
		// Move the positional scanner to the end of the file
		fmt.Println("i is: ", i)
		l.pos += Pos(len(l.input[l.pos:]))
	} else {
		l.pos += Pos(i)
	}
	l.ignore()
	return lexCode
}

// lexMultilineComment scans for a multiline comment block ('/*', '*/') and discards it
// The left comment marker ('/*') has already been consumed
func lexMultilineComment(l *lexer) stateFunc {
	rightComment := "*/"
	i := strings.Index(l.input[l.pos:], rightComment)
	if i < 0 {
		return l.errorf("Multiline comment is not closed")
	}
	l.pos += Pos(i + len(rightComment))
	l.ignore()
	return lexCode
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
