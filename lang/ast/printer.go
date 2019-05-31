package ast

import "strings"

// Printer is an example of how to implement the ast.Visitor interface
type Printer struct{}

// Print returns the string value of the given AST via the Node.accept() method
func (v *Printer) Print(expr Expr) string {
	val := expr.Accept(v)
	s, ok := val.(string)
	if !ok {
		panic("NOT STRING VALUE")
	}
	return s
}

func (v *Printer) VisitNameDeclStmt(n *NameDeclStmt) interface{} {
	return nil
}

func (v *Printer) VisitExprStmt(n *ExprStmt) interface{} {
	return nil
}

func (v *Printer) VisitNameExpr(n *NameExpr) interface{} {
	return nil
}

func (v *Printer) VisitGrpExpr(n *GrpExpr) interface{} {
	return v.surroundBracket("group", n.Expression)
}
func (v *Printer) VisitBinExpr(n *BinExpr) interface{} {
	return v.surroundBracket(n.Op.Value, n.Left, n.Right)
}
func (v *Printer) VisitUnExpr(n *UnExpr) interface{} {
	return v.surroundBracket(n.Op.Value, n.Operand)
}
func (v *Printer) VisitBasicLit(n *BasicLit) interface{} {
	return n.Text
}

func (v *Printer) surroundBracket(name string, exprs ...Expr) string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteString(" ")
		sb.WriteString(v.Print(expr))
	}
	sb.WriteString(")")
	return sb.String()
}
