package token

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Most of this package is listed and adapted from
// https://golang.org/src/text/template/parse/lex.go and partially taken from
// https://golang.org/src/go/scanner/scanner.go

// NewLexer creates a new lexer for the input string
func NewLexer(name, input string) *Lexer {
	l := &Lexer{
		Name:    name,
		Input:   input,
		line:    1,
		col:     0,
		prevCol: 0,
	}
	return l
}

// Lexer scans the entire input string and tokenises it, storing the tokens in
// a channel of Tokens
type Lexer struct {
	Name  string // name of the input; used only for error reporting
	Input string // string being scanned

	// current state to track & emit info
	line    uint32 // 1 + number of newlines seen
	col     uint32 // 1 + current column number
	prevCol uint32 // previous column number seen (ensure backup() is correct)

	// Internal lexer state
	start        int       // start position of the current token
	pos          int       // current position
	runeWidth    int       // runeWidth of the last rune read from input
	prevTokTyp   Type      // previous Token type used for automatic semicolon insertion
	bracketStack runeStack // a stack of runes used to keep track of all '(', '[' and '{'
}

const eof = -1

type runeStack []rune

func (rs *runeStack) empty() bool { return len(*rs) == 0 }

// push a rune to the top of the stack
func (rs *runeStack) push(r rune) { *rs = append(*rs, r) }

// pop removes a rune from the top of the stack, you should always check if
// the stack is empty prior to popping
func (rs *runeStack) pop() (r rune) {
	r, *rs = (*rs)[len(*rs)-1], (*rs)[:len(*rs)-1]
	return
}

// peek looks at the top of the stack you should always check if the stack is
// empty prior to peeking
func (rs *runeStack) peek() rune { return (*rs)[len(*rs)-1] }

// next returns the next rune in the input
// next increases newline count
func (l *Lexer) next() rune {
	if int(l.pos) >= len(l.Input) {
		l.runeWidth = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.Input[l.pos:])
	l.runeWidth = w
	l.pos += l.runeWidth
	// handle columns and lines seen
	if r == '\n' {
		l.line++
		l.col = 1
	} else {
		l.prevCol = l.col
		l.col++
	}
	return r
}

// peek returns but does not consume next rune in the input
func (l *Lexer) peek() rune {
	if l.pos >= len(l.Input) {
		return eof
	}
	r, _ := utf8.DecodeRuneInString(l.Input[l.pos:])
	return r
}

// backup steps back one rune, can only be called once per call of next
func (l *Lexer) backup() {
	l.pos -= l.runeWidth
	l.col = l.prevCol
	if l.runeWidth == 1 && l.Input[l.pos] == '\n' {
		l.line--
	}
}

// nextToken returns the next token at the lexer's current position
// this will also update the last seen emitted Token type
func (l *Lexer) nextToken(typ Type) Token {
	tkn := Token{typ, l.Input[l.start:l.pos], newPos(l.line, l.col)}
	l.start = l.pos
	l.prevTokTyp = typ
	return tkn
}

// ignore skips over the pending input before this point
func (l *Lexer) ignore() { l.start = l.pos }

