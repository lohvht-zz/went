package lang

/**
 * NOTE: This file is deprecated and should be cleaned up (i.e. Deleted) as soon
 * as the main language is written
 * It's being kept around as a reference while the lox implementation is being
 * written
 */

// NodeWalker is the interface to implement for all walkers/visitors to the AST
type NodeWalker interface {

	// Statements
	visitExprStmt(*ExprStmt) WType
	visitAssignStmt(*AssignStmt) WType
	visitPlusAssignStmt(*PlusAssignStmt) WType
	visitMinusAssignStmt(*MinusAssignStmt) WType
	visitDivAssignStmt(*DivAssignStmt) WType
	visitMultAssignStmt(*MultAssignStmt) WType
	visitModAssignStmt(*ModAssignStmt) WType

	// Expressions

	// Binary Expressions

	visitBinExpr(*BinExpr) WType
	// visitAdd(*AddExpr) WType
	// visitSubtract(*SubtractExpr) WType
	// visitMult(*MultExpr) WType
	// visitDiv(*DivExpr) WType
	// visitMod(*ModExpr) WType
	// visitEq(*EqExpr) WType
	// visitSm(*SmExpr) WType
	// visitGr(*GrExpr) WType
	// visitIn(*InExpr) WType
	// visitAnd(*AndExpr) WType
	// visitOr(*OrExpr) WType

	// Unary Expressions

	visitUnExpr(*UnExpr) WType
	// visitPlus(*PlusExpr) WType
	// visitMinus(*MinusExpr) WType
	// visitNot(*NotExpr) WType

	// visit literals
	// visitNum(*Num) WType
	// visitStr(*Str) WType
	// visitNull(*Null) WType
	// visitBool(*Bool) WType

	visitBasicLit(*BasicLit) WType
	visitList(*List) WType
	visitDict(n *Dict) WType
	visitID(*Ident) WType
}
