package lang

import (
	"github.com/lohvht/went/lang/token"
)

var textFormat = "%s" // change to "%q" in tests for better error messages

// Interfaces
type (
	// Node is the AST node interface
	Node interface {
		NScope() Scope           // returns the scope of a Node
		Pos() token.Pos          // starting position in the code which represents this node and its children
		End() token.Pos          // end position in the code which represents this node and its children
		accept(NodeWalker) WType // Accepts and marshalls the Nodewalker to the correct visit function
	}

	// Stmt interface.
	Stmt interface {
		Node
		stmt()
	}

	// Expr interface, expressions evaluate to a value
	Expr interface {
		Node
		expr()
	}
)

// Statements
type (
	// ExprStmt is an expression statement, it can have a comma separated
	// series of expressions
	ExprStmt struct {
		Scope
		exprs []Expr
	}
	// AssignStmt is the assignment statement
	AssignStmt struct {
		Scope
		left  []Expr
		right []Expr
	}
	// PlusAssignStmt is the assignment statement
	PlusAssignStmt struct {
		Scope
		left  []Expr
		right []Expr
	}
	// MinusAssignStmt is the assignment statement
	MinusAssignStmt struct {
		Scope
		left  []Expr
		right []Expr
	}
	// DivAssignStmt is the assignment statement
	DivAssignStmt struct {
		Scope
		left  []Expr
		right []Expr
	}
	// MultAssignStmt is the assignment statement
	MultAssignStmt struct {
		Scope
		left  []Expr
		right []Expr
	}
	// ModAssignStmt is the assignment statement
	ModAssignStmt struct {
		Scope
		left  []Expr
		right []Expr
	}
)

func (n *ExprStmt) accept(nw NodeWalker) WType        { return nw.visitExprStmt(n) }
func (n *AssignStmt) accept(nw NodeWalker) WType      { return nw.visitAssignStmt(n) }
func (n *PlusAssignStmt) accept(nw NodeWalker) WType  { return nw.visitPlusAssignStmt(n) }
func (n *MinusAssignStmt) accept(nw NodeWalker) WType { return nw.visitMinusAssignStmt(n) }
func (n *DivAssignStmt) accept(nw NodeWalker) WType   { return nw.visitDivAssignStmt(n) }
func (n *MultAssignStmt) accept(nw NodeWalker) WType  { return nw.visitMultAssignStmt(n) }
func (n *ModAssignStmt) accept(nw NodeWalker) WType   { return nw.visitModAssignStmt(n) }

func (n *ExprStmt) stmt()        {}
func (n *AssignStmt) stmt()      {}
func (n *PlusAssignStmt) stmt()  {}
func (n *MinusAssignStmt) stmt() {}
func (n *DivAssignStmt) stmt()   {}
func (n *MultAssignStmt) stmt()  {}
func (n *ModAssignStmt) stmt()   {}

// func newExprStmt(expressions []Expr, tkn token.Token) *ExprStmt {
// 	return &ExprStmt{exprs: expressions, Token: tkn}
// }
// func newAssignStmt(left, right []Expr, tkn token.Token) *AssignStmt {
// 	return &AssignStmt{left: left, right: right, Token: tkn}
// }
// func newPlusAssignStmt(left, right []Expr, tkn token.Token) *PlusAssignStmt {
// 	return &PlusAssignStmt{left: left, right: right, Token: tkn}
// }
// func newMinusAssignStmt(left, right []Expr, tkn token.Token) *MinusAssignStmt {
// 	return &MinusAssignStmt{left: left, right: right, Token: tkn}
// }
// func newDivAssignStmt(left, right []Expr, tkn token.Token) *DivAssignStmt {
// 	return &DivAssignStmt{left: left, right: right, Token: tkn}
// }
// func newMultAssignStmt(left, right []Expr, tkn token.Token) *MultAssignStmt {
// 	return &MultAssignStmt{left: left, right: right, Token: tkn}
// }
// func newModAssignStmt(left, right []Expr, tkn token.Token) *ModAssignStmt {
// 	return &ModAssignStmt{left: left, right: right, Token: tkn}
// }

// Expressions
// An expression is represented by a tree consisting of one or more of
// the following concrete expression nodes.
type (
	// BinExpr holds a binary operator between left and right expressions
	BinExpr struct {
		op    token.Token
		opPos token.Pos
		Scope
		left  Expr
		right Expr
	}
	// UnExpr holds a unary operator over its operand expression
	UnExpr struct {
		op    token.Token
		opPos token.Pos
		Scope
		operand Expr
	}
)

func (n *BinExpr) accept(nw NodeWalker) WType { return nw.visitBinExpr(n) }
func (n *UnExpr) accept(nw NodeWalker) WType  { return nw.visitUnExpr(n) }

func (n *BinExpr) expr() {}
func (n *UnExpr) expr()  {}

func (n *BinExpr) Pos() token.Pos { return n.left.Pos() }
func (n *UnExpr) Pos() token.Pos  { return n.opPos }

func (n *BinExpr) End() token.Pos { return n.right.End() }
func (n *UnExpr) End() token.Pos  { return n.operand.End() }

func newBinExpr(left, right Expr, op token.Token) *BinExpr {
	return &BinExpr{op: op, opPos: op.Pos, left: left, right: right}
}
func newUnExpr(operand Expr, op token.Token) *UnExpr {
	return &UnExpr{op: op, opPos: op.Pos, operand: operand}
}

// // Atom expressions
// type funcCall struct {
// }

// Literals
type (
	// BasicLit node represents a literal of basic type
	BasicLit struct {
		token.Token // token.INT, token.FLOAT, token.STR, token.BOOL, token.NULL
		Scope
		Text string
	}

	// List holds a list of literal nodes
	List struct {
		LSqPos token.Pos // the position of the opening square bracket "["
		RSqPos token.Pos // the position of the closing square bracket "]"
		Scope
		elements []Expr
	}
	// Ident node represents Identifier/Name nodes
	Ident struct {
		token.Token
		Scope
		Name string
	}
)

func (n *BasicLit) accept(nw NodeWalker) WType { return nw.visitBasicLit(n) }
func (n *List) accept(nw NodeWalker) WType     { return nw.visitList(n) }
func (n *Ident) accept(nw NodeWalker) WType    { return nw.visitID(n) }

func (n *BasicLit) Pos() token.Pos { return n.Token.Pos }
func (n *List) Pos() token.Pos     { return n.LSqPos }
func (n *Ident) Pos() token.Pos    { return n.Token.Pos }

func (n *BasicLit) End() token.Pos { return token.AddOffset(n.Token.Pos, len(n.Text)) }
func (n *List) End() token.Pos     { return n.RSqPos }
func (n *Ident) End() token.Pos    { return token.AddOffset(n.Token.Pos, len(n.Name)) }

func (n *BasicLit) expr() {}
func (n *List) expr()     {}
func (n *Ident) expr()    {}

func newBasicLit(tkn token.Token) *BasicLit {
	return &BasicLit{Token: tkn, Text: tkn.Value}
}

func newList(elems []Expr, leftSquare, rightSquare token.Token) *List {
	return &List{elements: elems, LSqPos: leftSquare.Pos, RSqPos: rightSquare.Pos}
}

func newID(tkn token.Token) *Ident { return &Ident{Token: tkn, Name: tkn.Value} }