// accept consumes the next rune if its from the valid set
func (l *Lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set
func (l *Lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

// errorf emits an error Token returns an error token
// TODO: Make it such that lexical analysis is not terminated when error is reached
// to improve usability (EXTENSION WITH ERRORLISTS)
func (l *Lexer) errorf(format string, args ...interface{}) Token {
	return Token{ILLEGAL, fmt.Sprintf(format, args...), newPos(l.line, l.col)}
}

func (l *Lexer) scan2(runeToScan rune, typ0, typ1 Type) Token {
	if l.peek() != runeToScan {
		return l.nextToken(typ1)
	}
	l.next()
	return l.nextToken(typ0)
}

// Scan scans for the next token and returns it (Type, string Val and Pos in string)
// end of source is indicated by a Token of Type EOF.
//
func (l *Lexer) Scan() Token {
ScanAgain:
	l.skipWhitespace()

	switch r := l.next(); {
	case isLetter(r):
		l.backup()
		return l.lexIdentifier()
	case '0' <= r && r <= '9':
		return l.lexNumber()
	case r == eof:
		if !l.bracketStack.empty() {
			r := l.bracketStack.pop()
			// TODO: Rewrite errorhandling code
			l.errorf("unclosed left bracket: %#U", r)
		}
		return l.nextToken(EOF)
	case r == '\n':
		insertSemicolon := false
		l.skipNewlines(&insertSemicolon)
		if insertSemicolon {
			return l.nextToken(SEMICOLON)
		}
		goto ScanAgain
	case r == '\'':
		return l.lexQuotedString()
	case r == '`':
		return l.lexRawString()
	case r == ':':
		return l.nextToken(COLON)
	case r == '.':
		if r := l.peek(); r < '0' || r > '9' { // if its not a number
			return l.nextToken(DOT)
		}
		return l.lexNumber()
	case r == ',':
		return l.nextToken(COMMA)
	case r == ';':
		return l.nextToken(SEMICOLON)
	case r == '(':
		l.bracketStack.push('(')
		return l.nextToken(LROUND)
	case r == ')':
		// TODO: rewrite error handling
		if l.bracketStack.empty() {
			l.errorf("unexpected right bracket %#U", r)
		} else if toCheck := l.bracketStack.pop(); toCheck != '(' {
			l.errorf("unexpected right bracket %#U", r)
		}
		return l.nextToken(RROUND)
	case r == '[':
		l.bracketStack.push('[')
		return l.nextToken(LSQUARE)
	case r == ']':
		// TODO: rewrite error handling
		if l.bracketStack.empty() {
			l.errorf("unexpected right bracket %#U", r)
		} else if toCheck := l.bracketStack.pop(); toCheck != '[' {
			l.errorf("unexpected right bracket %#U", r)
		}
		return l.nextToken(RSQUARE)
	case r == '{':
		l.bracketStack.push('{')
		return l.nextToken(LCURLY)
	case r == '}':
		switch {
		// TODO: rewrite error
		case l.bracketStack.empty():
			l.errorf("unexpected right bracket %#U", r)
		case l.bracketStack.pop() != '{':
			l.errorf("unexpected right bracket %#U", r)
		case l.prevTokTyp != SEMICOLON:
			return l.nextToken(SEMICOLON)
		}
		return l.nextToken(RCURLY)
	case r == '|':
		if l.peek() != '|' {
			// TODO: Rewrite errorhandling code
			l.errorf("Unexpected token: %#U", r)
		}
		l.next()
		return l.nextToken(LOGICALOR)
	case r == '&':
		if l.peek() != '&' {
			// TODO: Rewrite errorhandling code
			l.errorf("Unexpected token: %#U", r)
		}
		l.next()
		return l.nextToken(LOGICALAND)
	case r == '+':
		return l.scan2('=', PLUS, PLUSASSIGN)
	case r == '-':
		return l.scan2('=', MINUS, MINUSASSIGN)
	case r == '*':
		return l.scan2('=', MULT, MULTASSIGN)
	case r == '%':
		return l.scan2('=', MOD, MODASSIGN)
	case r == '=':
		return l.scan2('=', ASSIGN, EQ)
	case r == '!':
		return l.scan2('=', LOGICALNOT, NEQ)
	case r == '<':
		return l.scan2('=', SM, SMEQ)
	case r == '>':
		return l.scan2('=', GR, GREQ)
	case r == '/':
		// handle for '/', can be comment or divide sign
		switch r := l.peek(); {
		case r == '/':
			l.skipSingleLineComment()
		case r == '*':
			l.skipMultilineComment()
		default:
			return l.scan2('=', DIV, DIVASSIGN)
		}
		goto ScanAgain
	default:
		l.errorf("illegal character: %#U", r)
		return l.nextToken(ILLEGAL)
	}

}

// lexQuotedString scans a quoted string, can be escaped using the '\' character
func (l *Lexer) lexQuotedString() Token {
	l.ignore() // ignore the opening quote
Loop:
	for {
		switch l.next() {
		case '\\': // single '\' character as escape character
			if r := l.next(); r == '\n' || r == eof {
				// TODO: Rewrite error handling
				l.errorf("unterminated quoted string")
			}
		case '\'':
			l.backup() // move back before the closing quote
			break Loop
		}
	}
	tkn := l.nextToken(STR)
	l.next()
	l.ignore() // now consume and ignore the closing quote
	return tkn
}

// lexRawString scans a raw string delimited by '`' character
func (l *Lexer) lexRawString() Token {
	l.ignore() // ignore the opening quote
	startLine := l.line
	startCol := l.col
Loop:
	for {
		switch l.next() {
		case eof:
			// restore line and col number to the location of the opening quote
			// will error out, okay to overwrite l.line
			l.line = startLine
			l.col = startCol
			l.errorf("Unterminated raw string")
		case '`':
			l.backup() // move back before the closing quote
			break Loop
		}
	}
	tkn := l.nextToken(STR)
	l.next()
	l.ignore() // now consume and ignore the closing quote
	return tkn
}

// scanSignificand scans for all numbers (of the given base) up to a non-number
func (l *Lexer) scanSignificand(base int) {
	for digitValue(l.next()) < base {
	}
	l.backup()
}

// lexNumber scans for a number, assumes that the lexer has not consumed the start
// of the number (either number or a dot)
func (l *Lexer) lexNumber() Token {
	l.backup() // backup to see the '.' or numerical runes
	emitTyp := INT
	// Seen decimal point --> is a float (i.e. .1234E10 for example)
	if l.peek() == '.' {
		goto FRACTION
	}
	// Leading 0 ==> hexadecimal ("0x"/"0X") or octal 0
	// if l.peek() == '0' {
	if l.accept("0") {
		if l.accept("xX") {
			// hexadecimal int
			l.scanSignificand(16)
			if l.pos-l.start <= 2 {
				// Only scanned "0x" or "0X"
				// TODO: Rewrite error handling
				l.errorf("illegal hexadecimal number: %q", l.Input[l.start:l.pos])
			}
		} else {
			l.scanSignificand(8)
			if l.accept("89") {
				// error, illegal octal int/float
				l.scanSignificand(10)
				// TODO: Rewrite error handling
				l.errorf("illegal octal number: %q", l.Input[l.start:l.pos])
			}
			if r := l.peek(); r == '.' || r == 'e' || r == 'E' {
				// NOTE: ".eEi" including imaginary number, if we wanna support it in the future
				// Octal float
				goto FRACTION
			}
		}
		return l.nextToken(emitTyp)
	}
	// Decimal integer/float
	l.scanSignificand(10)
FRACTION: // handles all other floating point lexing
	if l.accept(".") {
		emitTyp = FLOAT
		if r := l.peek(); !(r >= '0' && r <= '9') {
			// NOTE: we prohibit trailing decimal points with no numbers as we would
			// eventually support method overloading for numbers etc.
			// TODO: Rewrite error handling
			l.errorf("Illegal trailing decimal point after number")
		}
		l.scanSignificand(10)
	}
	if l.accept("eE") {
		emitTyp = FLOAT
		l.accept("+-")
		if digitValue(l.peek()) < 10 {
			l.scanSignificand(10)
		} else {
			// TODO: Rewrite error handling
			l.errorf("Illegal floating-point exponent: %q", l.Input[l.start:l.pos])
		}
	}
	return l.nextToken(emitTyp)
}

// lexIdentifier scans an alphanumeric word
func (l *Lexer) lexIdentifier() Token {
	r := l.next()
	for isLetter(r) || isDigit(r) {
		r = l.next()
	}
	l.backup()
	word := l.Input[l.start:l.pos]
	var typ Type
	if keywordBegin+1 <= keywords[word] && keywords[word] < keywordEnd {
		typ = keywords[word]
	} else {
		typ = NAME
	}
	return l.nextToken(typ)
}

func (l *Lexer) skipWhitespace() {
	for isSpace(l.next()) {
	}
	l.backup()
	l.ignore()
}

// skipNewlines ignores consecutive newlines and sets the state for
// automatic semicolon insertion via the following rules:
// 1. the Token is an identifier, or string/number literal
// 2. the Token is a `break`, `return` or `continue`
// 3. Token closes a round, square, or curly bracket (')', ']', '}')
func (l *Lexer) skipNewlines(insertSemicolon *bool) {
	l.ignore() // ignore the 1st newline
	for r := l.next(); r != '\n'; r = l.next() {
		// advance head of the lexer, go to next iteration
		// We do this ignore bit here so that whenever we emit a semicolon,
		// the string literal emitted will be condensed to a single \n
		l.ignore()
	}
	l.backup()
	switch l.prevTokTyp {
	case NAME, STR, INT, FLOAT,
		BREAK, CONT, RETURN,
		RROUND, RSQUARE, RCURLY:
		*insertSemicolon = true
	default:
		l.ignore() // do not count the spaces as the next() already adds
	}
}

// skipSingleLineComment skips over the while single line comment
func (l *Lexer) skipSingleLineComment() {
	for r := l.next(); !(r == '\n' || r == eof); r = l.next() {
	}
	l.ignore()
}

// skipMultilineComment skips over the whole multiline comment
// The left comment marker ('/*') has already been consumed
func (l *Lexer) skipMultilineComment() {
	// NOTE: we ignore this as unclosed multiline comments shouldnt be an error
	// if i := strings.Index(l.Input[l.pos:], "*/"); i < 0 {
	// 	// TODO: rewrite error handling
	// 	l.errorf("Multiline comment is not closed")
	// }
	var left, right rune
	right = l.next()
	for {
		left, right = right, l.next()
		if left == '*' && right == '/' {
			break
		} else if left == eof || right == eof {
			break
		}
	}
	l.ignore()
}

// Utility Functions

func isSpace(r rune) bool { return r == ' ' || r == '\t' || r == '\r' }

func isDigit(r rune) bool {
	return '0' <= r && r <= '9' || r >= utf8.RuneSelf && unicode.IsDigit(r)
}

func isLetter(r rune) bool {
	return 'a' <= r && r <= 'z' ||
		'A' <= r && r <= 'Z' ||
		r == '_' ||
		r >= utf8.RuneSelf && unicode.IsLetter(r)
}

func digitValue(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= ch && ch <= 'f':
		return int(ch - 'a' + 10)
	case 'A' <= ch && ch <= 'F':
		return int(ch - 'A' + 10)
	}
	return 16 // larger than any legal digit val
}
