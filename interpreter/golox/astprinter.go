package golox

import (
	"errors"
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (ap *AstPrinter) Print(expr Expr[any]) string {
	value, err := expr.accept(ap)
	if err != nil {
		return err.Error()
	}
	valueString, ok := value.(string)
	if !ok {
		return "Unexpected error: the expression is not a string"
	}
	return valueString
}

func (ap *AstPrinter) visitBinaryExpr(expr *Binary[any]) (any, error) {
	return ap.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (ap *AstPrinter) visitGroupingExpr(expr *Grouping[any]) (any, error) {
	return ap.parenthesize("group", expr.expression)
}

func (ap *AstPrinter) visitLiteralExpr(expr *Literal[any]) (any, error) {
	if expr.value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.value), nil
}

func (ap *AstPrinter) visitUnaryExpr(expr *Unary[any]) (any, error) {
	return ap.parenthesize(expr.operator.lexeme, expr.right)
}

func (ap *AstPrinter) visitVariableExpr(expr *Variable[any]) (any, error) {
	return expr.name, nil
}

func (ap *AstPrinter) parenthesize(name string, exprs ...Expr[any]) (string, error) {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		e, err := expr.accept(ap)
		if err != nil {
			return "", err
		}
		eString, ok := e.(string)
		if !ok {
			return "", errors.New("Unexpected error: the expression is not a string")
		}
		builder.WriteString(eString)
	}
	builder.WriteString(")")
	return builder.String(), nil
}
