package runtime

import (
	"fmt"

	"github.com/lohvht/went/lang/ast"
	"github.com/lohvht/went/lang/token"
)

// Interpreter implements the ast.Visitor interface
type Interpreter struct {
	inputName string
	errors    token.ErrorList // runtime errors
}

func NewInterpreter(inputName string) *Interpreter {
	return &Interpreter{inputName: inputName}
}

// errorf formats the message and its arguments and should be favoured over using p.error
func (v *Interpreter) errorf(pos token.Pos, message string, msgArgs ...interface{}) {
	v.errors.Add(NewRuntimeError(v.inputName, pos, fmt.Sprintf(message, msgArgs...)))
	// log.Fatalln(p.errors[len(p.errors)-1])
}

func (v *Interpreter) Run(stmts []ast.Stmt) {
	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			fmt.Println(err.Error())
		}
	}()
	for _, stmt := range stmts {
		v.execute(stmt)
	}
}

func (v *Interpreter) execute(stmt ast.Stmt) {
	stmt.Accept(v)
}

func (v *Interpreter) evaluate(expr ast.Expr) interface{} { return expr.Accept(v) }

func (v *Interpreter) VisitExprStmt(stmt *ast.ExprStmt) interface{} {
	val := v.evaluate(stmt.Expression)
	// TODO: Add support to interpreter if not running in REPL mode to not print
	fmt.Println(stringify(val))
	return nil
}

func (v *Interpreter) VisitGrpExpr(n *ast.GrpExpr) interface{} {
	return v.evaluate(n.Expression)
}

func (v *Interpreter) VisitBinExpr(n *ast.BinExpr) interface{} {
	left := v.evaluate(n.Left)
	right := v.evaluate(n.Right)
	switch n.Op.Type {
	case token.PLUS:
		leftV, okl := left.(float64)
		rightV, okr := right.(float64)
		if okl && okr {
			return leftV + rightV
		}
		leftS, okl := left.(string)
		rightS, okr := right.(string)
		if okl && okr {
			return leftS + rightS
		}
		v.errorf(n.Op.Pos, "operands must be two numbers or two strings")
		panic(v.errors[len(v.errors)-1])
	case token.MINUS, token.DIV, token.MULT, token.GR, token.GREQ, token.SM, token.SMEQ:
		// TODO: Handle MOD types (change representation to separate between int and float?)
		fs, hasErr := v.checkFloatOperands(n.Op, left, right)
		if hasErr {
			panic(v.errors[len(v.errors)-1])
		}
		leftV := fs[0]
		rightV := fs[1]
		switch n.Op.Type {
		case token.MINUS:
			return leftV - rightV
		case token.DIV:
			// TODO: throw error here for ZeroDivisionError
			// One possible test is this: (0 / 0) == (0 / 0)
			// as per IEEE standard, any operation on NaN is false
			return leftV / rightV
		case token.MULT:
			return leftV * rightV
		case token.GR:
			return leftV > rightV
		case token.GREQ:
			return leftV >= rightV
		case token.SM:
			return leftV < rightV
		case token.SMEQ:
			return leftV <= rightV
		}
	case token.EQ:
		return v.isEqual(left, right)
	case token.NEQ:
		return !v.isEqual(left, right)
	}
	// Should be unreachable
	return nil
}

func (v *Interpreter) VisitUnExpr(n *ast.UnExpr) interface{} {
	operandVal := v.evaluate(n.Operand)
	switch n.Op.Type {
	case token.MINUS:
		fs, hasErr := v.checkFloatOperands(n.Op, operandVal)
		if hasErr {
			panic(v.errors[len(v.errors)-1])
		}
		return -fs[0]
	case token.PLUS:
		fs, hasErr := v.checkFloatOperands(n.Op, operandVal)
		if hasErr {
			panic(v.errors[len(v.errors)-1])
		}
		return fs[0]
	case token.LOGICALNOT:
		return !v.isTruthy(operandVal)
	}
	return nil
}

func (v *Interpreter) VisitBasicLit(n *ast.BasicLit) interface{} {
	return n.Value
}

func (v *Interpreter) isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}
	switch castVal := val.(type) {
	case bool:
		return castVal
	}
	return true
}

func (v *Interpreter) isEqual(a, b interface{}) bool { return a == b }

func (v *Interpreter) checkFloatOperands(op token.Token, operandVals ...interface{}) ([]float64, bool) {
	result := make([]float64, len(operandVals))
	for i, operandVal := range operandVals {
		f, ok := operandVal.(float64)
		if !ok {
			var s string
			var a string
			if len(operandVals) <= 1 {
				s = ""
				a = "a "
			} else {
				s = "s "
				a = ""
			}
			v.errorf(op.Pos, "operand%smust be %snumber%s", s, a, s)
			return nil, true
		}
		result[i] = f
	}
	return result, false
}

func stringify(val interface{}) string {
	if val == nil {
		return "null"
	}
	return fmt.Sprintf("%v", val)
}
