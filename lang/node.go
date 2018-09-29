package lang

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lohvht/went/lang/token"
)

var textFormat = "%s" // change to "%q" in tests for better error messages

// Node is an element from the parse tree
type Node interface {
	Scope() Scope            // returns the scope of a Node
	Accept(NodeWalker) WType // Accepts and marshalls the Nodewalker to the correct visit function
	Start() token.Position   // position of the 1st character belonging to the node
	End() token.Position     // position of the 1st character immediately after the node
}

// Stmt interface, all statment nodes implements this
type Stmt interface {
	Node
	stmt()
}

// Expr interface, all expression nodes implements this
type Expr interface {
	Node
	expr()
}

// ExprStmt is an expression statement, it can have a comma separated
// series of expressions
type ExprStmt struct {
	token.Token
	scope Scope
	exprs []Expr
}

func newExprStmt(expressions []Expr, tkn token.Token) *ExprStmt {
	return &ExprStmt{exprs: expressions, Token: tkn}
}

// Scope returns the scope that the statement was in
func (n *ExprStmt) Scope() Scope { return n.scope }

// Accept marshalls the AST node walker to the correct visit method
func (n *ExprStmt) Accept(nw NodeWalker) WType { return nw.visitExprStmt(n) }
func (n *ExprStmt) stmt()                      {}

// AssignStmt is the assignment statement
type AssignStmt struct {
	token.Token
	scope Scope
	left  []Expr
	right []Expr
}

func newAssignStmt(left, right []Expr, tkn token.Token) *AssignStmt {
	return &AssignStmt{left: left, right: right, Token: tkn}
}

// Scope returns the scope that the statement was in
func (n *AssignStmt) Scope() Scope { return n.scope }

// Accept marshalls the AST node walker to the correct visit method
func (n *AssignStmt) Accept(nw NodeWalker) WType { return nw.visitAssignStmt(n) }
func (n *AssignStmt) stmt()                      {}

// PlusAssignStmt is the assignment statement
type PlusAssignStmt struct {
	token.Token
	scope Scope
	left  []Expr
	right []Expr
}

func newPlusAssignStmt(left, right []Expr, tkn token.Token) *PlusAssignStmt {
	return &PlusAssignStmt{left: left, right: right, Token: tkn}
}

// Scope returns the scope that the statement was in
func (n *PlusAssignStmt) Scope() Scope { return n.scope }

// Accept marshalls the AST node walker to the correct visit method
func (n *PlusAssignStmt) Accept(nw NodeWalker) WType { return nw.visitPlusAssignStmt(n) }
func (n *PlusAssignStmt) stmt()                      {}

// MinusAssignStmt is the assignment statement
type MinusAssignStmt struct {
	token.Token
	scope Scope
	left  []Expr
	right []Expr
}

func newMinusAssignStmt(left, right []Expr, tkn token.Token) *MinusAssignStmt {
	return &MinusAssignStmt{left: left, right: right, Token: tkn}
}

// Scope returns the scope that the statement was in
func (n *MinusAssignStmt) Scope() Scope { return n.scope }

// Accept marshalls the AST node walker to the correct visit method
func (n *MinusAssignStmt) Accept(nw NodeWalker) WType { return nw.visitMinusAssignStmt(n) }
func (n *MinusAssignStmt) stmt()                      {}

// DivAssignStmt is the assignment statement
type DivAssignStmt struct {
	token.Token
	scope Scope
	left  []Expr
	right []Expr
}

func newDivAssignStmt(left, right []Expr, tkn token.Token) *DivAssignStmt {
	return &DivAssignStmt{left: left, right: right, Token: tkn}
}

// Scope returns the scope that the statement was in
func (n *DivAssignStmt) Scope() Scope { return n.scope }

// Accept marshalls the AST node walker to the correct visit method
func (n *DivAssignStmt) Accept(nw NodeWalker) WType { return nw.visitDivAssignStmt(n) }
func (n *DivAssignStmt) stmt()                      {}

// MultAssignStmt is the assignment statement
type MultAssignStmt struct {
	token.Token
	scope Scope
	left  []Expr
	right []Expr
}

func newMultAssignStmt(left, right []Expr, tkn token.Token) *MultAssignStmt {
	return &MultAssignStmt{left: left, right: right, Token: tkn}
}

// Scope returns the scope that the statement was in
func (n *MultAssignStmt) Scope() Scope { return n.scope }

// Accept marshalls the AST node walker to the correct visit method
func (n *MultAssignStmt) Accept(nw NodeWalker) WType { return nw.visitMultAssignStmt(n) }
func (n *MultAssignStmt) stmt()                      {}

// ModAssignStmt is the assignment statement
type ModAssignStmt struct {
	token.Token
	scope Scope
	left  []Expr
	right []Expr
}

func newModAssignStmt(left, right []Expr, tkn token.Token) *ModAssignStmt {
	return &ModAssignStmt{left: left, right: right, Token: tkn}
}

// Scope returns the scope that the statement was in
func (n *ModAssignStmt) Scope() Scope { return n.scope }

// Accept marshalls the AST node walker to the correct visit method
func (n *ModAssignStmt) Accept(nw NodeWalker) WType { return nw.visitModAssignStmt(n) }
func (n *ModAssignStmt) stmt()                      {}

// An expression is represented by a tree consisting of one or more of
// the following concrete expression nodes.

type (
	// BinExpr holds a binary operator between left and right expressions
	BinExpr struct {
		operation token.Token
		left      Expr
		right     Expr
		scope     Scope
	}

	// UnExpr holds a unary operator over its operand expression
	UnExpr struct {
		operation token.Token
		scope     Scope
		operand   Expr
	}
)

