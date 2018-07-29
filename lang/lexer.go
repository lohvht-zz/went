package lang

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
	line  LinePos   // Line number at the start of this item
}

func (tok token) String() string {
	switch {
	case tok.typ == tokenEOF:
		return "EOF"
	case tok.typ == tokenError:
		return fmt.Sprintf("<err: %s>", tok.value)
	case tok.typ == tokenSemicolon:
		return ";"
	case tok.typ > tokenKeyword:
		return fmt.Sprintf("<%s>", tok.value)
	}
	return fmt.Sprintf("%q", tok.value)
}

type tokenType int

const (
	tokenError tokenType = iota // error occurred; value is the text of error
	tokenEOF
	tokenDot         // Dot character '.'
	tokenIdentifier  // alphanumeric identifier
	tokenLeftRound   // left bracket '('
	tokenRightRound  // right bracket ')'
	tokenLeftCurly   // left curly bracket '{'
	tokenRightCurly  // right curly bracket '}'
	tokenLeftSquare  // left square bracket '['
	tokenRightSquare // right square bracket ']'
	tokenColon       // colon symbol ':'
	tokenSemicolon   // semi colon symbol ';'
	tokenComma       // comma symbol ','

	// Literal tokens (not including object, array)
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
	tokenFalse   // 'false' constant, treated as a keyword
	tokenTrue    // 'True' constant, treated as a keyword
	tokenWhile   // 'while' keyword
	tokenReturn  // 'return' keyword
	tokenIn      // 'in' keyword
	tokenBreak   // 'break' keyword
	tokenCont    // 'continue' keyword
)

var tokenNames = map[tokenType]string{
	tokenError:       "error",
	tokenEOF:         "EOF",
	tokenIdentifier:  "identifier",
	tokenLeftRound:   "(",
	tokenRightRound:  ")",
	tokenLeftCurly:   "{",
	tokenRightCurly:  "}",
	tokenLeftSquare:  "[",
	tokenRightSquare: "]",
	tokenColon:       ":",

	// Literal tokens (not including object, array)
	tokenNumber:       "number",
	tokenQuotedString: "string",
	tokenRawString:    "raw string",

	// Arithmetic Operators
	tokenPlus:  "+",
	tokenMinus: "-",
	tokenDiv:   "/",
	tokenMult:  "*",
	tokenMod:   "%",
	// Assignment Operators
	tokenAssign:      "=",
	tokenPlusAssign:  "+=",
	tokenMinusAssign: "-=",
	tokenDivAssign:   "/=",
	tokenMultAssign:  "*=",
	tokenModAssign:   "%=",
	// Comparison Operators
	tokenEquals:        "==",
	tokenNotEquals:     "!=",
	tokenGreater:       ">",
	tokenSmaller:       "<",
	tokenGreaterEquals: ">=",
	tokenSmallerEquals: "<=",
	// Logical Operators
	tokenLogicalNot: "!",
	tokenOr:         "||",
	tokenAnd:        "&&",

	// Keywords after all the rest
	tokenFunc:   "func",
	tokenIf:     "if",
	tokenElse:   "else",
	tokenElseIf: "elif",
	tokenFor:    "for",
	tokenNull:   "null",
	tokenFalse:  "false",
	tokenTrue:   "true",
	tokenWhile:  "while",
	tokenReturn: "return",
	tokenIn:     "in",
}

func (i tokenType) String() string {
	s := tokenNames[i]
	if s == "" {
		return fmt.Sprintf("token%d", int(i))
	}
	return s
}

var keyMap = map[string]tokenType{
	"func":     tokenFunc,
	"if":       tokenIf,
	"else":     tokenElse,
	"elif":     tokenElseIf,
	"for":      tokenFor,
	"null":     tokenNull,
	"false":    tokenFalse,
	"true":     tokenTrue,
	"while":    tokenWhile,
	"return":   tokenReturn,
	"in":       tokenIn,
	"break":    tokenBreak,
	"continue": tokenCont,
}

var parenMap = map[rune]rune{
	')': '(',
	']': '[',
	'}': '{',
}

const eof = -1

/**
 * lexer Definition
 */

type runeStack []rune

func (rs *runeStack) empty() bool {
	return len(*rs) == 0
}

// push a rune to the top of the stack
func (rs *runeStack) push(r rune) {
	*rs = append(*rs, r)
}

// pop removes a rune from the top of the stack, you should always check if
// the stack is empty prior to popping
func (rs *runeStack) pop() (r rune) {
	r, *rs = (*rs)[len(*rs)-1], (*rs)[:len(*rs)-1]
	return
}

// peek looks at the top of the stack you should always check if the stack is
// empty prior to peeking
func (rs *runeStack) peek() rune {
	return (*rs)[len(*rs)-1]
}

