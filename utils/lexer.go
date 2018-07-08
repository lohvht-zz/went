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

const (
	tokenError tokenType = iota // error occurred; value is the text of error
	tokenEOF
	tokenDot         // Dot character '.'
	tokenIdentifier  // alphanumeric identifier
	tokenLeftParen   // left parenthesis '('
	tokenRightParan  // right parenthesis ')'
	tokenLeftBrace   // left brace '{'
	tokenRightBrace  // right brace '}'
	tokenLeftSquare  // left square bracket '['
	tokenRightSquare // right square bracket ']'
	tokenColon       // colon symbol ':'
	tokenSemicolon   // semi colon symbol ';', used to terminate some grammars

	// Literal tokens (not including object, array)
	tokenBool         // boolean literal (true, false)
	tokenNumber       // Integer64 or float64 numbers
	tokenQuotedString // Singly quoted ('\'') strings, escaped using a single '\' char
	tokenRawString    // tilde quoted ('`') strings, intepreted as-is, with no way of escaping

	// tokenOperators // Only used to delimit Operators below
	// Operators
	// Arithmetic Operators
	tokenPlus  // '+', can be used for strings
	tokenMinus // '-'
	tokenDiv   // '/'
	tokenMult  // '*'
	tokenMod   // '%'
	// Assignment Operators
	tokenAssign      // Equals ('=') sign for assigning
	tokenPlusAssign  // '+=', addition then assign, can be used for strings
	tokenMinusAssign // '-=', subtract then assign
	tokenDivAssign   // '/=', divide then assign
	tokenMultAssign  // '*=', multiply then assign
	tokenModAssign   // '%=', modulo then assign
	// Comparison Operators
	tokenEquals        // '==', test for value equality
	tokenNotEquals     // '!=', test for value inequality
	tokenGreater       // '>', test for greater than
	tokenSmaller       // '<', test for smaller than
	tokenGreaterEquals // '>=', test for greater than or equals to
	tokenSmallerEquals // '<=', test for smaller than or equals to
	// Logical Operators
	tokenLogicalNot // exclamation mark ('!') as logical not
	tokenOr         // OR symbol, represented by ('||')
	tokenAnd        // AND sumbol, represented by ('&&')

	// Keywords after all the rest
	tokenKeyword // Only used to delimit the keywords below
	tokenFunc    // 'func' keyword for function definition
	tokenIf      // 'if' keyword
	tokenElse    // 'else' keyword
	tokenElseIf  // 'elif' keyword
	tokenFor     // 'for' keyword, for loops
	tokenNull    // 'null' constant, treated as a keyword
	tokenWhile   // 'while' keyword
	tokenReturn  // 'return' keyword
	tokenIn      // 'in' keyword
	tokenBreak   // 'break' keyword
	tokenCont    // 'continue' keyword
)

var keyMap = map[string]tokenType{
	"func":     tokenFunc,
	"if":       tokenIf,
	"else":     tokenElse,
	"elif":     tokenElseIf,
	"for":      tokenFor,
	"null":     tokenNull,
	"while":    tokenWhile,
	"return":   tokenReturn,
	"in":       tokenIn,
	"break":    tokenBreak,
	"continue": tokenCont,
}

const eof = -1

/**
 * lexer Definition
 */
