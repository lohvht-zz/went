package lang

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/lohvht/went/lang/token"
)

var textFormat = "%s" // change to "%q" in tests for better error messages

// Node is an element from the parse tree
type Node interface {
	Scope() Scope // returns the scope of a Node
	String() string
	Accept(NodeWalker) WType // Accepts and marshalls the Nodewalker to the correct visit function
	Position() token.Pos     // byte position of start of the node, in full original input string
	LinePosition() token.LinePos
}

// Stmt is a statement node
type Stmt interface {
	Node
	Statement()
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

func (n *ExprStmt) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(n.exprs[0].String())
	for _, expr := range n.exprs[1:] {
		buffer.WriteString(",")
		buffer.WriteString(expr.String())
	}
	return buffer.String()
}

// Accept marshalls the AST node walker to the correct visit method
func (n *ExprStmt) Accept(nw NodeWalker) WType { return nw.visitExprStmt(n) }

// Statement node
func (n *ExprStmt) Statement() {}

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

func (n *AssignStmt) String() string { return "=" }

// Accept marshalls the AST node walker to the correct visit method
func (n *AssignStmt) Accept(nw NodeWalker) WType { return nw.visitAssignStmt(n) }

// Statement node
func (n *AssignStmt) Statement() {}

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

func (n *PlusAssignStmt) String() string { return "=" }

// Accept marshalls the AST node walker to the correct visit method
func (n *PlusAssignStmt) Accept(nw NodeWalker) WType { return nw.visitPlusAssignStmt(n) }

// Statement node
func (n *PlusAssignStmt) Statement() {}

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

func (n *MinusAssignStmt) String() string { return "=" }

// Accept marshalls the AST node walker to the correct visit method
func (n *MinusAssignStmt) Accept(nw NodeWalker) WType { return nw.visitMinusAssignStmt(n) }

// Statement node
func (n *MinusAssignStmt) Statement() {}

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

func (n *DivAssignStmt) String() string { return "=" }

// Accept marshalls the AST node walker to the correct visit method
func (n *DivAssignStmt) Accept(nw NodeWalker) WType { return nw.visitDivAssignStmt(n) }

// Statement node
func (n *DivAssignStmt) Statement() {}

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

func (n *MultAssignStmt) String() string { return "=" }

// Accept marshalls the AST node walker to the correct visit method
func (n *MultAssignStmt) Accept(nw NodeWalker) WType { return nw.visitMultAssignStmt(n) }

// Statement node
func (n *MultAssignStmt) Statement() {}

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

func (n *ModAssignStmt) String() string { return "=" }

// Accept marshalls the AST node walker to the correct visit method
func (n *ModAssignStmt) Accept(nw NodeWalker) WType { return nw.visitModAssignStmt(n) }

// Statement node
func (n *ModAssignStmt) Statement() {}

// Expr is an expression node
type Expr interface {
	Node
	Expression()
}

// binOpExpr holds a binary operator between a left and right expressions
// This struct is meant to be embedded within all other binary op structs
type binOpExpr struct {
	token.Token
	scope Scope
	left  Expr
	right Expr
}

func (n binOpExpr) Scope() Scope { return n.scope }

func (n binOpExpr) Expression() {}

// Arithmetic Binary Operators

// AddExpr holds a '+' operator between its 2 children
type AddExpr struct{ binOpExpr }

// newAdd returns a pointer to a AddExpr
func newAdd(left Expr, right Expr, tkn token.Token) *AddExpr {
	return &AddExpr{binOpExpr{left: left, right: right, Token: tkn}}
}

// Accept marshalls the AST node walker to the correct visit method
func (n *AddExpr) Accept(nw NodeWalker) WType { return nw.visitAdd(n) }

func (n *AddExpr) String() string { return "+" }

// SubtractExpr holds a '-' operator between its 2 children
type SubtractExpr struct{ binOpExpr }

// newSubtract returns a pointer to a SubtractExpr
func newSubtract(left Expr, right Expr, tkn token.Token) *SubtractExpr {
	return &SubtractExpr{binOpExpr{left: left, right: right, Token: tkn}}
}

