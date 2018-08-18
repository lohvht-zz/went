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
// also links the AST to the appropriate scopes
type Parser struct {
	Name string
	Root Node // top-level root of the AST tree
	// symtab *SymbolTable // the entire symbol table, global scope, local scope, functions etc.
	// currentScope *Scope
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
	format = fmt.Sprintf("%s line %d: SyntaxError - %s", p.Name, p.currentToken.line, format)
	panic(fmt.Errorf(format, args...))
}

// error terminates the processing.
func (p *Parser) error(err error) { p.errorf("%s", err) }

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
	p := &Parser{Name: tokeniser.name, Root: nil, tokeniser: tokeniser,
		input: tokeniser.input}
	return p
}

func (p *Parser) stopParse() { p.tokeniser = nil }

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
	if p.peek().typ == tokenSemicolon {
		p.next() // just consume the semicolon for now
	}
	p.expect("End of File", tokenEOF)
}

// Grammar rules

// func (p *Parser) input() Node {
// 	for p.peek().typ != tokenEOF {

// 	}
// }

// func (p *Parser) stmt() Node {
// 	for
// }

// // exprStmt: exprList (augAssign exprList | ('=' exprList)*);
// // augAssign: "+=" | "-=" | "/=" | "*=" | "%=";
// func (p *Parser) exprStmt() Stmt {
// 	// get the first token, this used to give the line position and position
// 	// of the start of the statement
// 	firstTkn := p.peek()
// 	exprs := p.exprList()
// 	switch tkntyp := p.peek().typ; tkntyp {
// 	case tokenPlusAssign, tokenMinusAssign, tokenDivAssign, tokenMultAssign,
// 		tokenModAssign, tokenAssign:
// 		return p.assignStmt(exprs, tkntyp)
// 	default:
// 		// // TODO: Do not accept if its not an assignment statement and it has more
// 		// // than 1 expression to evaluate
// 		// if len(exprs) > 1 {
// 		// 	err := fmt.Errorf("cannot have more than 1 expression")
// 		// }
// 		return newExprStmt(exprs, firstTkn)
// 	}
// }

// func (p *Parser) assignStmt(lhs []Expr, typ tokenType) Stmt {
// 	// Possible errors:
// 	// mismatched LHS/RHS number of values
// 	// LHS is not addressable (i.e. not a NAME/SLICE)
// 	for _, lhExpr := range lhs {
// 		switch expr := lhExpr.(type) {
// 		case ID:
// 		default:
// 			// Error, not an addressable
// 			p.errorf()
// 		}
// 	}

// 	rhs := p.exprList()

// }

// orEval: andEval ("||" orEval)*;
func (p *Parser) orEval() Expr {
	node := p.andEval()
	for p.peek().typ == tokenOr {
		tkn := p.next()
		node = newOr(node, p.andEval(), tkn)
	}
	return node
}

// andEval: notEval ("&&" notEval)*;
func (p *Parser) andEval() Expr {
	node := p.notEval()
	for p.peek().typ == tokenAnd {
		tkn := p.next()
		node = newAnd(node, p.notEval(), tkn)
	}
	return node
}

// notEval: "!" notEval | comparison;
func (p *Parser) notEval() Expr {
	switch p.peek().typ {
	case tokenLogicalNot:
		tkn := p.next()
		return newNot(p.notEval(), tkn)
	default:
		return p.comparison()
	}
}

// comparison: smExpr (compOp smExpr)*;
// compOp: compOp: "==" | "!=" | "<" | ">" | "<=" | ">=" | ["!"] "in";
func (p *Parser) comparison() Expr {
	node := p.smExpr()
Loop:
	for {
		switch p.peek().typ {
		case tokenEquals, tokenNotEquals:
			tkn := p.next()
			node = newEq(node, p.smExpr(), tkn.typ == tokenNotEquals, tkn)
		case tokenSmaller, tokenSmallerEquals:
			tkn := p.next()
			node = newSm(node, p.smExpr(), tkn.typ == tokenSmallerEquals, tkn)
		case tokenGreater, tokenGreaterEquals:
			tkn := p.next()
			node = newGr(node, p.smExpr(), tkn.typ == tokenGreaterEquals, tkn)
		case tokenIn:
			tkn := p.next()
			node = newIn(node, p.smExpr(), tkn)
		default:
			break Loop
		}
	}
	return node
}