type lexer struct {
	name         string     // name of the input; used only for error reporting
	input        string     // string being scanned
	pos          Pos        // current position
	start        Pos        // start position of this token
	width        Pos        // width of the last rune read from input
	tokens       chan token // channel of the scanned items
	prevTokTyp   tokenType  // previous token type used for automatic semicolon insertion
	bracketStack runeStack  // a stack of runes used to keep track of all '(', '[' and '{'
	line         LinePos    // 1 + number of newlines seen
}

// next returns the next rune in the input
// next increases newline count
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
	case tokenRawString:
		l.line += LinePos(strings.Count(l.input[l.start:l.pos], "\n"))
	}
	l.start = l.pos
	l.prevTokTyp = typ
}

// ignore skips over the pending input before this point
func (l *lexer) ignore(countSpace bool) {
	if countSpace {
		l.line += LinePos(strings.Count(l.input[l.start:l.pos], "\n"))
	}
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
	// Check if the next rune is alphanumeric
	// The next number can't be digits as we have already scanned all the digits
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
		eof, '=', // EOF character and assignment/declaration ('='), or equality check ('==')
		'.', ',', // DOT ('.') to denote .property, or commas
		'|', '&', // OR ('||'), or AND ('&&')
		'(', ')', '[', ']', '{', '}', // Parenthesis, square, curly and normal
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
	case isEndOfLine(r): // detects \r or \n
		l.backup()
		return lexNewline
	case r == ':':
		l.emit(tokenColon)
	case r == ',':
		l.emit(tokenComma)
	case r == '|':
		if l.next() == '|' {
			l.emit(tokenOr)
		} else {
			l.errorf("expected token %#U", r)
		}
	case r == '&':
		if l.next() == '&' {
			l.emit(tokenAnd)
		} else {
			l.errorf("expected token %#U", r)
		}
	case r == '\'':
		l.ignore(false) // ignore the opening quote
		return lexQuotedString
	case r == '`':
		l.ignore(false) // ignore the opening quote
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
	case r == '+', r == '-', r == '*', r == '%', // Math signs
		r == '=', r == '!', r == '<', r == '>': //
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
	case r == '(', r == '{', r == '[': // opening brackets
		switch r {
		case '(':
			l.emit(tokenLeftRound)
		case '{':
			l.emit(tokenLeftCurly)
		case '[':
			l.emit(tokenLeftSquare)
		}
		l.bracketStack.push(r)
	case r == ')', r == '}', r == ']':
		if l.bracketStack.empty() {
			return l.errorf("unexpected right paren %#U", r)
		} else if toCheck := l.bracketStack.pop(); toCheck != parenMap[r] {
			return l.errorf("unexpected right paren %#U", r)
		}
		switch r {
		case ')':
			l.emit(tokenRightRound)
		case '}':
			l.emit(tokenRightCurly)
		case ']':
			l.emit(tokenRightSquare)
		}
	default:
		return l.errorf("Unrecognised character in code: %#U", r)
	}
	return lexCode
}

// lexEOF emits the EOF token and handles the termination of the main lexCode loop
func lexEOF(l *lexer) stateFunc {
	if !l.bracketStack.empty() {
		r := l.bracketStack.pop()
		return l.errorf("unclosed left paren: %#U", r)
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
	l.ignore(false)
	return lexCode
}

// lexNewline scans for a run of newline chars ('\n')
// This method also does the automatic semicolon insertions with the following
// rules:
// 1. the token is an identifier, or string/boolean/number literal
// 2. the token is a `break`, `return` or `continue`
// 3. token closes a round or square bracket (')', ']')
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
	switch l.prevTokTyp {
	case tokenIdentifier, tokenRawString, tokenQuotedString, tokenFalse,
		tokenTrue, tokenNumber, tokenBreak, tokenCont, tokenReturn,
		tokenRightRound, tokenRightSquare:
		l.emit(tokenSemicolon)
	default:
		l.ignore(false) // do not count the spaces as the next() already adds
	}
	return lexCode
}

// lexQuotedString scans a quoted string, can be escaped using the '\' character
func lexQuotedString(l *lexer) stateFunc {
Loop:
	for {
		switch l.next() {
		case '\\': // single '\' character as escape character
			if r := l.next(); r == '\n' || r == eof {
				return l.errorf("unterminated quoted string")
			}
		case '\'':
			l.backup() // move back before the closing quote
			break Loop
		}
	}
	l.emit(tokenQuotedString)
	l.next()
	l.ignore(false) // now consume and ignore the closing quote
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
			l.backup() // move back before the closing quote
			break Loop
		}
	}
	l.emit(tokenRawString)
	l.next()
	l.ignore(false) // now consume and ignore the closing quote
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
	l.ignore(false)
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
	l.ignore(true)
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
