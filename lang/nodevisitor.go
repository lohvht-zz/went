package lang

// NodeWalker is the interface to implement for all walkers/visitors to the AST
type NodeWalker interface {
	// Binary Operators
	visitAdd(*AddNode) interface{}
	visitSubtract(*SubtractNode) interface{}
	visitMult(*MultNode) interface{}
	visitDiv(*DivNode) interface{}
	visitMod(*ModNode) interface{}
	visitEq(*EqNode) interface{}
	visitSm(*SmNode) interface{}
	visitGr(*GrNode) interface{}
	visitIn(*InNode) interface{}
	visitAnd(*AndNode) interface{}
	visitOr(*OrNode) interface{}
	// Unary Operators
	visitPlus(*PlusNode) interface{}
	visitMinus(*MinusNode) interface{}
	visitNot(*NotNode) interface{}
	// visit literals
	visitNum(*NumberNode) interface{}
	visitStr(*StringNode) interface{}
	visitNull(*NullNode) interface{}
	visitBool(*BoolNode) interface{}
}

// Interpreter contains the already parsed inputs, as well as
// the various already defined variables/functions
// TODO: scopes
type Interpreter struct {
	Root Node
	name string // name of the interpreter, used for debugging purposes
}

// // Formats the error string before passing into errorf for panicking
// func (i *Interpreter) typeErrorf(format string, node Node, args ...interface{}) {
// 	format = fmt.Sprintf("Type Error: %s:%d: %s", i.name, node.Position(), format)
// 	i.errorf("%s", format)
// }

// func (i *Interpreter) errorf(format string, args ...interface{}) {
// 	i.Root = nil // Discard the AST
// 	format = fmt.Sprintf()
// }

// // NewInterpreter creates a new interpreter object with the root as the Node
// // being passed in
// func NewInterpreter(rootNode Node) *Interpreter {
// 	i := &Interpreter{Root: rootNode}
// 	return i
// }

// // Interpret walks the tree from its root, exploring its children while making
// // its walk downwards
// func (i *Interpreter) Interpret() {
// 	visit(i.Root, i)
// }

func (i *Interpreter) visitAdd(node *AddNode) interface{} {
	a, aOk := node.left.Accept(i).(string)
	b, bOk := node.right.Accept(i).(string)
	if aOk && bOk { // if they're both strings
		return a + b
	}
	c, cOk := node.left.Accept(i).(*NumberNode)
	d, dOk := node.right.Accept(i).(*NumberNode)
	if cOk && dOk {
		if isIntOp(c, d) {
			return c.Int64 + d.Int64
		}
		return c.Float64 + d.Float64
	}
	return -1
}

func (i *Interpreter) visitSubtract(node *SubtractNode) interface{} { return -1 }
func (i *Interpreter) visitMult(node *MultNode) interface{}         { return -1 }
func (i *Interpreter) visitDiv(node *DivNode) interface{}           { return -1 }
func (i *Interpreter) visitMod(node *ModNode) interface{}           { return -1 }
func (i *Interpreter) visitEq(node *EqNode) interface{}             { return -1 }
func (i *Interpreter) visitSm(node *SmNode) interface{}             { return -1 }
func (i *Interpreter) visitGr(node *GrNode) interface{}             { return -1 }
func (i *Interpreter) visitIn(node *InNode) interface{}             { return -1 }
func (i *Interpreter) visitAnd(node *AndNode) interface{}           { return -1 }
func (i *Interpreter) visitOr(node *OrNode) interface{}             { return -1 }

// Unary Operators
func (i *Interpreter) visitPlus(node *PlusNode) interface{}   { return -1 }
func (i *Interpreter) visitMinus(node *MinusNode) interface{} { return -1 }
func (i *Interpreter) visitNot(node *NotNode) interface{}     { return -1 }

// visit literals
func (i *Interpreter) visitNum(node *NumberNode) interface{} { return -1 }
func (i *Interpreter) visitStr(node *StringNode) interface{} { return node.String() }
func (i *Interpreter) visitNull(node *NullNode) interface{}  { return -1 }
func (i *Interpreter) visitBool(node *BoolNode) interface{}  { return -1 }

// Helper functions

// isIntOp determines if both number nodes should be evaluated as a omt64
// as opposed to a float
func isIntOp(a *NumberNode, b *NumberNode) bool {
	return a.IsInt && b.IsInt
}