// smExpr: term (("+" | "-") term)*;
func (p *Parser) smExpr() Expr {
	node := p.term()
Loop:
	for {
		switch p.peek().typ {
		case tokenPlus:
			tkn := p.next()
			node = newAdd(node, p.term(), tkn)
		case tokenMinus:
			tkn := p.next()
			node = newSubtract(node, p.term(), tkn)
		default:
			break Loop
		}
	}
	return node
}

// term: factor (("*" | "/" | "%") factor)*;
func (p *Parser) term() Expr {
	node := p.factor()
Loop:
	for {
		switch p.peek().typ {
		case tokenMult:
			tkn := p.next()
			node = newMult(node, p.factor(), tkn)
		case tokenDiv:
			tkn := p.next()
			node = newDiv(node, p.factor(), tkn)
		case tokenMod:
			tkn := p.next()
			node = newMod(node, p.factor(), tkn)
		default:
			break Loop
		}
	}
	return node
}

// factor: ("+" | "-") factor | atom;
func (p *Parser) factor() Expr {
	switch p.peek().typ {
	case tokenPlus:
		tkn := p.next()
		return newPlus(p.factor(), tkn)
	case tokenMinus:
		tkn := p.next()
		return newMinus(p.factor(), tkn)
	default:
		return p.atom()
	}
}

// TODO: Implement me!
func (p *Parser) atomExpr() Expr {
	n := p.atom()
	switch p.peek().typ {
	case tokenDot: // map reference
		return nil
	case tokenLeftRound: // function call
		return nil
	case tokenLeftSquare: // slice / index reference
		return nil
	default:
		return n
	}
}

// atom: "[" [exprList] "]" | "{" mapList "}" | "(" expr ")" | ID | NUM | STR |
// RAWSTR | "null" | "false" | "true";
// mapList: keyval ("," keyval)* [","];
// keyval: (ID | STR) ":" expr;
func (p *Parser) atom() Expr {
	tkn := p.expectRange("atom type check", tokenName, tokenNumber,
		tokenRawString, tokenQuotedString, tokenNull, tokenFalse, tokenTrue,
		tokenLeftRound, tokenLeftSquare,
	)
	switch tkn.typ {
	case tokenName:
		return newID(tkn.value, tkn)
	case tokenNumber:
		n, err := newNumber(tkn.value, tkn)
		if err != nil {
			p.error(err)
		}
		return n
	case tokenRawString, tokenQuotedString:
		return newString(tkn.value, tkn)
	case tokenNull:
		return newNull(tkn.value, tkn)
	case tokenFalse, tokenTrue:
		n, err := newBool(tkn.value, tkn.typ, tkn)
		if err != nil {
			p.error(err)
		}
		return n
	case tokenLeftRound:
		n := p.orEval()
		p.expect("closing brackets, expected ')'", tokenRightRound)
		return n
	case tokenLeftSquare:
		elements := p.exprList()
		p.expect("closing square brackets, expected ']'", tokenRightSquare)
		return newList(elements, tkn)
	// case tokenLeftCurly:

	default:
		p.unexpected("atom", tkn)
		return nil
	}
}

// exprList: orEval ("," orEval)* [","];
func (p *Parser) exprList() []Expr {
	elements := []Expr{p.orEval()}
	for p.peek().typ == tokenComma {
		p.next() // consume the comma token
		// if the following token isn't ']' handles dangling commas as well
		if p.peek().typ != tokenRightSquare {
			elements = append(elements, p.orEval())
		}
	}
	return elements
}