func newBinExpr(left, right Expr, op token.Token) *BinExpr {
	return &BinExpr{operation: op, left: left, right: right}
}

// Scope returns the scope the current node is in
func (n *BinExpr) Scope() Scope               { return n.scope }
func (n *BinExpr) Accept(nw NodeWalker) WType { return nw.visitBinExpr(n) }
func (n *BinExpr) expr()                      {}

func newUnExpr(operand Expr, op token.Token) *UnExpr {
	return &UnExpr{operation: op, operand: operand}
}

func (n *UnExpr) Scope() Scope               { return n.scope }
func (n *UnExpr) Accept(nw NodeWalker) WType { return nw.visitUnExpr(n) }
func (n *UnExpr) expr()                      {}

// // Atom expressions
// type funcCall struct {
// }

/* Literals */

type literal struct {
	token.Token
	scope Scope
}

func (n literal) Scope() Scope { return n.scope }
func (n literal) expr()        {}

// func (n literal) Accept(nw NodeWalker) WType { return nw.visitLiteral(n) }

// Num holds a numerical constant: signed integer or float
type Num struct {
	literal
	IsInt   bool    // number has an int value
	IsFloat bool    // number has floating point value
	Int64   int64   // signed integer value
	Float64 float64 // floating point value
	Text    string  // Original text representation from input
}

// newNumber creates a new pointer to the Num
func newNumber(text string, tkn token.Token) (*Num, error) {
	n := &Num{literal: literal{Token: tkn}, Text: text}
	i, err := strconv.ParseInt(text, 0, 64)
	// If an int extraction succeeded, promote the float
	if err == nil {
		n.IsInt = true
		n.Int64 = i
	}
	if n.IsInt {
		// If an integer extraction is successful, promote the float
		n.IsFloat = true
		n.Float64 = float64(n.Int64)
	} else {
		// Else an integer extraction was initially unsuccessful, process the float
		f, err := strconv.ParseFloat(text, 64)
		if err == nil {
			// If we parsed it as a float, but looks like an integer,
			// it's a huge number too large to fit in an integer. Reject it
			if !strings.ContainsAny(text, ".eE") {
				return nil, fmt.Errorf("Integer overflow: %q", text)
			}
			n.IsFloat = true
			n.Float64 = f
			// If a floating-point extraction succeeded, extract the int if needed.
			if !n.IsInt && float64(int64(f)) == f {
				n.IsInt = true
				n.Int64 = int64(f)
			}
		}
	}
	if !n.IsInt && !n.IsFloat {
		return nil, fmt.Errorf("illegal number syntax: %q", text)
	}
	return n, nil
}

// Accept marshalls the AST node walker to the correct visit method
func (n *Num) Accept(nw NodeWalker) WType { return nw.visitNum(n) }

// Str holds a string literal: both raw and quoted
type Str struct {
	literal
	Value string
}

// newString creates a new pointer to the Str
func newString(text string, tkn token.Token) *Str {
	return &Str{literal: literal{Token: tkn}, Value: text}
}

// Accept marshalls the AST node walker to the correct visit method
func (n *Str) Accept(nw NodeWalker) WType { return nw.visitStr(n) }

// Null holds a null literal
type Null struct {
	literal
	Text string
}

// newNull creates a new pointer to the Null
func newNull(text string, tkn token.Token) *Null {
	return &Null{literal: literal{Token: tkn}, Text: text}
}

// Accept marshalls the AST node walker to the correct visit method
func (n *Null) Accept(nw NodeWalker) WType { return nw.visitNull(n) }

// Bool holds a boolean literal
type Bool struct {
	literal
	Value bool
	Text  string
}

// newBool creates a new pointer to the Bool
func newBool(text string, tknTyp token.Type, tkn token.Token) (*Bool, error) {
	switch tknTyp {
	case token.TRUE:
		return &Bool{literal: literal{Token: tkn}, Value: true, Text: text}, nil
	case token.FALSE:
		return &Bool{literal: literal{Token: tkn}, Value: false, Text: text}, nil
	default:
		return nil, fmt.Errorf("illegal bool syntax: %q", text)
	}
}

// Accept marshalls the AST node walker to the correct visit method
func (n *Bool) Accept(nw NodeWalker) WType { return nw.visitBool(n) }

// List holds a list of Nodes
type List struct {
	literal
	elements []Expr
}

func newList(elems []Expr, tkn token.Token) *List {
	return &List{literal: literal{Token: tkn}, elements: elems}
}

// Accept marshalls the AST node walker to the correct visit method
func (n *List) Accept(nw NodeWalker) WType { return nw.visitList(n) }

// ID node represents Identifier/Name nodes
type ID struct {
	literal
	value string
}

func newID(value string, tkn token.Token) *ID { return &ID{literal{Token: tkn}, value} }

// Scope : self-explanatory
func (n *ID) Scope() Scope { return n.scope }

// Accept marshalls the AST node walker to the correct visit method
func (n *ID) Accept(nw NodeWalker) WType { return nw.visitID(n) }

// TODO: fix all the errors
// 1) Remove individual types for basic literals (add in basic literal kinds)
// 		1.1) this is for constant literals like FLOAT, INT, STRING, RAWSTRING,
//		BOOL & NULL
// 		1.2) Implement them as a list of literals
// 2) implement the Pos Start and End methods for what we have now
// 3) Implement Composite literal types (Lists, Maps, Classes)
// 4) Investigate how golang traverses its ast
// 		4.1) Refactor traversing and decouple it to any 1 implementation
// 5) Refactor how errors are being handled in the first place!
// (parse/syntax errors should not terminate the parsing --> Collect errors)
