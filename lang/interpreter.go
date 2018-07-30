package lang

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
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

// error panics a general error
func (i *Interpreter) error(err error) {
	i.errorf("%s", err)
}

// typeError panics a type error
func (i *Interpreter) typeError(node Node, err error) {
	i.typeErrorf("%s", node, err)
}

func (i *Interpreter) recover(erri *error) {
	e := recover()
	if e != nil {
		if _, ok := e.(runtime.Error); ok {
			panic(e)
		}
		*erri = e.(error)
	}
}

// initInterp creates a new interpreter object with the root as the Node
// being passed in
func initInterp(rootNode Node) *Interpreter {
	i := &Interpreter{Root: rootNode}
	return i
}

// Interpret interprets the AST tree from its root
func Interpret(rootNode Node) (interp *Interpreter, err error) {
	i := initInterp(rootNode)
	defer i.recover(&err)
	i.interpret()
	return i, nil
}

// interpret walks the tree from its root, exploring its children while making
// its walk downwards
func (i *Interpreter) interpret() {
	res := i.Root.Accept(i)
	fmt.Printf("result is: %v of type %T\n", res, res)
}

// NOTE: Should we allow functional overloading for arithmetic expressions?

// additiveOp handles visit method for "additive" operators such as
// '+', '-', '*' for arithmetic operations
func (i *Interpreter) additiveOp(leftRes, rightRes WType, node Node) WType {
	a, aOk := leftRes.(WNum)
	b, bOk := rightRes.(WNum)
	if aOk && bOk {
		switch node.(type) {
		case *AddNode:
			return a + b
		case *SubtractNode:
			return a - b
		case *MultNode:
			return a * b
		}
	}
	// If reached here, force a type error, especially if they're adding in
	// incompatible types
	typ1Str := reflect.TypeOf(leftRes).Name()
	typ2Str := reflect.TypeOf(rightRes).Name()
	i.typeErrorf("unsupported operand type(s) for %s: '%s' and '%s'",
		node, node, typ1Str, typ2Str,
	)
	// Should not reach here as typeErrorf will panic
	return WNull{}
}

// divisiveOp handles visit method for "divisive" operators such as
// '/' and '%' for arithmetic operations such that they handle zero divisions
// properly
func (i *Interpreter) divisiveOp(leftRes, rightRes WType, node Node) WType {
	a, aOk := leftRes.(WNum)
	b, bOk := rightRes.(WNum)
	if aOk && bOk {
		if b.IsZeroValue() {
			if b.IsInt() {
				i.zeroDivisionErrorf("int division by zero", node)
			} else {
				i.zeroDivisionErrorf("float division by zero", node)
			}
		}
		switch node.(type) {
		case *DivNode:
			return a / b
		case *ModNode:
			if a.IsInt() && b.IsInt() {
				return WNum(int64(a) % int64(b))
			}
			return WNum(math.Mod(float64(a), float64(b)))
		}
	}
	// If reached here, force a type error, especially if they're adding in
	// incompatible types
	typ1Str := reflect.TypeOf(leftRes).Name()
	typ2Str := reflect.TypeOf(rightRes).Name()
	i.typeErrorf("unsupported operand type(s) for %s: '%s' and '%s'",
		node, node, typ1Str, typ2Str,
	)
	// Should not reach here as typeErrorf will panic
	return WNull{}
}

func (i *Interpreter) visitAdd(node *AddNode) WType {
	leftRes := node.left.Accept(i)
	rightRes := node.right.Accept(i)
	a, aOk := leftRes.(WString)
	b, bOk := rightRes.(WString)
	if aOk && bOk { // if they're both strings
		return a + b
	}
	return i.additiveOp(leftRes, rightRes, node)
}

func (i *Interpreter) visitSubtract(node *SubtractNode) WType {
	leftRes := node.left.Accept(i)
	rightRes := node.right.Accept(i)
	return i.additiveOp(leftRes, rightRes, node)
}

func (i *Interpreter) visitMult(node *MultNode) WType {
	leftRes := node.left.Accept(i)
	rightRes := node.right.Accept(i)
	return i.additiveOp(leftRes, rightRes, node)
}

