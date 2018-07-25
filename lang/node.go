package lang

import (
	"fmt"
	"strconv"
	"strings"
)

/**
* TODO: Refactor NewXxxNode to a factory pattern
 */
var textFormat = "%s" // change to "%q" in tests for better error messages

// Node is an element from the parse tree
type Node interface {
	Scope() Scope // returns the scope of a Node
	String() string
	Position() Pos // byte position of start of the node, in full original input string
}

// Pos represents the byte position in the original input text from which
// this template was parsed
type Pos int

// Position returns itself, provides an easy default implementation for
// embedding in a Node. Embedded in all non-trivial Nodes
func (p Pos) Position() Pos { return p }

// binaryOpNode holds a binary operator between a left and right node
// This struct is meant to be embedded within all other binary operation
// structs
// Logical Binary Operator: "||", "&&"
type binaryOpNode struct {
	Pos
	scope Scope
	left  Node
	right Node
}

func (n binaryOpNode) Scope() Scope { return n.scope }

// Arithmetic Binary Operators

// AddNode holds a '+' operator between its 2 children
type AddNode struct{ binaryOpNode }

// newAdd returns a pointer to a AddNode
func newAdd(left Node, right Node, pos Pos) *AddNode {
	return &AddNode{binaryOpNode{left: left, right: right, Pos: pos}}
}

func (n *AddNode) String() string { return "+" }

// SubtractNode holds a '-' operator between its 2 children
type SubtractNode struct{ binaryOpNode }

// newSubtract returns a pointer to a SubtractNode
func newSubtract(left Node, right Node, pos Pos) *SubtractNode {
	return &SubtractNode{binaryOpNode{left: left, right: right, Pos: pos}}
}

func (n *SubtractNode) String() string { return "-" }

// MultNode holds a '*' operator between its 2 children
type MultNode struct{ binaryOpNode }

// newMult returns a pointer to a MultNode
func newMult(left Node, right Node, pos Pos) *MultNode {
	return &MultNode{binaryOpNode{left: left, right: right, Pos: pos}}
}

func (n *MultNode) String() string { return "*" }

// DivNode holds a '/' operator between its 2 children
type DivNode struct{ binaryOpNode }

// newDiv returns a pointer to a DivNode
func newDiv(left Node, right Node, pos Pos) *DivNode {
	return &DivNode{binaryOpNode{left: left, right: right, Pos: pos}}
}

func (n *DivNode) String() string { return "/" }

// ModNode holds a '%' operator between its 2 children
type ModNode struct{ binaryOpNode }

// newMod returns a pointer to a ModNode
func newMod(left Node, right Node, pos Pos) *ModNode {
	return &ModNode{binaryOpNode{left: left, right: right, Pos: pos}}
}

func (n *ModNode) String() string { return "%" }

// Comparative Binary Operators

// EqNode holds either the '!=' or '==' operator between its 2 children
type EqNode struct {
	binaryOpNode
	IsNot bool
}

// newEq returns a pointer to a EqNode
func newEq(left Node, right Node, pos Pos, isNot bool) *EqNode {
	return &EqNode{binaryOpNode: binaryOpNode{left: left, right: right, Pos: pos}, IsNot: isNot}
}

func (n *EqNode) String() string {
	if n.IsNot {
		return "!="
	}
	return "=="
}

// SmNode holds either the '<' or '<=' operator between its 2 children
type SmNode struct {
	binaryOpNode
	AndEq bool
}

// newSm returns a pointer to a SmNode
func newSm(left Node, right Node, pos Pos, andEq bool) *SmNode {
	return &SmNode{binaryOpNode: binaryOpNode{left: left, right: right, Pos: pos}, AndEq: andEq}
}

func (n *SmNode) String() string {
	if n.AndEq {
		return "<="
	}
	return "<"
}

// GrNode holds either the '<' or '<=' operator between its 2 children
type GrNode struct {
	binaryOpNode
	AndEq bool
}

// newGr returns a pointer to a GrNode
func newGr(left Node, right Node, pos Pos, andEq bool) *GrNode {
	return &GrNode{binaryOpNode: binaryOpNode{left: left, right: right, Pos: pos}, AndEq: andEq}
}

func (n *GrNode) String() string {
	if n.AndEq {
		return ">="
	}
	return ">"
}

// InNode holds either the '!in' or 'in' operator between its 2 children
type InNode struct {
	binaryOpNode
	IsNot bool
}

// newIn returns a pointer to a InNode
func newIn(left Node, right Node, pos Pos, isNot bool) *InNode {
	return &InNode{binaryOpNode: binaryOpNode{left: left, right: right, Pos: pos}, IsNot: isNot}
}

