package lang

// NodeWalker is the interface to implement for all walkers/visitors to the AST
type NodeWalker interface {
	// Binary Operators
	visitAdd(*AddNode) WType
	visitSubtract(*SubtractNode) WType
	visitMult(*MultNode) WType
	visitDiv(*DivNode) WType
	visitMod(*ModNode) WType
	visitEq(*EqNode) WType
	visitSm(*SmNode) WType
	visitGr(*GrNode) WType
	visitIn(*InNode) WType
	visitAnd(*AndNode) WType
	visitOr(*OrNode) WType
	// Unary Operators
	visitPlus(*PlusNode) WType
	visitMinus(*MinusNode) WType
	visitNot(*NotNode) WType
	// visit literals
	visitNum(*NumberNode) WType
	visitStr(*StringNode) WType
	visitNull(*NullNode) WType
	visitBool(*BoolNode) WType
}
