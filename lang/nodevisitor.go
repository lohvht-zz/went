package lang

import (
	"fmt"
)

// NodeWalker is the interface to implement for all walkers/visitors to the AST
type NodeWalker interface {
	// Binary Operations
	visitAdd(*AddNode)
	visitUnaryOp(*UnaryOpNode)
	// visit literals
	visitNum(*NumberNode)
	visitStr(*StringNode)
	visitNull(*NullNode)
	visitBool(*BoolNode)
}

// Top level visit function, marshals the NodeWalker to their correct visitXxx
// method
func visit(node Node, nv NodeWalker) {
	// fmt.Printf("nv type: %T", nv)
	switch typedNode := node.(type) {
	case *AddNode:
		nv.visitAdd(typedNode)
	case *UnaryOpNode:
		nv.visitUnaryOp(typedNode)
	case *NumberNode:
		nv.visitNum(typedNode)
	case *StringNode:
		nv.visitStr(typedNode)
	case *NullNode:
		nv.visitNull(typedNode)
	case *BoolNode:
		nv.visitBool(typedNode)
	}
}

// Interpreter contains the already parsed inputs, as well as
// the various already defined variables/functions
// TODO: scopes
type Interpreter struct {
	Root Node
}

// NewInterpreter creates a new interpreter object with the root as the Node
// being passed in
func NewInterpreter(rootNode Node) *Interpreter {
	i := &Interpreter{Root: rootNode}
	return i
}

// Interpret walks the tree from its root, exploring its children while making
// its walk downwards
func (i *Interpreter) Interpret() {
	visit(i.Root, i)
}

func (i *Interpreter) visitAdd(node *AddNode) {
	visit(node.left, i)
	fmt.Print(node)
	visit(node.right, i)
	fmt.Println()
}

func (i *Interpreter) visitUnaryOp(node *UnaryOpNode) {
	fmt.Print(node)
	visit(node.expr, i)
}

func (i *Interpreter) visitNum(node *NumberNode) {
	fmt.Print(node.Float64)
}

func (i *Interpreter) visitStr(node *StringNode) {
	fmt.Print(node.Value)
}

func (i *Interpreter) visitNull(node *NullNode) {
	fmt.Print(node.Value)
}

func (i *Interpreter) visitBool(node *BoolNode) {
	fmt.Print(node.Value)
}