// Accept marshalls the AST node walker to the correct visit method
func (n *SubtractExpr) Accept(nw NodeWalker) WType { return nw.visitSubtract(n) }

func (n *SubtractExpr) String() string { return "-" }

// MultExpr holds a '*' operator between its 2 children
type MultExpr struct{ binOpExpr }

// newMult returns a pointer to a MultExpr
func newMult(left Expr, right Expr, tkn token.Token) *MultExpr {
	return &MultExpr{binOpExpr{left: left, right: right, Token: tkn}}
}

// Accept marshalls the AST node walker to the correct visit method
func (n *MultExpr) Accept(nw NodeWalker) WType { return nw.visitMult(n) }

func (n *MultExpr) String() string { return "*" }

// DivExpr holds a '/' operator between its 2 children
type DivExpr struct{ binOpExpr }

// newDiv returns a pointer to a DivExpr
func newDiv(left Expr, right Expr, tkn token.Token) *DivExpr {
	return &DivExpr{binOpExpr{left: left, right: right, Token: tkn}}
}

// Accept marshalls the AST node walker to the correct visit method
func (n *DivExpr) Accept(nw NodeWalker) WType { return nw.visitDiv(n) }

func (n *DivExpr) String() string { return "/" }

// ModExpr holds a '%' operator between its 2 children
type ModExpr struct{ binOpExpr }

// newMod returns a pointer to a ModExpr
func newMod(left Expr, right Expr, tkn token.Token) *ModExpr {
	return &ModExpr{binOpExpr{left: left, right: right, Token: tkn}}
}

// Accept marshalls the AST node walker to the correct visit method
func (n *ModExpr) Accept(nw NodeWalker) WType { return nw.visitMod(n) }

func (n *ModExpr) String() string { return "%" }

// Comparative Binary Operators

// EqExpr holds either the '!=' or '==' operator between its 2 children
type EqExpr struct {
	binOpExpr
	IsNot bool
}

// newEq returns a pointer to a EqExpr
func newEq(left Expr, right Expr, isNot bool, tkn token.Token) *EqExpr {
	return &EqExpr{binOpExpr: binOpExpr{left: left, right: right, Token: tkn}, IsNot: isNot}
}

// Accept marshalls the AST node walker to the correct visit method
func (n *EqExpr) Accept(nw NodeWalker) WType { return nw.visitEq(n) }

func (n *EqExpr) String() string {
	if n.IsNot {
		return "!="
	}
	return "=="
}

// SmExpr holds either the '<' or '<=' operator between its 2 children
type SmExpr struct {
	binOpExpr
	OrEq bool
}

// newSm returns a pointer to a SmExpr
func newSm(left Expr, right Expr, OrEq bool, tkn token.Token) *SmExpr {
	return &SmExpr{binOpExpr: binOpExpr{left: left, right: right, Token: tkn}, OrEq: OrEq}
}

// Accept marshalls the AST node walker to the correct visit method
func (n *SmExpr) Accept(nw NodeWalker) WType { return nw.visitSm(n) }

func (n *SmExpr) String() string {
	if n.OrEq {
		return "<="
	}
	return "<"
}

// GrExpr holds either the '<' or '<=' operator between its 2 children
type GrExpr struct {
	binOpExpr
	OrEq bool
}

// Accept marshalls the AST node walker to the correct visit method
func (n *GrExpr) Accept(nw NodeWalker) WType { return nw.visitGr(n) }

// newGr returns a pointer to a GrExpr
func newGr(left Expr, right Expr, OrEq bool, tkn token.Token) *GrExpr {
	return &GrExpr{binOpExpr: binOpExpr{left: left, right: right, Token: tkn}, OrEq: OrEq}
}

func (n *GrExpr) String() string {
	if n.OrEq {
		return ">="
	}
	return ">"
}

// InExpr holds either the '!in' or 'in' operator between its 2 children
type InExpr struct{ binOpExpr }

// newIn returns a pointer to a InExpr
func newIn(left Expr, right Expr, tkn token.Token) *InExpr {
	return &InExpr{binOpExpr: binOpExpr{left: left, right: right, Token: tkn}}
}

