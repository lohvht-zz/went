package lang

import (
	"fmt"
)

// NodeWalker is the interface to implement for all walkers/visitors to the AST
type NodeWalker interface {
	visitNum(*NumberNode)
	visitBinOp(*BinaryOpNode)
	visitUnaryOp(*UnaryOpNode)
}

// Top level visit function, marshals the NodeWalker to their correct visitXxx
// method
func visit(node Node, nv NodeWalker) {
	// fmt.Printf("nv type: %T", nv)
	switch typedNode := node.(type) {
	case *NumberNode:
		// fmt.Printf("Go to number: %T %v", typedNode, typedNode)
		nv.visitNum(typedNode)
	case *BinaryOpNode:
		// fmt.Printf("Go to bin op: %T %v", typedNode, typedNode)
		nv.visitBinOp(typedNode)
	case *UnaryOpNode:
		// fmt.Printf("Go to un op: %T %v", typedNode, typedNode)
		nv.visitUnaryOp(typedNode)
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

func (i *Interpreter) visitNum(node *NumberNode) {
	fmt.Print(node.Float64)
}

func (i *Interpreter) visitBinOp(node *BinaryOpNode) {
	visit(node.left, i)
	fmt.Print(node)
	visit(node.right, i)
	fmt.Println()
}

func (i *Interpreter) visitUnaryOp(node *UnaryOpNode) {
	fmt.Print(node)
	visit(node.expr, i)
}
