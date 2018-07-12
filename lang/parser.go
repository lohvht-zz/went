package lang

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
	Name         string
	Root         Node   // top-level root of the tree
	input        string // input text to be parsed
	tokeniser    *lexer
	tokens       tokenList // list of token lookaheads
	currentToken token     // the local that we are currently looking at (Not a lookahead)
}

// next consumes and returns the next token
func (p *Parser) next() token {
	// take a token from the bottom of the stack
	if !p.tokens.empty() {
		p.currentToken = p.tokens.shift()
	} else {
		p.currentToken = p.tokeniser.nextToken()
	}
	return p.currentToken
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
	format = fmt.Sprintf("Syntax Error: %s:%d: %s", p.Name, p.currentToken.line, format)
	panic(fmt.Errorf(format, args...))
}

// error terminates the processing.
func (p *Parser) error(err error) {
	p.errorf("%s", err)
}

// expect consumes the next token and guarantees it has the required type.
func (p *Parser) expect(context string, expected tokenType) token {
	tkn := p.next()
	if tkn.typ != expected {
		p.unexpected(context, tkn)
	}
	return tkn
}

// expectRange consumes the next token and guarantees it has one of the required types.
func (p *Parser) expectRange(context string, expectedTypes ...tokenType) (tkn token) {
	tkn = p.next()
	for _, exTyp := range expectedTypes {
		if tkn.typ == exTyp {
			return
		}
	}
	p.unexpected(context, tkn)
	return
}

// unexpected complains about the token and terminates processing
func (p *Parser) unexpected(context string, tkn token) {
	p.errorf("unexpected %s in %s", tkn, context)
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
	p.Root = p.expr()
	p.expect("End of File", tokenEOF)
}

// Grammar rules

// expr: term (("+" | "-") term)*;
func (p *Parser) expr() Node {
	node := p.term()
	for p.peek().typ == tokenPlus || p.peek().typ == tokenMinus {
		tkn := p.next()
		node = NewBinaryOp(node, tkn, p.term(), tkn.pos)
	}
	return node
}

// term: factor (("*" | "/" | "%") factor)*;
func (p *Parser) term() Node {
	node := p.factor()
	for p.peek().typ == tokenMult || p.peek().typ == tokenDiv || p.peek().typ == tokenMod {
		tkn := p.next()
		node = NewBinaryOp(node, tkn, p.factor(), tkn.pos)
	}
	return node
}

// factor: ("+" | "-") factor | atom;
func (p *Parser) factor() Node {
	switch p.peek().typ {
	case tokenPlus, tokenMinus:
		tkn := p.next()
		return NewUnaryOp(tkn, p.factor(), tkn.pos)
	default:
		return p.atom()
	}
}

// atom: ID | NUM | STR | RAWSTR | "null" | "false" | "true";
// Will complain and terminate if token is not an identifier, number, string,
// null or boolean.
func (p *Parser) atom() Node {
	tkn := p.expectRange("atom type check", tokenIdentifier, tokenNumber,
		tokenRawString, tokenQuotedString, tokenNull, tokenBool,
	)
	switch tkn.typ {
	case tokenNumber:
		number, err := NewNumber(tkn.pos, tkn.value)
		if err != nil {
			p.error(err)
		}
		return number
	default:
		p.unexpected("atom", tkn)
		return nil
	}
}
