package utils

import (
	"fmt"
	"runtime"
)

// tokenList is the stack of tokens the bottom of the stack is index 0, while
// top of stack is last index of the slice
type tokenList []token

func (tl *tokenList) empty() bool {
	return len(*tl) == 0
}

// push a series of tokens in sequence to the top of the stack
func (tl *tokenList) push(tkns ...token) {
	*tl = append(*tl, tkns...)
}

// pop removes a token from the top of the stack, you should always check if
// the stack is empty prior to popping
func (tl *tokenList) pop() (tkn token) {
	tkn, *tl = (*tl)[len(*tl)-1], (*tl)[:len(*tl)-1]
	return
}

// peekTop looks at the top of the stack without consuming the token, you should always
// check if the stack is empty prior to peeking
func (tl *tokenList) peekTop() token {
	return (*tl)[len(*tl)-1]
}

// unshift pushes a series of tokens to the bottom of the stack
func (tl *tokenList) unshift(tkns ...token) {
	*tl = append(tkns, (*tl)...)
}

// shift removes a token from the bottom of the stack, you should always check if
// the stack is empty prior to shifting
func (tl *tokenList) shift() (tkn token) {
	tkn, *tl = (*tl)[0], (*tl)[1:]
	return
}

// peekBottom looks at the bottom of the stack without consuming the token
// you should always check if the stack is empty prior to peeking
func (tl *tokenList) peekBottom() token {
	return (*tl)[0]
}

// Parser parses the input string (file or otherwise) and creates an AST at its Root
type Parser struct {
	Name      string
	Root      *Node  // top-level root of the tree
	input     string // input text to be parsed
	tokeniser *lexer
	tokens    tokenList // list of token lookaheads
	peekCount int
}

// next consumes and returns the next token
func (p *Parser) next() token {
	// take a token from the bottom of the stack
	if !p.tokens.empty() {
		return p.tokens.shift()
	}
	return p.tokeniser.nextToken()
}

// backup backs up a series of tokens to the bottom of the tokenList
// you should backup in the same order to preserve the proper token order from
// the lexer (i.e. if given 3 tokens in this order: tkn1, tkn2, tkn3, you should
// call backup(tkn1, tkn2, tkn3))
func (p *Parser) backup(tkns ...token) {
	p.tokens.unshift(tkns...)
}

// peek returns but does not consume the next token. If the there are no tokens left,
// grab one from the channel and add it into the tokens.
func (p *Parser) peek() token {
	if !p.tokens.empty() {
		return p.tokens.peekBottom()
	}
	p.tokens.push(p.tokeniser.nextToken())
	return p.tokens.peekBottom()
}

// Parsing

// errorf formats the error and terminates processing.
func (p *Parser) errorf(format string, args ...interface{}) {
	p.Root = nil
	format = fmt.Sprintf("Syntax Error: %s:%d: %s", p.Name, p.tokens[0].line, format)
	panic(fmt.Errorf(format, args...))
}

// error terminates the processing.
func (p *Parser) error(err error) {
	p.errorf("%s", err)
}

// recover is the handler that turns panics into returns from the top level
// of Parse
func (p *Parser) recover(errp *error) {
	e := recover()
	if e != nil {
		if _, ok := e.(runtime.Error); ok {
			panic(e)
		}
		if p != nil {
			p.tokeniser.drain()
			p.stopParse()
		}
		*errp = e.(error)
	}
}

// initParser initialises the parser, using the lexer
func initParser(tokeniser *lexer) *Parser {
	p := &Parser{Name: tokeniser.name, Root: nil, tokeniser: tokeniser, input: tokeniser.input}
	return p
}

func (p *Parser) stopParse() {
	p.tokeniser = nil
}

// Parse parses the input string to construct an AST
func Parse(name, input string) (parser *Parser, err error) {
	p := initParser(tokenise(name, input))
	defer p.recover(&err)
	p.parse()
	p.stopParse()
	return p, nil
}

func (p *Parser) parse() {

}
