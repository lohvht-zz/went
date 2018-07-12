package lang

import (
	"fmt"
	"strconv"
	"strings"
)

var textFormat = "%s" // change to "%q" in tests for better error messages

// Node is an element from the parse tree
type Node interface {
	String() string
	Position() Pos // byte position of start of the node, in full original input string
	// Most nodes should also implement a CopyXxx method for each specific NodeType
}

// Pos represents the byte position in the original input text from which
// this template was parsed
type Pos int

// Position returns itself, provides an easy default implementation for
// embedding in a Node. Embedded in all non-trivial Nodes
func (p Pos) Position() Pos {
	return p
}

// NumberNode holds a numerical constant: signed integer or float
// value is parsed and stored under all the types that can be represent
// the value.
// NOTE: Do not create a Node directly using the struct definition, newXxx methods
// exist to format a new node for you already.
type NumberNode struct {
	Pos
	IsInt   bool    // number has an int value
	IsFloat bool    // number has floating point value
	Int64   int64   // signed integer value
	Float64 float64 // floating point value
	Text    string  // Original text representation from input
}

// NewNumber creates a new NumberNode
// Should we include Uint? Complex numbers?
// https://golang.org/src/text/template/parse/node.go ::: line 533
func NewNumber(pos Pos, text string) (*NumberNode, error) {
	n := &NumberNode{Pos: pos, Text: text}
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
		return nil, fmt.Errorf("Illegal number syntax: %q", text)
	}
	return n, nil
}

func (n *NumberNode) String() string {
	return n.Text
}

// BinaryOpNode holds a binary operator between a left and right node
// such operations can include the following:
// Arithmetic: "+", "-", "/", "*", "%"
// Comparison: "==", "!=", "<", "<=", ">", ">="
// Membership check: "!in", "in"
// Logical Binary Operator: "||", "&&"
// NOTE: Do not create a Node directly using the struct definition, newXxx methods
// exist to format a new node for you already.
type BinaryOpNode struct {
	Pos
	Op    token
	left  Node
	right Node
}

// NewBinaryOp creates a new BinaryOpNode
func NewBinaryOp(left Node, op token, right Node, pos Pos) *BinaryOpNode {
	n := &BinaryOpNode{Pos: pos}
	n.left = left
	n.Op = op
	n.right = right
	return n
}

func (n *BinaryOpNode) String() string {
	return n.Op.String()
}

// UnaryOpNode holds a unary operator as well as an expression node
type UnaryOpNode struct {
	Pos
	Op   token
	expr Node
}

// NewUnaryOp creates a new UnaryOpNode
func NewUnaryOp(op token, expr Node, pos Pos) *UnaryOpNode {
	n := &UnaryOpNode{Pos: pos}
	n.Op = op
	n.expr = expr
	return n
}

func (n *UnaryOpNode) String() string {
	return n.Op.String()
}
