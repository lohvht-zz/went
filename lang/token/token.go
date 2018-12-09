package token

import (
	"fmt"
	"strconv"
)

// Pos describes a source position via its line and col location, it is represented
// by concatenating two 32-bit integers representing line and col.
type Pos uint64

func newPos(line uint32, col uint32) Pos {
	return Pos(uint64(line)<<32 | uint64(col))
}

// decompose Pos into line and col
func (p Pos) decompose() (line int, col int) {
	line = int(p >> 32)
	col = int(0xffffffff & p)
	return
}

// String returns the string representation of the position line:col
func (p Pos) String() string {
	line, col := p.decompose()
	return fmt.Sprintf("%d:%d", line, col)
}

// Pos helpers

// AddOffset returns a new Pos by adding an offset to the col to a given Pos
func AddOffset(p Pos, offset int) Pos {
	line, newCol := p.decompose()
	newCol = newCol + offset
	if newCol < 0 {
		// if the offset reduces the col value to less than zero, we set to zero
		// NOTE: update if running into issues relating to debugging
		newCol = 0
	}
	return newPos(uint32(line), uint32(newCol))
}

// Token represents a Token of the Went programming language
// It holds the type the token, its value in terms of string val, as well
// as its position within the source input
type Token struct {
	Type
	Value string // value of this item
	Pos
}

// Tkn returns itself, to be used to provide a default implementation
// for embedding in a node. Embedded in all Nodes
func (tok Token) Tkn() Token { return tok }

func (tok Token) String() string {
	switch {
	case tok.Type == EOF:
		return "EOF"
	case tok.Type == ERROR:
		return fmt.Sprintf("<err: %s>", tok.Value)
	case tok.Type == SEMICOLON:
		return ";"
	case tok.Type == NAME:
		return fmt.Sprintf("<NAME:%q>", tok.Value)
	case keywordBegin+1 <= tok.Type && tok.Type < keywordEnd:
		return fmt.Sprintf("<%s>", tok.Value)
	}
	return fmt.Sprintf("%q", tok.Value)
}

// Type represents the type of the token
type Type int

// Types of tokens
const (
	ERROR Type = iota // error occurred; value is the text of error
	EOF

	DOT       // .
	COLON     // :
	SEMICOLON // ;
	COMMA     // ,

	LROUND  // (
	LCURLY  // {
	LSQUARE // [

	RROUND  // )
	RCURLY  // }
	RSQUARE // ]

	//Literal tokens (not including object, array)
	NAME
	INT   // Integer64
	FLOAT // float64 numbers
	STR   // Singly quoted ('\'') strings, escaped using a single '\' char

	operatorStart
	PLUS  // +
	MINUS // -
	DIV   // /
	MULT  // *
	MOD   // %

	ASSIGN      // =
	PLUSASSIGN  // +=
	MINUSASSIGN // -=
	DIVASSIGN   // /=
	MULTASSIGN  // *=
	MODASSIGN   // %=

	EQ   // ==, test for value equality
	NEQ  // !=, test for value inequality
	GR   // >, test for greater than
	SM   // <, test for smaller than
	GREQ // >=, test for greater than or equals to
	SMEQ // <=, test for smaller than or equals to

	LOGICALNOT // !
	LOGICALOR  // ||
	LOGICALAND // &&
	operatorEnd

	keywordBegin
	FUNC   // func keyword for function definition
	IF     // if keyword
	ELSE   // else keyword
	ELIF   // elif keyword
	FOR    // for keyword, for loops
	NULL   // null constant, treated as a keyword
	FALSE  // false constant, treated as a keyword
	TRUE   // True constant, treated as a keyword
	WHILE  // while keyword
	RETURN // return keyword
	IN     // in keyword
	BREAK  // break keyword
	CONT   // continue keyword
	VAR    // var keyword (variable declaration)
	keywordEnd
)

var tokenTypes = [...]string{
	ERROR:       "ERROR",
	EOF:         "EOF",
	DOT:         "DOT",
	COLON:       ":",
	SEMICOLON:   ";",
	COMMA:       ",",
	LROUND:      "(",
	LCURLY:      "{",
	LSQUARE:     "[",
	RROUND:      ")",
	RCURLY:      "}",
	RSQUARE:     "]",
	NAME:        "NAME",
	INT:         "INTEGER",
	FLOAT:       "FLOAT",
	STR:         "STRING",
	PLUS:        "+",
	MINUS:       "-",
	DIV:         "/",
	MULT:        "*",
	MOD:         "%",
	ASSIGN:      "=",
	PLUSASSIGN:  "+=",
	MINUSASSIGN: "-=",
	DIVASSIGN:   "/=",
	MULTASSIGN:  "*=",
	MODASSIGN:   "%=",
	EQ:          "==",
	NEQ:         "!=",
	GR:          ">",
	SM:          "<",
	GREQ:        ">=",
	SMEQ:        "<=",
	LOGICALNOT:  "!",
	LOGICALOR:   "||",
	LOGICALAND:  "&&",
	FUNC:        "func",
	IF:          "if",
	ELSE:        "else",
	ELIF:        "elif",
	FOR:         "for",
	NULL:        "null",
	FALSE:       "false",
	TRUE:        "true",
	WHILE:       "while",
	RETURN:      "return",
	IN:          "in",
	BREAK:       "break",
	CONT:        "continue",
	VAR:         "var",
}

func (t Type) String() string {
	s := ""
	if 0 <= t && t < Type(len(tokenTypes)) {
		s = tokenTypes[t]
	}
	if s == "" {
		s = fmt.Sprintf("token(%s)", strconv.Itoa(int(t)))
	}
	return s
}

var keywords map[string]Type

func init() {
	keywords = make(map[string]Type)
	for i := keywordBegin + 1; i < keywordEnd; i++ {
		keywords[tokenTypes[i]] = i
	}
}

// List is the stack of tokens the bottom of the stack is index 0, while
// top of stack is last index of the slice
type List []Token

// Empty checks if a token list is empty
func (tl *List) Empty() bool { return len(*tl) == 0 }

// Push a series of tokens in sequence to the top of the stack
func (tl *List) Push(tkns ...Token) { *tl = append(*tl, tkns...) }

// Pop removes a Token from the top of the stack, you should always check if
// the stack is empty prior to popping
func (tl *List) Pop() (tkn Token) {
	tkn, *tl = (*tl)[len(*tl)-1], (*tl)[:len(*tl)-1]
	return
}

// PeekTop looks at the top of the stack without consuming the Token, you should always
// check if the stack is empty prior to peeking
func (tl *List) PeekTop() Token {
	return (*tl)[len(*tl)-1]
}

// Unshift pushes a series of tokens to the bottom of the stack
func (tl *List) Unshift(tkns ...Token) {
	*tl = append(tkns, (*tl)...)
}

// Shift removes a Token from the bottom of the stack, you should always check if
// the stack is empty prior to shifting
func (tl *List) Shift() (tkn Token) {
	tkn, *tl = (*tl)[0], (*tl)[1:]
	return
}

// PeekBottom looks at the bottom of the stack without consuming the Token
// you should always check if the stack is empty prior to peeking
func (tl *List) PeekBottom() Token { return (*tl)[0] }
