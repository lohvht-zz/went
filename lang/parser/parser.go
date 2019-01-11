package parser

import (
	"fmt"

	"github.com/lohvht/went/lang/ast"
	"github.com/lohvht/went/lang/token"
)

type Parser struct {
	name      string
	tokeniser *token.Lexer
	errors    token.ErrorList

	currentToken token.Token // next token to be consumed
	tokens       token.List  // lookahead tokens
}

func New(name, input string) (p *Parser) {
	p = &Parser{}
	eh := func(filename string, pos token.Pos, msg string) {
		p.errors.Add(filename, pos, msg)
		// NOTE: print to log for convenience, remove when no longer needed for debug
		// log.Fatalln(p.errors[len(p.errors)-1])
	}
	p.name = name
	p.tokeniser = token.NewLexer(name, input, eh)
	return
}

//===================================================================
// Parsing support

// errorf formats the message and its arguments and should be favoured over using p.error
func (p *Parser) errorf(pos token.Pos, message string, msgArgs ...interface{}) {
	p.errors.Add(p.name, pos, fmt.Sprintf(message, msgArgs...))
	// log.Fatalln(p.errors[len(p.errors)-1])
}

// next consumes and returns the next token
func (p *Parser) next() token.Token {
	// take a token from the bottom of the stack
	if !p.tokens.Empty() {
		p.currentToken = p.tokens.Shift()
	} else {
		p.currentToken = p.tokeniser.Scan()
	}
	return p.currentToken
}

// backup backs up a series of tokens to the bottom of the tokenList
// you should backup in the same order to preserve the proper token order from
// the token.Lexer (i.e. if given 3 tokens in this order: tkn1, tkn2, tkn3, you should
// call backup(tkn1, tkn2, tkn3))
func (p *Parser) backup(tkns ...token.Token) { p.tokens.Unshift(tkns...) }

// peek returns but does not consume the next token. If the there are no tokens left,
// grab one from the channel and add it into the tokens.
func (p *Parser) peek() token.Token {
	if !p.tokens.Empty() {
		return p.tokens.PeekBottom()
	}
	p.tokens.Push(p.tokeniser.Scan())
	return p.tokens.PeekBottom()
}

// match checks against the given list of token.Type. if the next token is
// of the same token.Type, it consumes it and return true.
func (p *Parser) match(types ...token.Type) bool {
	for _, typ := range types {
		if p.check(typ) {
			p.next()
			return true
		}
	}
	return false
}

// check returns true if the lookahead token matches the same type
func (p *Parser) check(typ token.Type) bool {
	tkn := p.peek()
	if tkn.Type == token.EOF {
		return false
	}
	return tkn.Type == typ
}

func (p *Parser) errorExpected(pos token.Pos, message string) {
	message = "expected " + message
	if pos == p.peek().Pos {
		// error happened at current position, make message more specific
		switch {
		case p.peek().Type == token.SEMICOLON && p.peek().Value == "\n":
			message += ", found newline"
		default:
			message += ", found '" + p.peek().Type.String() + "'"
		}
	}
	p.errorf(pos, message)
}

func (p *Parser) expect(typ token.Type) (token.Token, bool) {
	expected := p.check(typ)
	if !expected {
		p.errorExpected(p.peek().Pos, "'"+typ.String()+"'")
	}
	return p.next(), expected
}

func (p *Parser) sync() {
	for ; p.currentToken.Type != token.EOF; p.next() {
		switch p.currentToken.Type {
		case token.SEMICOLON: // end of expressions, discard semicolon and return
			p.next()
			return
		case token.CLASS, token.FUNC, token.VAR, // start of statements
			token.FOR, token.IF, token.WHILE, token.RETURN:
			return
		}
	}
}

//===================================================================
// Rules

func (p *Parser) Run() (expr ast.Expr, err error) {
	defer func() {
		if r := recover(); r != nil {
			err, _ = r.(error)
		}
	}()
	expr = p.expression()
	return
}

func (p *Parser) expression() ast.Expr {
	return p.equalityExpr()
}

func (p *Parser) equalityExpr() ast.Expr {
	expr := p.comparisonExpr()
	for p.match(token.EQ, token.NEQ) {
		op := p.currentToken
		r := p.comparisonExpr()
		expr = &ast.BinExpr{Left: expr, Op: op, Right: r}
	}
	return expr
}

func (p *Parser) comparisonExpr() ast.Expr {
	expr := p.addExpr()
	for p.match(token.SM, token.SMEQ, token.GR, token.GREQ) {
		op := p.currentToken
		r := p.addExpr()
		expr = &ast.BinExpr{Left: expr, Op: op, Right: r}
	}
	return expr
}

func (p *Parser) addExpr() ast.Expr {
	expr := p.multExpr()
	for p.match(token.PLUS, token.MINUS) {
		op := p.currentToken
		r := p.multExpr()
		expr = &ast.BinExpr{Left: expr, Op: op, Right: r}
	}
	return expr
}

func (p *Parser) multExpr() ast.Expr {
	expr := p.arithUnExpr()
	for p.match(token.MULT, token.DIV, token.MOD) {
		op := p.currentToken
		r := p.arithUnExpr()
		expr = &ast.BinExpr{Left: expr, Op: op, Right: r}
	}
	return expr
}

func (p *Parser) arithUnExpr() ast.Expr {
	if p.match(token.PLUS, token.MINUS) {
		op := p.currentToken
		operand := p.arithUnExpr()
		return &ast.UnExpr{Op: op, Operand: operand}
	}
	return p.primaryExpr()
}

func (p *Parser) primaryExpr() ast.Expr {
	var n ast.Expr
	switch {
	case p.match(token.FALSE):
		n = &ast.BasicLit{Value: false, Token: p.currentToken}
	case p.match(token.TRUE):
		n = &ast.BasicLit{Value: true, Token: p.currentToken}
	case p.match(token.NULL):
		n = &ast.BasicLit{Value: nil, Token: p.currentToken}
	case p.match(token.INT, token.FLOAT, token.STR):
		n = &ast.BasicLit{Value: p.currentToken.Value, Token: p.currentToken}
	case p.match(token.LROUND):
		lround := p.currentToken
		expr := p.expression()
		rround, ok := p.expect(token.RROUND)
		if !ok {
			panic(p.errors)
		}
		n = &ast.GrpExpr{LeftRound: lround, Expression: expr, RightRound: rround}
	}
	if n == nil {
		p.errorExpected(p.peek().Pos, "expression")
		panic(p.errors)
	}
	return n
}
