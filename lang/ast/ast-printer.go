package ast

import "strings"

type AstPrinter struct{}

func (v *AstPrinter) Print(expr Expr) string {
	val := expr.accept(v)
	s, ok := val.(string)
	if !ok {
		panic("NOT STRING VALUE")
	}
	return s
}

func (v *AstPrinter) visitGrpExpr(n *GrpExpr) interface{} {
	return v.surroundBracket("group", n.Expression)
}
func (v *AstPrinter) visitBinExpr(n *BinExpr) interface{} {
	return v.surroundBracket(n.Op.Value, n.Left, n.Right)
}
func (v *AstPrinter) visitUnExpr(n *UnExpr) interface{} {
	return v.surroundBracket(n.Op.Value, n.Operand)
}
func (v *AstPrinter) visitBasicLit(n *BasicLit) interface{} {
	return n.Text
}

func (v *AstPrinter) surroundBracket(name string, exprs ...Expr) string {
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
