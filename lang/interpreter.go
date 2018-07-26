package lang

import (
	"fmt"
	"math"
	"reflect"
)

// Interpreter implements NodeWalker
// TODO: scopes
type Interpreter struct {
	Root Node
	name string // name of the interpreter, used for debugging purposes
}

// typeErrorf formats the error string before passing into errorf() for panicking
func (i *Interpreter) typeErrorf(format string, node Node, args ...interface{}) {
	format = fmt.Sprintf("Type Error: %s:%d: %s", i.name, node.LinePosition(), format)
	i.errorf(format, args...)
}

// zeroDivisionErrorf formats the error string before passing into errorf() for panicking
func (i *Interpreter) zeroDivisionErrorf(format string, node Node, args ...interface{}) {
	format = fmt.Sprintf("Zero Division Error: %s:%d: %s", i.name, node.LinePosition(), format)
	i.errorf(format, args...)
}

func (i *Interpreter) errorf(format string, args ...interface{}) {
	i.Root = nil // Discard the AST
	panic(fmt.Errorf(format, args...))
}

func (i *Interpreter) error(err error) {
	i.errorf("%s", err)
}

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

// NOTE: Should we allow functional overloading for arithmetic expressions?

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
	// If reached here, force a type error, especially if they're adding in
	// incompatible types
	typ1 := reflect.TypeOf(node.left).Name()
	typ2 := reflect.TypeOf(node.right).Name()
	i.typeErrorf(
		"unsupported operand types(s) for %s: '%s' and '%s'",
		node, node, typ1, typ2,
	)
	// Should not reach here as typeErrorf will panic
	return -1
}

func (i *Interpreter) visitSubtract(node *SubtractNode) interface{} {
	a, aOk := node.left.Accept(i).(*NumberNode)
	b, bOk := node.right.Accept(i).(*NumberNode)
	if aOk && bOk {
		if isIntOp(a, b) {
			return a.Int64 - b.Int64
		}
		return a.Float64 - b.Float64
	}
	// If reached here, force a type error, especially if they're adding in
	// incompatible types
	typ1 := reflect.TypeOf(node.left).Name()
	typ2 := reflect.TypeOf(node.right).Name()
	i.typeErrorf(
		"unsupported operand types(s) for %s: '%s' and '%s'",
		node, node, typ1, typ2,
	)
	// Should not reach here as typeErrorf will panic
	return -1
}

func (i *Interpreter) visitMult(node *MultNode) interface{} {
	a, aOk := node.left.Accept(i).(*NumberNode)
	b, bOk := node.right.Accept(i).(*NumberNode)
	if aOk && bOk {
		if isIntOp(a, b) {
			return a.Int64 * b.Int64
		}
		return a.Float64 * b.Float64
	}
	// If reached here, force a type error, especially if they're adding in
	// incompatible types
	typ1 := reflect.TypeOf(node.left).Name()
	typ2 := reflect.TypeOf(node.right).Name()
	i.typeErrorf(
		"unsupported operand types(s) for %s: '%s' and '%s'",
		node, node, typ1, typ2,
	)
	// Should not reach here as typeErrorf will panic
	return -1
}

func (i *Interpreter) visitDiv(node *DivNode) interface{} {
	a, aOk := node.left.Accept(i).(*NumberNode)
	b, bOk := node.right.Accept(i).(*NumberNode)
	if aOk && bOk {
		switch {
		case b.Float64 == 0:
			i.zeroDivisionErrorf("float division by zero", node)
		case b.Int64 == 0:
			i.zeroDivisionErrorf("int division by zero", node)
		case isIntOp(a, b):
			return a.Int64 / b.Int64
		default:
			return a.Float64 / b.Float64
		}
	}
	// If reached here, force a type error, especially if they're adding in
	// incompatible types
	typ1 := reflect.TypeOf(node.left).Name()
	typ2 := reflect.TypeOf(node.right).Name()
	i.typeErrorf(
		"unsupported operand types(s) for %s: '%s' and '%s'",
		node, node, typ1, typ2,
	)
	// Should not reach here as typeErrorf will panic
	return -1
}

func (i *Interpreter) visitMod(node *ModNode) interface{} {
	a, aOk := node.left.Accept(i).(*NumberNode)
	b, bOk := node.right.Accept(i).(*NumberNode)
	if aOk && bOk {
		switch {
		case b.Float64 == 0:
			i.zeroDivisionErrorf("float modulo by zero", node)
		case b.Int64 == 0:
			i.zeroDivisionErrorf("int modulo by zero", node)
		case isIntOp(a, b):
			return a.Int64 % b.Int64
		default:
			return math.Mod(a.Float64, b.Float64)
		}
	}
	// If reached here, force a type error, especially if they're adding in
	// incompatible types
	typ1 := reflect.TypeOf(node.left).Name()
	typ2 := reflect.TypeOf(node.right).Name()
	i.typeErrorf(
		"unsupported operand types(s) for %s: '%s' and '%s'",
		node, node, typ1, typ2,
	)
	// Should not reach here as typeErrorf will panic
	return -1
}

func (i *Interpreter) visitEq(node *EqNode) interface{}   { return -1 }
func (i *Interpreter) visitSm(node *SmNode) interface{}   { return -1 }
func (i *Interpreter) visitGr(node *GrNode) interface{}   { return -1 }
func (i *Interpreter) visitIn(node *InNode) interface{}   { return -1 }
func (i *Interpreter) visitAnd(node *AndNode) interface{} { return -1 }
func (i *Interpreter) visitOr(node *OrNode) interface{}   { return -1 }

// Unary Operators
func (i *Interpreter) visitPlus(node *PlusNode) interface{}   { return -1 }
func (i *Interpreter) visitMinus(node *MinusNode) interface{} { return -1 }
func (i *Interpreter) visitNot(node *NotNode) interface{}     { return -1 }

// visit literals

// visitNode for interpreter
func (i *Interpreter) visitNum(node *NumberNode) interface{} { return node }
func (i *Interpreter) visitStr(node *StringNode) interface{} { return node.String() }
func (i *Interpreter) visitNull(node *NullNode) interface{}  { return -1 }
func (i *Interpreter) visitBool(node *BoolNode) interface{}  { return -1 }

// Helper functions

// isIntOp determines if both number nodes should be evaluated as a omt64
// as opposed to a float
func isIntOp(a *NumberNode, b *NumberNode) bool {
	return a.IsInt && b.IsInt
}
