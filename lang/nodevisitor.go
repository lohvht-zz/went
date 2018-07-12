package utils

import (
	"fmt"
)

// NodeWalker is the interface to implement for all walkers/visitors to the AST
type NodeWalker interface {
	visit(Node) // Generic visit function to marshal visits into the right Node types
	visitNum(*NumberNode)
	visitBinOp(*BinaryOpNode)
	visitUnaryOp(*UnaryOpNode)
}

// NodeVisitor provides a default implementation for visit() over all other visitors
// the default implementation of the other methods visitXxx should be overriden
// when embedding NodeVisitor in other Visitor structs (If not overriden, visitXxx
// will not do any action and terminate the walking there)
type NodeVisitor struct {
}

func (nv *NodeVisitor) visit(node Node) {
	switch typedNode := node.(type) {
	case *NumberNode:
		nv.visitNum(typedNode)
	case *BinaryOpNode:
		nv.visitBinOp(typedNode)
	case *UnaryOpNode:
		nv.visitUnaryOp(typedNode)
	}
}

func (nv *NodeVisitor) visitNum(node *NumberNode)      {}
func (nv *NodeVisitor) visitBinOp(node *BinaryOpNode)  {}
func (nv *NodeVisitor) visitUnaryOp(node *UnaryOpNode) {}

// Interpreter contains the already parsed inputs, as well as
// the various already defined variables/functions
// TODO: scopes
type Interpreter struct {
	NodeVisitor
}

func (i *Interpreter) visitNum(node *NumberNode) {
	fmt.Print(node.Float64)
}

func (i *Interpreter) visitBinOp(node *BinaryOpNode) {
	i.visit(node.left)
	fmt.Print(node)
	i.visit(node.right)
	fmt.Println()
}

func (i *Interpreter) visitUnaryOp(node *UnaryOpNode) {
	fmt.Print(node)
	i.visit(node.expr)
}
