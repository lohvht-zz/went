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
