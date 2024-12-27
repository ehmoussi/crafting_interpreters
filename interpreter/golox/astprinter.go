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

func (ap *AstPrinter) Print(statements []Stmt[any]) string {
	errorMsg := "The construction of the AST failed"
	printStatements := ""
	for _, statement := range statements {
		value, err := statement.accept(ap)
		if err != nil {
			return errorMsg
		}
		valueString, ok := value.(string)
		if !ok {
			return errorMsg
		}
		printStatements += valueString
	}
	return printStatements
}

func (ap *AstPrinter) visitBlockStmt(expr *Block[any]) (any, error) {
	var builder strings.Builder
	builder.WriteString("{")
	for _, statement := range expr.statements {
		value, err := statement.accept(ap)
		if err != nil {
			return "", err
		}
		valueString, ok := value.(string)
		if !ok {
			return nil, errors.New("The construction of the AST failed")
		}
		builder.WriteString(valueString)
	}
	builder.WriteString("}")
	return builder.String(), nil
}

func (ap *AstPrinter) visitExpressionStmt(expr *Expression[any]) (any, error) {
	return expr.expression.accept(ap)
}

func (ap *AstPrinter) visitPrintStmt(expr *Print[any]) (any, error) {
	return ap.parenthesize("print", expr.expression)
}

func (ap *AstPrinter) visitVarStmt(expr *Var[any]) (any, error) {
	return ap.parenthesize("var "+expr.name.lexeme, expr.initializer)
}

func (ap *AstPrinter) visitAssignExpr(expr *Assign[any]) (any, error) {
	var builder strings.Builder
	builder.WriteString("(=")
	assignement, err := ap.parenthesize(expr.name.lexeme, expr.value)
	if err != nil {
		return "", err
	}
	builder.WriteString(assignement)
	builder.WriteString(")")
	return builder.String(), nil
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
			return "", errors.New("unexpected error: the expression is not a string")
		}
		builder.WriteString(eString)
	}
	builder.WriteString(")")
	return builder.String(), nil
}