func (n *InNode) String() string {
	if n.IsNot {
		return "!in"
	}
	return "in"
}

// AndNode holds the '&&' operator between its 2 children
type AndNode struct{ binaryOpNode }

// newAnd returns a pointer to a AndNode
func newAnd(left Node, right Node, pos Pos) *AndNode {
	return &AndNode{binaryOpNode{left: left, right: right, Pos: pos}}
}

func (n *AndNode) String() string { return "&&" }

// OrNode holds the '||' operator between its 2 children
type OrNode struct{ binaryOpNode }

// newOr returns a pointer to a OrNode
func newOr(left Node, right Node, pos Pos) *OrNode {
	return &OrNode{binaryOpNode{left: left, right: right, Pos: pos}}
}

func (n *OrNode) String() string { return "||" }

// Unary Operators

// unaryOpNode holds a unary operator as well as an operand node
type unaryOpNode struct {
	Pos
	scope   Scope
	operand Node
}

func (n unaryOpNode) Scope() Scope {
	return n.scope
}

// PlusNode holds a unary positive ('+') operator and its operand
type PlusNode struct{ unaryOpNode }

// newPlus returns a pointer to a PlusNode
func newPlus(operand Node, pos Pos) *PlusNode {
	return &PlusNode{unaryOpNode{operand: operand, Pos: pos}}
}

func (n *PlusNode) String() string { return "+" }

// MinusNode holds a unary negative ('-') operator and its operand
type MinusNode struct{ unaryOpNode }

// newMinus returns a pointer to a MinusNode
func newMinus(operand Node, pos Pos) *MinusNode {
	return &MinusNode{unaryOpNode{operand: operand, Pos: pos}}
}

func (n *MinusNode) String() string { return "-" }

// NotNode holds a unary logical not ('!') operator and its operand
type NotNode struct{ unaryOpNode }

// newNot returns a pointer to a NotNode
func newNot(operand Node, pos Pos) *NotNode {
	return &NotNode{unaryOpNode{operand: operand, Pos: pos}}
}

func (n *NotNode) String() string { return "!" }

// Literals

type litNode struct {
	Pos
	scope Scope
}

func (n litNode) Scope() Scope {
	return n.scope
}

// NumberNode holds a numerical constant: signed integer or float
type NumberNode struct {
	litNode
	IsInt   bool    // number has an int value
	IsFloat bool    // number has floating point value
	Int64   int64   // signed integer value
	Float64 float64 // floating point value
	Text    string  // Original text representation from input
}

// newNumber creates a new pointer to the NumberNode
func newNumber(pos Pos, text string) (*NumberNode, error) {
	n := &NumberNode{litNode: litNode{Pos: pos}, Text: text}
	i, err := strconv.ParseInt(text, 0, 64)
	// If an int extraction succeeded, promote the float
	if err == nil {
		n.IsInt = true
		n.Int64 = i
	}
	// If an integer extraction is successful, promote the float
	if n.IsInt {
		n.IsFloat = true
		n.Float64 = float64(n.Int64)
	} else {
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

func (n *NumberNode) String() string {
	return n.Text
}

// StringNode holds a string literal: both raw and quoted
type StringNode struct {
	litNode
	Value string
}

// newString creates a new pointer to the StringNode
func newString(pos Pos, text string) *StringNode {
	return &StringNode{litNode: litNode{Pos: pos}, Value: text}
}

func (n *StringNode) String() string {
	return n.Value
}

// NullNode holds a null literal
type NullNode struct {
	litNode
	Text string
}

// newNull creates a new pointer to the NullNode
func newNull(pos Pos, text string) *NullNode {
	return &NullNode{litNode: litNode{Pos: pos}, Text: text}
}

func (n *NullNode) String() string {
	return n.Text
}

// BoolNode holds a boolean literal
type BoolNode struct {
	litNode
	Value bool
	Text  string
}

// newBool creates a new pointer to the BoolNode
func newBool(pos Pos, text string, tknTyp tokenType) (*BoolNode, error) {
	switch tknTyp {
	case tokenTrue:
		return &BoolNode{litNode: litNode{Pos: pos}, Value: true, Text: text}, nil
	case tokenFalse:
		return &BoolNode{litNode: litNode{Pos: pos}, Value: false, Text: text}, nil
	default:
		return nil, fmt.Errorf("illegal bool syntax: %q", text)
	}
}

func (n *BoolNode) String() string { return n.Text }