// Accept marshalls the AST node walker to the correct visit method
func (n *InExpr) Accept(nw NodeWalker) WType { return nw.visitIn(n) }

func (n *InExpr) String() string { return "in" }

// AndExpr holds the '&&' operator between its 2 children
type AndExpr struct{ binOpExpr }

// newAnd returns a pointer to a AndExpr
func newAnd(left Expr, right Expr, tkn token.Token) *AndExpr {
	return &AndExpr{binOpExpr{left: left, right: right, Token: tkn}}
}

// Accept marshalls the AST node walker to the correct visit method
func (n *AndExpr) Accept(nw NodeWalker) WType { return nw.visitAnd(n) }

func (n *AndExpr) String() string { return "&&" }

// OrExpr holds the '||' operator between its 2 children
type OrExpr struct{ binOpExpr }

// newOr returns a pointer to a OrExpr
func newOr(left Expr, right Expr, tkn token.Token) *OrExpr {
	return &OrExpr{binOpExpr{left: left, right: right, Token: tkn}}
}

// Accept marshalls the AST node walker to the correct visit method
func (n *OrExpr) Accept(nw NodeWalker) WType { return nw.visitOr(n) }

func (n *OrExpr) String() string { return "||" }

// Unary Operators

// unOpExpr holds a unary operator as well as an operand node
type unOpExpr struct {
	token.Token
	scope   Scope
	operand Expr
}

func (n unOpExpr) Scope() Scope { return n.scope }

func (n unOpExpr) Expression() {}

// PlusExpr holds a unary positive ('+') operator and its operand
type PlusExpr struct{ unOpExpr }

// newPlus returns a pointer to a PlusExpr
func newPlus(operand Expr, tkn token.Token) *PlusExpr {
	return &PlusExpr{unOpExpr{operand: operand, Token: tkn}}
}

// Accept marshalls the AST node walker to the correct visit method
func (n *PlusExpr) Accept(nw NodeWalker) WType { return nw.visitPlus(n) }

func (n *PlusExpr) String() string { return "+" }

// MinusExpr holds a unary negative ('-') operator and its operand
type MinusExpr struct{ unOpExpr }

// newMinus returns a pointer to a MinusExpr
func newMinus(operand Expr, tkn token.Token) *MinusExpr {
	return &MinusExpr{unOpExpr{operand: operand, Token: tkn}}
}

// Accept marshalls the AST node walker to the correct visit method
func (n *MinusExpr) Accept(nw NodeWalker) WType { return nw.visitMinus(n) }

func (n *MinusExpr) String() string { return "-" }

// NotExpr holds a unary logical not ('!') operator and its operand
type NotExpr struct{ unOpExpr }

// newNot returns a pointer to a NotExpr
func newNot(operand Expr, tkn token.Token) *NotExpr {
	return &NotExpr{unOpExpr{operand: operand, Token: tkn}}
}

// Accept marshalls the AST node walker to the correct visit method
func (n *NotExpr) Accept(nw NodeWalker) WType { return nw.visitNot(n) }

func (n *NotExpr) String() string { return "!" }

// // Atom expressions
// type funcCall struct {
// }

/* Literals */

type literal struct {
	token.Token
	scope Scope
}

func (n literal) Scope() Scope { return n.scope }

func (n literal) Expression() {}

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

func (n *Num) String() string { return n.Text }

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

func (n *Str) String() string {
	return n.Value
}

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

func (n *Null) String() string { return n.Text }

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

func (n *Bool) String() string { return n.Text }

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

func (n *List) String() string { return fmt.Sprintf("%v", n.elements) }

// ID node represents Identifier/Name nodes
type ID struct {
	literal
	value string
}

func newID(value string, tkn token.Token) *ID { return &ID{literal{Token: tkn}, value} }

// Scope : self-explanatory
func (n *ID) Scope() Scope   { return n.scope }
func (n *ID) String() string { return n.value }

// Accept marshalls the AST node walker to the correct visit method
func (n *ID) Accept(nw NodeWalker) WType { return nw.visitID(n) }