type lexer struct {
	name             string     // name of the input; used only for error reporting
	input            string     // string being scanned
	pos              Pos        // current position
	start            Pos        // start position of this token
	width            Pos        // width of the last rune read from input
	tokens           chan token // channel of the scanned items
	prevTokTyp       tokenType  // previous token type used for automatic semicolon insertion
	paranthesisDepth int        // nesting depth of () brackets
	bracesDepth      int        // nesting depth of {} brackets
	squareDepth      int        //nesting depth of [] brackets
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
// this will also update the last seen emitted token type
func (l *lexer) emit(typ tokenType) {
	l.tokens <- token{typ, l.start, l.input[l.start:l.pos], l.line}
	// Some of the tokens contain text internally, if so, count their newlines
	switch typ {
	case tokenRawString, tokenQuotedString:
		l.line += strings.Count(l.input[l.start:l.pos], "\n")
	}
	l.start = l.pos
	l.prevTokTyp = typ
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

// does not accept leading +=
func (l *lexer) scanNumber() bool {
	digits := "0123456789"
	leadingSigns := "+-"
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
	// Check if the next rune is alphanumeric (if so then its not a number anymore)
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
	case r == ':':
		l.emit(tokenColon)
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
				l.emit(tokenDot)
				return lexCode // emit the dot '.' and go back to lexCode
			}
		}
		fallthrough // '.' can start a number, especially next rune is a number
	case '0' <= r && r <= '9':
		l.backup()
		return lexNumber
	case r == '+' || r == '-' || r == '*' || r == '%' || // Math signs
		r == '=' || r == '!' || r == '<' || r == '>': //
		return lexOperator
	case r == '/':
		// Special lookahead for '*' or '/', for comment check
		if int(l.pos) < len(l.input) {
			switch r := l.input[l.pos]; {
			case r == '/':
				return lexSinglelineComment
			case r == '*':
				return lexMultilineComment
			}
		}
		return lexOperator
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
	case r == '[':
		l.emit(tokenLeftSquare)
		l.squareDepth++
	case r == ']':
		l.emit(tokenRightSquare)
		l.squareDepth--
		if l.squareDepth < 0 {
			return l.errorf("Unexpected right square bracket %#U", r)
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
		return l.errorf("Unclosed left brace '{'")
	} else if l.squareDepth != 0 {
		return l.errorf("Unclosed left square bracket '['")
	}
	l.emit(tokenEOF)
	return nil
}

// lexSpace scans a run of space characters, One space has already been seen
// Ignore spaces seen
func lexSpace(l *lexer) stateFunc {
	for isSpace(l.peek()) {
		l.next()
	}
	l.ignore()
	return lexCode
}

// lexNewline scans for a run of newline characters (either \r\n OR \n)
func lexNewline(l *lexer) stateFunc {
Loop:
	for {
		switch r := l.next(); {
		case r == '\n':
			// Absorb and go to next iteration
		default:
			l.backup()
			break Loop
		}
	}
	// This is the automatic semicolon insertion for newlines, is inserted if the
	// token before the newline fits one of the following rules:
	// 1. the token is an identifier, or string/boolean/number literal
	// 2. the token is a `break`, `return` or `continue`
	// 3. token closes a bracket (either parenthesis, square brackets, or braces)
	switch l.prevTokTyp {
	case tokenIdentifier, tokenRawString, tokenQuotedString, tokenBool, tokenNumber, // identifiers and literals
		tokenBreak, tokenCont, tokenReturn, // keywords such as 'break', 'continue', 'return'
		tokenRightParan, tokenRightSquare, tokenRightBrace: // closing brackets ')', ']', '}'
		l.emit(tokenSemicolon)
	default:
		l.ignore()
	}
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

// lexOperator scans for a potential operator
// The first character ('+', '-', '/', '%', '*', '=', '!', '>', '<') has already
// been consumed
func lexOperator(l *lexer) stateFunc {
	r := l.input[int(l.start)] // store the 1st character somewhere
	if l.next() != '=' {
		l.backup() // go back to capture 'r' only
		switch r {
		case '+':
			l.emit(tokenPlus)
		case '-':
			l.emit(tokenMinus)
		case '/':
			l.emit(tokenDiv)
		case '%':
			l.emit(tokenMod)
		case '*':
			l.emit(tokenMult)
		case '=':
			l.emit(tokenAssign)
		case '!':
			l.emit(tokenLogicalNot)
		case '>':
			l.emit(tokenGreater)
		case '<':
			l.emit(tokenSmaller)
		}
	} else {
		// capture both r and the equal sign '='
		switch r {
		case '+':
			l.emit(tokenPlusAssign)
		case '-':
			l.emit(tokenMinusAssign)
		case '/':
			l.emit(tokenDivAssign)
		case '%':
			l.emit(tokenModAssign)
		case '*':
			l.emit(tokenMultAssign)
		case '=':
			l.emit(tokenEquals)
		case '!':
			l.emit(tokenNotEquals)
		case '>':
			l.emit(tokenGreaterEquals)
		case '<':
			l.emit(tokenSmallerEquals)
		}
	}
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
	return r == ' ' || r == '\t' || r == '\r'
}

func isEndOfLine(r rune) bool {
	return r == '\n'
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