func (i *Interpreter) visitDiv(node *DivNode) WType {
	leftRes := node.left.Accept(i)
	rightRes := node.right.Accept(i)
	return i.divisiveOp(leftRes, rightRes, node)
}

func (i *Interpreter) visitMod(node *ModNode) WType {
	leftRes := node.left.Accept(i)
	rightRes := node.right.Accept(i)
	return i.divisiveOp(leftRes, rightRes, node)
}

func (i *Interpreter) visitEq(node *EqNode) WType {
	leftRes := node.left.Accept(i)
	rightRes := node.right.Accept(i)
	if node.IsNot {
		return !leftRes.Equals(rightRes)
	}
	return leftRes.Equals(rightRes)
}

// visitSm evaluates '<' and '<=' operators
func (i *Interpreter) visitSm(node *SmNode) WType {
	leftRes := node.left.Accept(i)
	rightRes := node.right.Accept(i)

	smRes, err := leftRes.Sm(rightRes, node.OrEq)
	if err != nil {
		i.typeError(node, err)
	}
	return smRes
}

// visitGr evaluates '>' and '>=' operators
func (i *Interpreter) visitGr(node *GrNode) WType {
	leftRes := node.left.Accept(i)
	rightRes := node.right.Accept(i)

	grRes, err := leftRes.Gr(rightRes, node.OrEq)
	if err != nil {
		i.typeError(node, err)
	}
	return grRes
}

// TODO: confirm grammar spec for `in` keyword
func (i *Interpreter) visitIn(node *InNode) WType { return WNull{} }

// visitAnd evaluates '&&' operators
// if 'expr1 && expr2', expr1 if expr1 is false (i.e. zero-value), else expr2
func (i *Interpreter) visitAnd(node *AndNode) WType {
	leftRes := node.left.Accept(i)
	if leftRes.IsZeroValue() {
		return leftRes
	}
	return node.right.Accept(i)
}

// visitOr evaluates '||' operators
// if 'expr1 || expr2', expr2 if expr1 is false (i.e. zero-value), else expr2
func (i *Interpreter) visitOr(node *OrNode) WType {
	leftRes := node.left.Accept(i)
	if !leftRes.IsZeroValue() {
		return leftRes
	}
	return node.right.Accept(i)
}

// Unary Operators

// visitPlus evaluates a node
func (i *Interpreter) visitPlus(node *PlusNode) WType {
	switch v := node.operand.Accept(i).(type) {
	case WNum:
		return v
	default:
		typ := reflect.TypeOf(v).Name()
		i.typeErrorf("bad operand type for unary %s: '%s'", node, node, typ)
	}
	// Should not reach here as typeErrorf will panic
	return WNull{}
}

func (i *Interpreter) visitMinus(node *MinusNode) WType {
	switch v := node.operand.Accept(i).(type) {
	case WNum:
		return -v
	default:
		typ := reflect.TypeOf(v).Name()
		i.typeErrorf("bad operand type for unary %s: '%s'", node, node, typ)
	}
	// Should not reach here as typeErrorf will panic
	return WNull{}
}

// visitNot returns true if its operand are zero values (i.e. are false)
// else returns false
func (i *Interpreter) visitNot(node *NotNode) WType {
	switch v := node.operand.Accept(i).(type) {
	case WType:
		return v.IsZeroValue()
	default:
		typ := reflect.TypeOf(v).Name()
		i.typeErrorf("bad operand type for unary %s: '%s'", node, node, typ)
	}
	// Should not reach here as typeErrorf will panic
	return WNull{}
}

// visit literals ==> At its core, these will return WType values

// TODO: visit literals for maps
func (i *Interpreter) visitNum(n *NumberNode) WType { return WNum(n.Float64) }
func (i *Interpreter) visitStr(n *StringNode) WType { return WString(n.String()) }
func (i *Interpreter) visitNull(n *NullNode) WType  { return WNull{} }
func (i *Interpreter) visitBool(n *BoolNode) WType  { return WBool(n.Value) }

func (i *Interpreter) visitList(n *ListNode) WType {
	wl := WList{}
	for _, elNode := range n.elements {
		wl = append(wl, elNode.Accept(i))
	}
	return wl
}
