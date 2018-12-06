package lang

import (
	"fmt"
	"strconv"

	"github.com/lohvht/went/lang/token"
)

var textFormat = "%s" // change to "%q" in tests for better error messages

// Interfaces
type (
	// Node is the AST node interface
	Node interface {
		Tkn() token.Token
		NScope() Scope           // returns the scope of a Node
		accept(NodeWalker) WType // Accepts and marshalls the Nodewalker to the correct visit function
	}

	// Stmt interface, all statment nodes implements this
	Stmt interface {
		Node
		stmt()
	}

	// Expr interface, all expression nodes implements this
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
		token.Token
		Scope
		exprs []Expr
	}
	// AssignStmt is the assignment statement
	AssignStmt struct {
		token.Token
		Scope
		left  []Expr
		right []Expr
	}
	// PlusAssignStmt is the assignment statement
	PlusAssignStmt struct {
		token.Token
		Scope
		left  []Expr
		right []Expr
	}
	// MinusAssignStmt is the assignment statement
	MinusAssignStmt struct {
		token.Token
		Scope
		left  []Expr
		right []Expr
	}
	// DivAssignStmt is the assignment statement
	DivAssignStmt struct {
		token.Token
		Scope
		left  []Expr
		right []Expr
	}
	// MultAssignStmt is the assignment statement
	MultAssignStmt struct {
		token.Token
		Scope
		left  []Expr
		right []Expr
	}
	// ModAssignStmt is the assignment statement
	ModAssignStmt struct {
		token.Token
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

func newExprStmt(expressions []Expr, tkn token.Token) *ExprStmt {
	return &ExprStmt{exprs: expressions, Token: tkn}
}
func newAssignStmt(left, right []Expr, tkn token.Token) *AssignStmt {
	return &AssignStmt{left: left, right: right, Token: tkn}
}
func newPlusAssignStmt(left, right []Expr, tkn token.Token) *PlusAssignStmt {
	return &PlusAssignStmt{left: left, right: right, Token: tkn}
}
func newMinusAssignStmt(left, right []Expr, tkn token.Token) *MinusAssignStmt {
	return &MinusAssignStmt{left: left, right: right, Token: tkn}
}
func newDivAssignStmt(left, right []Expr, tkn token.Token) *DivAssignStmt {
	return &DivAssignStmt{left: left, right: right, Token: tkn}
}
func newMultAssignStmt(left, right []Expr, tkn token.Token) *MultAssignStmt {
	return &MultAssignStmt{left: left, right: right, Token: tkn}
}
func newModAssignStmt(left, right []Expr, tkn token.Token) *ModAssignStmt {
	return &ModAssignStmt{left: left, right: right, Token: tkn}
}

// Expressions
// An expression is represented by a tree consisting of one or more of
// the following concrete expression nodes.
type (
	// BinExpr holds a binary operator between left and right expressions
	BinExpr struct {
		token.Token
		Scope
		left  Expr
		right Expr
	}
	// UnExpr holds a unary operator over its operand expression
	UnExpr struct {
		token.Token
		Scope
		operand Expr
	}
)

func (n *BinExpr) accept(nw NodeWalker) WType { return nw.visitBinExpr(n) }
func (n *UnExpr) accept(nw NodeWalker) WType  { return nw.visitUnExpr(n) }

func (n *BinExpr) expr() {}
func (n *UnExpr) expr()  {}

func newBinExpr(left, right Expr, op token.Token) *BinExpr {
	return &BinExpr{Token: op, left: left, right: right}
}
func newUnExpr(operand Expr, op token.Token) *UnExpr {
	return &UnExpr{Token: op, operand: operand}
}

// // Atom expressions
// type funcCall struct {
// }

// Literals
type (
	// Num holds a numerical constant: signed integer or float
	Num struct {
		token.Token
		Scope
		IsInt   bool    // number has an int value
		IsFloat bool    // number has floating point value
		Int64   int64   // signed integer value
		Float64 float64 // floating point value
		Text    string  // Original text representation from input
	}
	// Str holds a string both raw and quoted
	Str struct {
		token.Token
		Scope
		Value string
	}
	// Null holds a null literal
	Null struct {
		token.Token
		Scope
		Text string
	}
	// Bool holds a boolean literal
	Bool struct {
		token.Token
		Scope
		Value bool
		Text  string
	}
	// List holds a list of Nodes
	List struct {
		token.Token
		Scope
		elements []Expr
	}
	// ID node represents Identifier/Name nodes
	ID struct {
		token.Token
		Scope
		value string
	}
)

func (n *Num) accept(nw NodeWalker) WType  { return nw.visitNum(n) }
func (n *Str) accept(nw NodeWalker) WType  { return nw.visitStr(n) }
func (n *Null) accept(nw NodeWalker) WType { return nw.visitNull(n) }
func (n *Bool) accept(nw NodeWalker) WType { return nw.visitBool(n) }
func (n *List) accept(nw NodeWalker) WType { return nw.visitList(n) }
func (n *ID) accept(nw NodeWalker) WType   { return nw.visitID(n) }

func (n *Num) expr()  {}
func (n *Str) expr()  {}
func (n *Null) expr() {}
func (n *Bool) expr() {}
func (n *List) expr() {}
func (n *ID) expr()   {}

// newNumber creates a new pointer to the Num
func newNumber(tkn token.Token) (*Num, error) {
	text := tkn.Value
	n := &Num{Token: tkn, Text: text}
	switch tkn.Type {
	case token.FLOAT:
		f, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return nil, fmt.Errorf("parse float error: %q", text)
		}
		n.IsFloat = true
		n.Float64 = f
		if float64(int64(f)) == f {
			n.IsInt = true
			n.Int64 = int64(f)
		}
	case token.INT:
		i, err := strconv.ParseInt(text, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("parse int error: %q", text)
		}
		n.IsInt = true
		n.Int64 = i
		n.IsFloat = true
		n.Float64 = float64(n.Int64)
	}
	if !n.IsInt && !n.IsFloat {
		return nil, fmt.Errorf("illegal number syntax: %q", text)
	}
	return n, nil
}

// newString creates a new pointer to the Str
func newString(tkn token.Token) *Str { return &Str{Token: tkn, Value: tkn.Value} }

// newNull creates a new pointer to the Null
func newNull(tkn token.Token) *Null { return &Null{Token: tkn, Text: tkn.Value} }

// newBool creates a new pointer to the Bool
func newBool(tkn token.Token) (*Bool, error) {
	switch tkn.Type {
	case token.TRUE:
		return &Bool{Token: tkn, Value: true, Text: tkn.Value}, nil
	case token.FALSE:
		return &Bool{Token: tkn, Value: false, Text: tkn.Value}, nil
	default:
		return nil, fmt.Errorf("illegal bool syntax: %q", tkn.Value)
	}
}

func newList(elems []Expr, tkn token.Token) *List {
	return &List{Token: tkn, elements: elems}
}

func newID(tkn token.Token) *ID { return &ID{Token: tkn, value: tkn.Value} }
