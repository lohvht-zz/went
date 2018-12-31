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

// nextComma checks if the parser's next token is a comma
func (p *Parser) nextComma(context string, followTyp token.Type) bool {
	typ := p.peek().Type
	if typ == token.COMMA {
		return true
	}
	if typ != followTyp {
		msg := "missing comma ','"
		if typ == token.SEMICOLON && p.peek().Value == "\n" {
			msg += " before newline"
		}
		p.errorf(msg + " in " + context)
		return true
	}
	return false
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
	// p.Root = p.orEval()
	// if p.peek().Type == token.SEMICOLON {
	// 	p.next() // just consume the semicolon for now
	// }
	// p.expect("End of File", token.EOF)
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

// Expr : NotExpr (("||" | "&&") NotExpr)*;
func (p *Parser) parseExpr() Expr {
	node := p.parseNotExpr()
	for p.peek().Type == token.LOGICALOR || p.peek().Type == token.LOGICALAND {
		tkn := p.next()
		node = newBinExpr(node, p.parseNotExpr(), tkn)
	}
	return node
}

// NotExpr : "!" NotExpr | ComparisonExpr;
func (p *Parser) parseNotExpr() Expr {
	switch p.peek().Type {
	case token.LOGICALNOT:
		tkn := p.next()
		return newUnExpr(p.parseNotExpr(), tkn)
	default:
		return p.parseComparisonExpr()
	}
}

// ComparisonExpr : AddExpr (compOp AddExpr)*;
// comparison_op: "==" | "!=" | "<" | ">" | "<=" | ">=" | ["!"] "in";
func (p *Parser) parseComparisonExpr() Expr {
	node := p.parseAddExpr()
Loop:
	for {
		switch p.peek().Type {
		case token.EQ, token.NEQ, token.SM, token.SMEQ,
			token.GR, token.GREQ, token.IN:
			tkn := p.next()
			node = newBinExpr(node, p.parseAddExpr(), tkn)
		default:
			break Loop
		}
	}
	return node
}

// AddExpr : MultExpr (("+" | "-") MultExpr)*;
func (p *Parser) parseAddExpr() Expr {
	node := p.parseMultExpr()
Loop:
	for {
		switch p.peek().Type {
		case token.PLUS, token.MINUS:
			tkn := p.next()
			node = newBinExpr(node, p.parseMultExpr(), tkn)
		default:
			break Loop
		}
	}
	return node
}

// MultExpr : UnExpr (("*" | "/" | "%") UnExpr)*;
func (p *Parser) parseMultExpr() Expr {
	node := p.parseUnExpr()
Loop:
	for {
		switch p.peek().Type {
		case token.MULT, token.DIV, token.MOD:
			tkn := p.next()
			node = newBinExpr(node, p.parseUnExpr(), tkn)
		default:
			break Loop
		}
	}
	return node
}

// UnExpr : ("+" | "-") UnExpr | PrimaryExpr;
func (p *Parser) parseUnExpr() Expr {
	switch p.peek().Type {
	case token.PLUS, token.MINUS:
		tkn := p.next()
		return newUnExpr(p.parseUnExpr(), tkn)
	default:
		return p.parsePrimaryExpr()
	}
}

// PrimaryExpr : Operand (Selector | Index | Slice | Args)*;
// Selector: "." Name;
// Index: "[" Expr "]";
// Slice: "[ [Expr] ":" [Expr] [":" [Expr]] "]";
// Args: "(" [Expr ("," Expr)* [","]] ")";
func (p *Parser) parsePrimaryExpr() Expr {
	return nil
}

// Operand : Literal | Name | "(" Expr ")";
// Literal: BasicLit | CompositeLit; // NOTE: FuncLit support in the future?
// BasicLit: int | float | str; // NOTE: imaginary support in the future?
// CompositeLit: Array | Dict;
func (p *Parser) parseOperand() Expr {
	switch p.peek().Type {
	case token.INT, token.FLOAT, token.STR: // BasicLit
		return newBasicLit(p.next())
	case token.LROUND:
		p.next() // consume
		n := p.parseExpr()
		p.expect("Operand", token.RROUND)
		return n
	case token.LSQUARE: // Array
		return p.parseArray()
	case token.LCURLY: // Dict
		return p.parseDict()

	}
	return nil
}

// Array : "[" [Expr ("," Expr)* [","]] "]";
func (p *Parser) parseArray() Expr {
	lsquare := p.expect("Array", token.LSQUARE)
	var elements []Expr
	if p.peek().Type != token.RSQUARE {
		for p.peek().Type != token.RSQUARE && p.peek().Type != token.EOF {
			elements = append(elements, p.parseExpr())
			if !p.nextComma("Array", token.RSQUARE) {
				break
			}
			p.next() // consume the comma
		}
	}
	rsquare := p.expect("Array", token.RSQUARE)
	return newList(elements, lsquare, rsquare)
}

// Dict: "{" [DictEl ("," DictEl)* [","]] "}";
// DictEl: Key ":" Expr | Name;
// Key: Name | str;
func (p *Parser) parseDict() Expr {
	lcurly := p.expect("Dictionary", token.LCURLY)
	var elements []Expr
	if p.peek().Type != token.RCURLY {
		for p.peek().Type != token.RSQUARE && p.peek().Type != token.EOF {
			elements = append(elements, p.parseDictEl())
			if !p.nextComma("Dictionary", token.RCURLY) {
				break
			}
			p.next()
		}
	}
	rcurly := p.expect("Dictionary", token.RCURLY)
	return newDict(elements, lcurly, rcurly)
}

func (p *Parser) parseDictEl() Expr {
	return nil
}
