package lang

import (
	"fmt"
	"runtime"

	"github.com/lohvht/went/lang/token"
)

// Parser parses the input string (file or otherwise) and creates an AST at its Root
// also links the AST to the appropriate scopes
type Parser struct {
	Name string
	Root Node // top-level root of the AST tree
	// symtab *SymbolTable // the entire symbol table, global scope, local scope, functions etc.
	// currentScope *Scope
	input        string // input text to be parsed
	tokeniser    *token.Lexer
	tokens       token.List  // list of token lookaheads
	currentToken token.Token // the local that we are currently looking at (Not a lookahead)
}

// next consumes and returns the next token
func (p *Parser) next() token.Token {
	// take a token from the bottom of the stack
	if !p.tokens.Empty() {
		p.currentToken = p.tokens.Shift()
	} else {
		p.currentToken = p.tokeniser.Next()
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
	p.tokens.Push(p.tokeniser.Next())
	return p.tokens.PeekBottom()
}

// Parsing

// errorf formats the error and terminates processing.
func (p *Parser) errorf(format string, args ...interface{}) {
	p.Root = nil
	format = fmt.Sprintf("%s: SyntaxError - %s", p.currentToken.Pos.String(), format)
	panic(fmt.Errorf(format, args...))
}

// error terminates the processing.
func (p *Parser) error(err error) { p.errorf("%s", err) }

// expect consumes the next token and guarantees it has the required type.
func (p *Parser) expect(context string, expected token.Type) token.Token {
	tkn := p.next()
	if tkn.Type != expected {
		p.unexpected(context, tkn)
	}
	return tkn
}

// expectRange consumes the next token and guarantees it has one of the required types.
func (p *Parser) expectRange(context string, expectedTypes ...token.Type) (tkn token.Token) {
	tkn = p.next()
	for _, exTyp := range expectedTypes {
		if tkn.Type == exTyp {
			return
		}
	}
	p.unexpected(context, tkn)
	return
}

// unexpected complains about the token and terminates processing
func (p *Parser) unexpected(context string, tkn token.Token) {
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
			p.tokeniser.Drain()
			p.stopParse()
		}
		*errp = e.(error)
	}
}

// initParser initialises the parser, using the token.Lexer
func initParser(tokeniser *token.Lexer) *Parser {
	p := &Parser{Name: tokeniser.Name, Root: nil, tokeniser: tokeniser,
		input: tokeniser.Input}
	return p
}

func (p *Parser) stopParse() { p.tokeniser = nil }

// Parse parses the input string to construct an AST
func Parse(name, input string) (parser *Parser, err error) {
	p := initParser(token.Tokenise(name, input))
	defer p.recover(&err)
	p.parse()
	p.stopParse()
	return p, nil
}

func (p *Parser) parse() {
	p.Root = p.orEval()
	if p.peek().Type == token.SEMICOLON {
		p.next() // just consume the semicolon for now
	}
	p.expect("End of File", token.EOF)
}

// Grammar rules

// func (p *Parser) input() Node {
// 	for p.peek().Type != tokenEOF {

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
// 	switch tkntyp := p.peek().Type; tkntyp {
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

// func (p *Parser) assignStmt(lhs []Expr, typ token.Type) Stmt {
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
	for p.peek().Type == token.LOGICALOR {
		tkn := p.next()
		node = newBinExpr(node, p.andEval(), tkn)
	}
	return node
}

// andEval: notEval ("&&" notEval)*;
func (p *Parser) andEval() Expr {
	node := p.notEval()
	for p.peek().Type == token.LOGICALAND {
		tkn := p.next()
		node = newBinExpr(node, p.andEval(), tkn)
	}
	return node
}

// notEval: "!" notEval | comparison;
func (p *Parser) notEval() Expr {
	switch p.peek().Type {
	case token.LOGICALNOT:
		tkn := p.next()
		return newUnExpr(p.notEval(), tkn)
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
		switch p.peek().Type {
		case token.EQ, token.NEQ,
			token.SM, token.SMEQ,
			token.GR, token.GREQ, token.IN:
			tkn := p.next()
			node = newBinExpr(node, p.smExpr(), tkn)
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
		switch p.peek().Type {
		case token.PLUS, token.MINUS:
			tkn := p.next()
			node = newBinExpr(node, p.term(), tkn)
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
		switch p.peek().Type {
		case token.MULT, token.DIV, token.MOD:
			tkn := p.next()
			node = newBinExpr(node, p.factor(), tkn)
		default:
			break Loop
		}
	}
	return node
}

// factor: ("+" | "-") factor | atom;
func (p *Parser) factor() Expr {
	switch p.peek().Type {
	case token.PLUS, token.MINUS:
		tkn := p.next()
		return newUnExpr(p.factor(), tkn)
	default:
		return p.atom()
	}
}

// TODO: Implement me!
// atomExpr: atom trailer*;
// trailer: "(" [argList] ")" | "[" slice "]" | "." NAME;
// slice: orEval | [orEval] ":" [orEval] [":" [orEval]];
// argList: arg ("," arg)* [","];
// arg: orEval | NAME "=" orEval;
func (p *Parser) atomExpr() Expr {
	n := p.atom()
TrailerLoop:
	for {
		switch p.peek().Type {
		case token.DOT:
		case token.LROUND:
		case token.LSQUARE:

		default:
			break TrailerLoop
		}
	}
	return n
}

// atom: "[" [exprList] "]" | "{" mapList "}" | "(" expr ")" | ID | NUM | STR |
// RAWSTR | "null" | "false" | "true";
// mapList: keyval ("," keyval)* [","];
// keyval: (ID | STR) ":" expr;
func (p *Parser) atom() Expr {
	tkn := p.expectRange("atom type check", token.NAME, token.INT, token.FLOAT,
		token.RAWSTR, token.STR, token.NULL, token.FALSE, token.TRUE,
		token.LROUND, token.LSQUARE)
	switch tkn.Type {
	case token.NAME:
		return newID(tkn.Value, tkn)
	case token.NUM:
		n, err := newNumber(tkn.Value, tkn)
		if err != nil {
			p.error(err)
		}
		return n
	case token.RAWSTR, token.STR:
		return newString(tkn.Value, tkn)
	case token.NULL:
		return newNull(tkn.Value, tkn)
	case token.FALSE, token.TRUE:
		n, err := newBool(tkn.Value, tkn.Type, tkn)
		if err != nil {
			p.error(err)
		}
		return n
	case token.LROUND:
		n := p.orEval()
		p.expect("closing brackets, expected ')'", token.RROUND)
		return n
	case token.LSQUARE:
		elements := p.exprList()
		p.expect("closing square brackets, expected ']'", token.RSQUARE)
		return newList(elements, tkn)
	// case token.LeftCurly:

	default:
		p.unexpected("atom", tkn)
		return nil
	}
}

// exprList: orEval ("," orEval)* [","];
func (p *Parser) exprList() []Expr {
	elements := []Expr{p.orEval()}
	for p.peek().Type == token.COMMA {
		p.next() // consume the comma token
		// if the following token isn't ']' handles dangling commas as well
		if p.peek().Type != token.RSQUARE {
			elements = append(elements, p.orEval())
		}
	}
	return elements
}
