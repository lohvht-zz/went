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
	p.Root = p.orEval()
	p.expect("End of File", tokenEOF)
}

// Grammar rules

// orEval: andEval ("||" orEval)*;
func (p *Parser) orEval() Node {
	node := p.notEval()
	for p.peek().typ == tokenOr {
		tkn := p.next()
		node = newOr(node, p.notEval(), tkn.pos)
	}
	return node
}

// andEval: notEval ("&&" notEval)*;
func (p *Parser) andEval() Node {
	node := p.notEval()
	for p.peek().typ == tokenAnd {
		tkn := p.next()
		node = newAnd(node, p.notEval(), tkn.pos)
	}
	return node
}

// notEval: "!" notEval | comparison;
func (p *Parser) notEval() Node {
	switch p.peek().typ {
	case tokenLogicalNot:
		tkn := p.next()
		return newNot(p.notEval(), tkn.pos)
	default:
		return p.comparison()
	}
}

// comparison: smExpr (compOp smExpr)*;
// compOp: compOp: "==" | "!=" | "<" | ">" | "<=" | ">=" | ["!"] "in";
func (p *Parser) comparison() Node {
	node := p.smExpr()
Loop:
	for {
		switch p.peek().typ {
		case tokenEquals, tokenNotEquals:
			tkn := p.next()
			node = newEq(node, p.smExpr(), tkn.pos, tkn.typ == tokenNotEquals)
		case tokenSmaller, tokenSmallerEquals:
			tkn := p.next()
			node = newSm(node, p.smExpr(), tkn.pos, tkn.typ == tokenSmallerEquals)
		case tokenGreater, tokenGreaterEquals:
			tkn := p.next()
			node = newGr(node, p.smExpr(), tkn.pos, tkn.typ == tokenGreaterEquals)
		case tokenIn:
			tkn := p.next()
			node = newIn(node, p.smExpr(), tkn.pos, false)
		case tokenLogicalNot:
			tkn := p.next()
			// Expects the next token to be tokenIn
			p.expect("not operator in comparison rule", tokenIn)
			node = newIn(node, p.smExpr(), tkn.pos, true)
		default:
			break Loop
		}
	}
	return node
}

// smExpr: term (("+" | "-") term)*;
func (p *Parser) smExpr() Node {
	node := p.term()
Loop:
	for {
		switch p.peek().typ {
		case tokenPlus:
			tkn := p.next()
			node = newAdd(node, p.term(), tkn.pos)
		case tokenMinus:
			tkn := p.next()
			node = newSubtract(node, p.term(), tkn.pos)
		default:
			break Loop
		}
	}
	return node
}

// term: factor (("*" | "/" | "%") factor)*;
func (p *Parser) term() Node {
	node := p.factor()
Loop:
	for {
		switch p.peek().typ {
		case tokenMult:
			tkn := p.next()
			node = newMult(node, p.factor(), tkn.pos)
		case tokenDiv:
			tkn := p.next()
			node = newDiv(node, p.factor(), tkn.pos)
		case tokenMod:
			tkn := p.next()
			node = newMod(node, p.factor(), tkn.pos)
		default:
			break Loop
		}
	}
	return node
}

// factor: ("+" | "-") factor | atom;
func (p *Parser) factor() Node {
	switch p.peek().typ {
	case tokenPlus:
		tkn := p.next()
		return newPlus(p.factor(), tkn.pos)
	case tokenMinus:
		tkn := p.next()
		return newMinus(p.factor(), tkn.pos)
	default:
		return p.atom()
	}
}

// atom: "[" [exprList] "]" | "expr" | ID | NUM | STR | RAWSTR | "null" | "false" | "true";
// Will complain and terminate if token is not an identifier, number, string,
// null or boolean.
func (p *Parser) atom() Node {
	tkn := p.expectRange("atom type check", tokenIdentifier, tokenNumber,
		tokenRawString, tokenQuotedString, tokenNull, tokenFalse, tokenTrue,
		tokenLeftParen, tokenLeftSquare,
	)
	switch tkn.typ {
	case tokenNumber:
		n, err := newNumber(tkn.pos, tkn.value)
		if err != nil {
			p.error(err)
		}
		return n
	case tokenRawString, tokenQuotedString:
		return newString(tkn.pos, tkn.value)
	case tokenNull:
		return newNull(tkn.pos, tkn.value)
	case tokenFalse, tokenTrue:
		n, err := newBool(tkn.pos, tkn.value, tkn.typ)
		if err != nil {
			p.error(err)
		}
		return n
	case tokenLeftParen:
		n := p.orEval()
		p.expect("closing parenthesis", tokenRightParan)
		return n
	// case tokenLeftSquare:

	default:
		p.unexpected("atom", tkn)
		return nil
	}
}
