package golox

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (ap *AstPrinter) Print(expr Expr[string]) string {
	return expr.accept(ap)
}

func (ap *AstPrinter) visitBinaryExpr(expr *Binary[string]) string {
	return ap.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (ap *AstPrinter) visitGroupingExpr(expr *Grouping[string]) string {
	return ap.parenthesize("group", expr.expression)
}

func (ap *AstPrinter) visitLiteralExpr(expr *Literal[string]) string {
	if expr.value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.value)
}

func (ap *AstPrinter) visitUnaryExpr(expr *Unary[string]) string {
	return ap.parenthesize(expr.operator.lexeme, expr.right)
}

func (ap *AstPrinter) parenthesize(name string, exprs ...Expr[string]) string {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.accept(ap))
	}
	builder.WriteString(")")
	return builder.String()
}
